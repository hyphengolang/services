package service

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hyphengolang/noughts-and-crosses/internal/conf"
	h "github.com/hyphengolang/prelude/http"
)

type Router interface {
	chi.Router

	Respond(w http.ResponseWriter, r *http.Request, data any, status int)
	Decode(w http.ResponseWriter, r *http.Request, data any) error

	SetLocation(w http.ResponseWriter, r *http.Request, location string)
	SetCookie(w http.ResponseWriter, r *http.Request, cookie *http.Cookie)

	Log(v ...any)
	Logf(format string, v ...any)

	ClientURI() string
}

func (*routerHandler) ClientURI() string { return conf.ClientURI }

func (*routerHandler) Decode(w http.ResponseWriter, r *http.Request, data any) error {
	return h.Decode(w, r, data)
}

func (*routerHandler) Respond(w http.ResponseWriter, r *http.Request, data any, status int) {
	h.Respond(w, r, data, status)
}

func (*routerHandler) SetLocation(w http.ResponseWriter, r *http.Request, location string) {
	var scheme string
	if r.TLS == nil {
		// the scheme was HTTP
		scheme = "http://"
	} else {
		// the scheme was HTTPS
		scheme = "https://"
	}

	w.Header().Set("Location", scheme+r.Host+location)
}

func (rh *routerHandler) SetCookie(w http.ResponseWriter, r *http.Request, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
}

func (r *routerHandler) Log(v ...any) {
	r.l.Println(v...)
}

func (r *routerHandler) Logf(format string, v ...any) {
	r.l.Printf(format, v...)
}

type routerHandler struct {
	chi.Router

	l *log.Logger
}

func NewRouter() Router {
	sh := routerHandler{
		Router: chi.NewRouter(),
		// NOTE - logger can be passed as an option
		l: log.Default(),
	}
	return &sh
}
