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

	// Activity logs routes
	r.Route("/api/logs", func(r chi.Router) {
		r.Get("/", h.ListLogs)
		r.Post("/", h.CreateLog)
		r.Get("/{id}", h.GetLog)
		r.Put("/{id}", h.UpdateLog)
		r.Delete("/{id}", h.DeleteLog)
	})

	// Master data routes
	r.Route("/api/master/job-titles", func(r chi.Router) {
		r.Get("/", h.ListMasterJobTitles)
		r.Post("/", h.CreateMasterJobTitle)
		r.Put("/{id}", h.UpdateMasterJobTitle)
		r.Delete("/{id}", h.DeleteMasterJobTitle)
	})

	r.Route("/api/master/pics", func(r chi.Router) {
		r.Get("/", h.ListMasterPICs)
		r.Post("/", h.CreateMasterPIC)
		r.Put("/{id}", h.UpdateMasterPIC)
		r.Delete("/{id}", h.DeleteMasterPIC)
	})

	r.Route("/api/master/statuses", func(r chi.Router) {
		r.Get("/", h.ListMasterStatuses)
		r.Post("/", h.CreateMasterStatus)
		r.Put("/{id}", h.UpdateMasterStatus)
		r.Delete("/{id}", h.DeleteMasterStatus)
	})

	r.Route("/api/master/categories", func(r chi.Router) {
		r.Get("/", h.ListMasterCategories)
		r.Post("/", h.CreateMasterCategory)
		r.Put("/{id}", h.UpdateMasterCategory)
		r.Delete("/{id}", h.DeleteMasterCategory)
	})

	// File uploads
	fileServer := http.FileServer(http.Dir(uploadDir))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	// Static web files
	r.Handle("/*", http.FileServer(webFS))

	return r
}
