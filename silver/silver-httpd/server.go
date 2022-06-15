package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"compress/gzip"

	"github.com/andybalholm/brotli"
	"github.com/gorilla/mux"
	"github.com/klauspost/compress/zstd"
	"github.com/zincium/zinc/modules/auth"
	"github.com/zincium/zinc/modules/base"
	"github.com/zincium/zinc/modules/env"
	"github.com/zincium/zinc/modules/shadow"
	"github.com/zincium/zinc/modules/simulate"
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
	ac      auth.AuthorizationClient
}

// NewServer todo
func NewServer(opt *Options) *Server {
	return &Server{root: opt.Root, gitPath: opt.GitPath, options: &ServerOptions{
		IdleTimeout:  opt.idleTimeout,
		ReadTimeout:  opt.readTimeout,
		WriteTimeout: opt.writeTimeout},
		ac: simulate.NewAuthorizationClient()}
}

func renderNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("repository not found"))
}

func renderInternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

type repoResult struct {
	RepoPath string
	Address  string
}

func (s *Server) doAuthFilter(w http.ResponseWriter, r *http.Request, action auth.Action) *repoResult {
	user := mux.Vars(r)["user"]
	repo := mux.Vars(r)["repo"]
	req := &auth.AuthRequest{
		RelativePath: path.Join(user, repo),
		Action:       action,
		Protocol:     auth.HTTP,
	}
	resp, err := s.ac.Authorize(req)
	if err != nil {
		renderInternalError(w, err)
		return nil
	}
	if !resp.Success() {
		resp.RenderHTML(w)
		return nil
	}
	gitdir := filepath.Join(s.root, resp.Location)
	if !strings.HasSuffix(gitdir, ".git") {
		gitdir += ".git"
	}
	if _, err := os.Stat(gitdir); err != nil {
		sugar.Errorf("access %s error %v", gitdir, err)
		renderNotFound(w)
		return nil
	}
	return &repoResult{
		RepoPath: gitdir,
		Address:  resp.Address,
	}
}

func (s *Server) lookupReferences(w http.ResponseWriter, r *http.Request, serviceName string, cmd *exec.Cmd) {
	cmd.Env = env.Environ()
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
	if _, err := io.Copy(w, stdout); err != nil && err != io.EOF {
		sugar.Errorf("io.Copy %v")
	}
}

func (s *Server) handleUploadPackRefs(w http.ResponseWriter, r *http.Request) {
	result := s.doAuthFilter(w, r, auth.Download)
	if result == nil {
		return
	}
	cmd := exec.Command(s.gitPath, "upload-pack", "--advertise-refs", "--stateless-rpc", result.RepoPath)
	s.lookupReferences(w, r, "upload-pack", cmd)
}

func (s *Server) handleReceivePackRefs(w http.ResponseWriter, r *http.Request) {
	result := s.doAuthFilter(w, r, auth.Upload)
	if result == nil {
		return
	}
	cmd := exec.Command(s.gitPath, "receive-pack", "--advertise-refs", "--stateless-rpc", result.RepoPath)
	s.lookupReferences(w, r, "receive-pack", cmd)
}

func (s *Server) exchangeInputOutput(w http.ResponseWriter, r *http.Request, serviceName string, cmd *exec.Cmd) {
	cmd.Env = env.Environ()
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
	defer func() {
		if stdin != nil {
			stdin.Close()
		}
	}()
	if err := cmd.Start(); err != nil {
		renderInternalError(w, err)
		return
	}
	var rc io.ReadCloser
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
	// Accept-Encoding: gzip, deflate
	// Content-Encoding: gzip
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		if rc, err = gzip.NewReader(r.Body); err != nil {
			return
		}
	case "br":
		rc = io.NopCloser(brotli.NewReader(r.Body))
	case "zstd":
		zr, err := zstd.NewReader(r.Body, nil)
		if err != nil {
			return
		}
		rc = zr.IOReadCloser()
	default:
		rc = io.NopCloser(r.Body)
	}
	defer rc.Close()
	if _, err := base.Copy(stdin, rc); err != nil {
		renderInternalError(w, err)
		return
	}
	if serviceName == "upload-pack" {
		stdin.Close()
		stdin = nil
	}
	w.Header().Set("Content-Type", "application/x-git-"+serviceName+"-result")
	if _, err := base.Copy(w, stdout); err != nil {
		sugar.Errorf("exchange %v", err)
	}
}

func (s *Server) handleUploadPack(w http.ResponseWriter, r *http.Request) {
	result := s.doAuthFilter(w, r, auth.Download)
	if result == nil {
		return
	}
	cmd := exec.Command(s.gitPath, "upload-pack", "--stateless-rpc", result.RepoPath)
	s.exchangeInputOutput(w, r, "upload-pack", cmd)
}

func (s *Server) handleReceivePack(w http.ResponseWriter, r *http.Request) {
	result := s.doAuthFilter(w, r, auth.Upload)
	if result == nil {
		return
	}
	cmd := exec.Command(s.gitPath, "receive-pack", "--stateless-rpc", result.RepoPath)
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
	sugar.Infof("%s %s %s status: %d body: %d spend: %v\n",
		hw.RequestID(), r.Method, r.RequestURI, hw.StatusCode(), hw.Written(), hw.Since())
}
