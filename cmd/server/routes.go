package server

import (
	"net/http"

	"github.com/DomZippilli/gcs-proxy-cloud-function/cmd/domain/file"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/justinas/alice"
)

type Handler struct {
	FileHandler file.Handler
	H2cHandler  http.Handler
}

func SetupRouter(r *chi.Mux, handler Handler) {
	middlewares := alice.New(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Method(http.MethodGet, "/healthcheck", middlewares.ThenFunc(handler.FileHandler.HealthCheck))
	r.Method(http.MethodGet, "/download/{id}", middlewares.ThenFunc(handler.FileHandler.DownloadFile))
	r.Method(http.MethodPost, "/upload", middlewares.ThenFunc(handler.FileHandler.UploadFile))
	r.Method(http.MethodPost, "/decodeToken", middlewares.ThenFunc(handler.FileHandler.VerifyAndDecodeToken))
	r.Method(http.MethodPost, "/upload/check", middlewares.ThenFunc(handler.FileHandler.UploadStatus))
	r.Method(http.MethodGet, "/*", handler.H2cHandler)
}
