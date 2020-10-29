package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	srv     *http.Server
	options *ServerOptions
}

func (srv *Server) handleUploadPackRefs(w http.ResponseWriter, r *http.Request) {

}

func (srv *Server) handleReceivePackRefs(w http.ResponseWriter, r *http.Request) {

}

func (srv *Server) handleUploadPack(w http.ResponseWriter, r *http.Request) {

}

func (srv *Server) handleReceivePack(w http.ResponseWriter, r *http.Request) {

}

// ListenAndServe listen and serve
func (srv *Server) ListenAndServe(listen string) error {
	r := mux.NewRouter()
	r.HandleFunc("/{user}/{repo}/info/refs", srv.handleUploadPackRefs).Queries("service", "git-upload-pack").Methods("GET")
	r.HandleFunc("/{user}/{repo}/info/refs", srv.handleReceivePackRefs).Queries("service", "git-receive-pack").Methods("GET")
	r.HandleFunc("/{user}/{repo}/git-upload-pack", srv.handleUploadPack).Methods("POST")
	r.HandleFunc("/{user}/{repo}/git-receive-pack", srv.handleReceivePack).Methods("POST")
	srv.r = r
	srv.srv = &http.Server{
		Addr:         listen,
		Handler:      srv,
		IdleTimeout:  srv.options.IdleTimeout,
		ReadTimeout:  srv.options.ReadTimeout,
		WriteTimeout: srv.options.WriteTimeout,
	}
	return srv.srv.ListenAndServe()
}

// ServeHTTP serve HTTP
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", ServerVersion)
	hw := shadow.NewResponseWriter(w)
	srv.r.ServeHTTP(hw, r)
	log.Printf("%s %s %s status: %d body: %d spend: %v\n", hw.RequestID(), r.Method, r.RequestURI, hw.StatusCode(), hw.Written(), hw.Since())
}
