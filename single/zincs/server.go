package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"compress/gzip"

	"github.com/andybalholm/brotli"
	"github.com/gorilla/mux"
	"github.com/zincium/zinc/modules/base"
	"github.com/zincium/zinc/modules/shadow"
)

// ServerOptions server options
type ServerOptions struct {
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Server server
type Server struct {
	r       *mux.Router
	s       *http.Server
	options *ServerOptions
	root    string
	gitPath string
}

// NewServer todo
func NewServer(opt *Options) *Server {
	return &Server{root: opt.Root, gitPath: opt.GitPath, options: &ServerOptions{
		IdleTimeout:  opt.idleTimeout,
		ReadTimeout:  opt.readTimeout,
		WriteTimeout: opt.writeTimeout}}
}

func (s *Server) resolvePath(u, r string) (string, error) {
	gitdir := filepath.Join(s.root, u, r)
	if !strings.HasSuffix(gitdir, ".git") {
		gitdir += ".git"
	}
	if _, err := os.Stat(gitdir); err != nil {
		sugar.Errorf("access %s error %v", gitdir, err)
		return "", err
	}
	return gitdir, nil
}

func renderNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("repository not found"))
}

func renderInternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func (s *Server) lookupReferences(w http.ResponseWriter, r *http.Request, serviceName string, cmd *exec.Cmd) {
	cmd.Env = os.Environ()
	version := r.Header.Get("Git-Protocol")
	if len(version) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+version)
	}
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		renderInternalError(w, err)
		return
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		renderInternalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/x-git-"+serviceName+"-advertisement")
	if len(version) == 0 {
		// smart protocol
		fmt.Fprintf(w, "%04x# service=git-%s\n0000", len(serviceName)+4+15, serviceName)
	}
	if _, err := io.Copy(w, stdout); err != nil && err != io.EOF {
		sugar.Errorf("io.Copy %v")
	}
}

func (s *Server) handleUploadPackRefs(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	gitdir, err := s.resolvePath(user, repo)
	if err != nil {
		renderNotFound(w)
		return
	}
	cmd := exec.Command(s.gitPath, "upload-pack", "--advertise-refs", "--stateless-rpc", gitdir)
	s.lookupReferences(w, r, "upload-pack", cmd)
}

func (s *Server) handleReceivePackRefs(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	gitdir, err := s.resolvePath(user, repo)
	if err != nil {
		renderNotFound(w)
		return
	}
	cmd := exec.Command(s.gitPath, "receive-pack", "--advertise-refs", "--stateless-rpc", gitdir)
	s.lookupReferences(w, r, "receive-pack", cmd)
}

func (s *Server) exchangeInputOutput(w http.ResponseWriter, r *http.Request, serviceName string, cmd *exec.Cmd) {
	cmd.Env = os.Environ()
	version := r.Header.Get("Git-Protocol")
	if len(version) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+version)
	}
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		renderInternalError(w, err)
		return
	}
	defer stdout.Close()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		renderInternalError(w, err)
		return
	}
	defer stdin.Close()
	if err := cmd.Start(); err != nil {
		renderInternalError(w, err)
		return
	}
	var in io.Reader
	enc := r.Header.Get("Content-Encoding")
	if strings.Contains(enc, "gzip") {
		gr, err := gzip.NewReader(r.Body)
		if err != nil {
			renderInternalError(w, err)
			return
		}
		defer gr.Close()
		in = gr
	} else if strings.Contains(enc, "br") {
		in = brotli.NewReader(r.Body)
	} else {
		in = r.Body
	}
	w.Header().Set("Content-Type", "application/x-git-"+serviceName+"-result")
	err = base.GroupExecute(
		func() error {
			_, err := base.Copy(w, stdout)
			return err
		},
		func() error {
			_, err := base.Copy(stdin, in)
			return err
		},
	)
	if err != nil {
		sugar.Errorf("exchange %v", err)
	}
}

func (s *Server) handleUploadPack(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	gitdir, err := s.resolvePath(user, repo)
	if err != nil {
		renderNotFound(w)
		return
	}
	cmd := exec.Command(s.gitPath, "upload-pack", "--stateless-rpc", gitdir)
	s.exchangeInputOutput(w, r, "upload-pack", cmd)
}

func (s *Server) handleReceivePack(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	gitdir, err := s.resolvePath(user, repo)
	if err != nil {
		renderNotFound(w)
		return
	}
	cmd := exec.Command(s.gitPath, "receive-pack", "--stateless-rpc", gitdir)
	s.exchangeInputOutput(w, r, "receive-pack", cmd)
}

// Shutdown shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}

// ListenAndServe listen and serve
func (s *Server) ListenAndServe(listen string) error {
	r := mux.NewRouter()
	r.HandleFunc("/{user}/{repo}/info/refs", s.handleUploadPackRefs).Queries("service", "git-upload-pack").Methods("GET")
	r.HandleFunc("/{user}/{repo}/info/refs", s.handleReceivePackRefs).Queries("service", "git-receive-pack").Methods("GET")
	r.HandleFunc("/{user}/{repo}/git-upload-pack", s.handleUploadPack).Methods("POST")
	r.HandleFunc("/{user}/{repo}/git-receive-pack", s.handleReceivePack).Methods("POST")
	s.r = r
	s.s = &http.Server{
		Addr:         listen,
		Handler:      s,
		IdleTimeout:  s.options.IdleTimeout,
		ReadTimeout:  s.options.ReadTimeout,
		WriteTimeout: s.options.WriteTimeout,
	}
	return s.s.ListenAndServe()
}

// ServeHTTP serve HTTP
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", ServerVersion)
	hw := shadow.NewResponseWriter(w)
	s.r.ServeHTTP(hw, r)
	sugar.Infof("%s %s %s status: %d body: %d spend: %v\n", hw.RequestID(), r.Method, r.RequestURI, hw.StatusCode(), hw.Written(), hw.Since())
}
