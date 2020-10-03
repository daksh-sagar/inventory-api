package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// App is the entrypoint into our application and what controls the context of each request
type App struct {
	mux *chi.Mux
	log *log.Logger
}

// NewApp constructs an App to handle a set of routes.
func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
	}
}

// Handle associates a handler function with an HTTP Method and URL pattern.
func (a *App) Handle(method, url string, h http.HandlerFunc) {
	a.mux.MethodFunc(method, url, h)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
