package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zincium/zinc/modules/uuid"
)

// App server
type App struct {
	r            *mux.Router
	srv          *http.Server
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (app *App) handleUploadPackRefs(w http.ResponseWriter, r *http.Request) {

}

func (app *App) handleReceivePackRefs(w http.ResponseWriter, r *http.Request) {

}

func (app *App) handleUploadPack(w http.ResponseWriter, r *http.Request) {

}

func (app *App) handleReceivePack(w http.ResponseWriter, r *http.Request) {

}

// ListenAndServe listen and serve
func (app *App) ListenAndServe(listen string) error {
	r := mux.NewRouter()
	r.HandleFunc("/{user}/{repo}/info/refs", app.handleUploadPackRefs).Queries("service", "git-upload-pack").Methods("GET")
	r.HandleFunc("/{user}/{repo}/info/refs", app.handleReceivePackRefs).Queries("service", "git-receive-pack").Methods("GET")
	r.HandleFunc("/{user}/{repo}/git-upload-pack", app.handleUploadPack).Methods("POST")
	r.HandleFunc("/{user}/{repo}/git-receive-pack", app.handleReceivePack).Methods("POST")
	app.r = r
	app.srv = &http.Server{Addr: listen, Handler: app, IdleTimeout: app.IdleTimeout, ReadTimeout: app.ReadTimeout, WriteTimeout: app.WriteTimeout}
	return app.srv.ListenAndServe()
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rid, err := uuid.NewRandom()
	if err == nil {
		w.Header().Add("X-Request-Id", rid.String())
	}
	w.Header().Set("Server", ServerVersion)
	now := time.Now()
	hw := NewHijackResponseWriter(w)
	app.r.ServeHTTP(hw, r)
	log.Printf("%s %s status: %d body: %d spend: %v\n", r.Method, r.RequestURI, hw.StatusCode(), hw.Written(), time.Since(now))
}
