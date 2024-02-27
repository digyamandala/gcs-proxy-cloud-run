package server

import (
	"net/http"

	"github.com/DomZippilli/gcs-proxy-cloud-function/cmd/domain/file"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

type Handler struct {
	FileHandler file.Handler
	H2cHandler  http.Handler
}

func SetupRouter(r *chi.Mux, handler Handler) {
	middlewares := alice.New(middleware.Recoverer)
	r.Method(http.MethodPost, "/upload", middlewares.ThenFunc(handler.FileHandler.UploadFile))
	r.Method(http.MethodGet, "/download/{id}", middlewares.ThenFunc(handler.FileHandler.DownloadFile))
	r.Method(http.MethodPost, "/decodeToken", middlewares.ThenFunc(handler.FileHandler.VerifyAndDecodeToken))
	r.Method(http.MethodGet, "/*", handler.H2cHandler)
}
