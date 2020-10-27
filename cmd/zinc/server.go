package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server server
type Server struct {
	srv          *http.Server
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ListenAndServe listen
func (s *Server) ListenAndServe(listen string) error {
	r := mux.NewRouter()
	r.HandleFunc("/{user}/{repo}/info/refs", s.handleUploadPackRefs).Queries("service", "git-upload-pack").Methods("GET")
	r.HandleFunc("/{user}/{repo}/info/refs", s.handleReceivePackRefs).Queries("service", "git-receive-pack").Methods("GET")
	r.HandleFunc("/{user}/{repo}/git-upload-pack", s.handleUploadPack).Methods("POST")
	r.HandleFunc("/{user}/{repo}/git-receive-pack", s.handleReceivePack).Methods("POST")
	s.srv = &http.Server{listen: listen, Handler: r, IdleTimeout: s.IdleTimeout, ReadTimeout: s.ReadTimeout, WriteTimeout: s.WriteTimeout}
	return s.srv.ListenAndServe()
}

func (s *Server) handleUploadPackRefs(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleReceivePackRefs(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleUploadPack(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleReceivePack(w http.ResponseWriter, r *http.Request) {

}
