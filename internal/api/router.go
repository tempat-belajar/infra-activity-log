package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler, uploadDir string, webFS http.FileSystem) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", h.Healthz)

	r.Route("/api/logs", func(r chi.Router) {
		r.Get("/", h.ListLogs)
		r.Post("/", h.CreateLog)
		r.Get("/{id}", h.GetLog)
		r.Put("/{id}", h.UpdateLog)
		r.Delete("/{id}", h.DeleteLog)
	})

	fileServer := http.FileServer(http.Dir(uploadDir))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	r.Handle("/*", http.FileServer(webFS))

	return r
}
