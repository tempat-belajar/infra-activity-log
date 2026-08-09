package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akr/infra-activity-log/internal/db"
	"github.com/akr/infra-activity-log/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	DB        *db.DB
	UploadDir string
	Logger    *slog.Logger
}

func New(database *db.DB, uploadDir string, logger *slog.Logger) *Handler {
	return &Handler{DB: database, UploadDir: uploadDir, Logger: logger}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListLogs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := db.ListFilter{
		DateFrom: q.Get("date_from"),
		DateTo:   q.Get("date_to"),
		PIC:      q.Get("pic"),
		Status:   q.Get("status"),
		Category: q.Get("category"),
		Search:   q.Get("search"),
	}
	logs, err := h.DB.ListLogs(r.Context(), f)
	if err != nil {
		h.Logger.Error("list logs failed", "err", err)
		writeErr(w, http.StatusInternalServerError, "failed to list logs")
		return
	}
	writeJSON(w, http.StatusOK, logs)
}

func (h *Handler) GetLog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	l, err := h.DB.GetLog(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "log not found")
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func validateRequired(tanggal, label string) string {
	if tanggal == "" || label == "" {
		return "tanggal and label are required"
	}
	return ""
}

const maxUploadSize = 5 << 20 // 5MB

func (h *Handler) saveUploadedFile(r *http.Request, field string) (string, error) {
	file, header, err := r.FormFile(field)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		return "", fmt.Errorf("file too large (max 5MB)")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", fmt.Errorf("only jpg/png allowed")
	}

	filename := uuid.New().String() + ext
	dstPath := filepath.Join(h.UploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return "/uploads/" + filename, nil
}

func (h *Handler) CreateLog(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize * 2); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid form data")
		return
	}

	l := models.ActivityLog{
		Tanggal:      r.FormValue("tanggal"),
		JobTitle:     r.FormValue("job_title"),
		PIC:          r.FormValue("pic"),
		Application:  r.FormValue("application"),
		Label:        r.FormValue("label"),
		OldValueText: r.FormValue("old_value_text"),
		NewValueText: r.FormValue("new_value_text"),
		Status:       r.FormValue("status"),
		Category:     r.FormValue("category"),
	}

	if msg := validateRequired(l.Tanggal, l.Label); msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	l.OldValueImageURL = oldImg

	newImg, err := h.saveUploadedFile(r, "new_value_image")
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	l.NewValueImageURL = newImg

	id, err := h.DB.CreateLog(r.Context(), l)
	if err != nil {
		h.Logger.Error("create log failed", "err", err)
		writeErr(w, http.StatusInternalServerError, "failed to create log")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *Handler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := r.ParseMultipartForm(maxUploadSize * 2); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid form data")
		return
	}

	l := models.ActivityLog{
		Tanggal:      r.FormValue("tanggal"),
		JobTitle:     r.FormValue("job_title"),
		PIC:          r.FormValue("pic"),
		Application:  r.FormValue("application"),
		Label:        r.FormValue("label"),
		OldValueText: r.FormValue("old_value_text"),
		NewValueText: r.FormValue("new_value_text"),
		Status:       r.FormValue("status"),
		Category:     r.FormValue("category"),
	}
	if msg := validateRequired(l.Tanggal, l.Label); msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}

	keepOld := true
	if oldImg, err := h.saveUploadedFile(r, "old_value_image"); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	} else if oldImg != "" {
		l.OldValueImageURL = oldImg
		keepOld = false
	}

	keepNew := true
	if newImg, err := h.saveUploadedFile(r, "new_value_image"); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	} else if newImg != "" {
		l.NewValueImageURL = newImg
		keepNew = false
	}

	if err := h.DB.UpdateLog(r.Context(), id, l, keepOld, keepNew); err != nil {
		h.Logger.Error("update log failed", "err", err)
		writeErr(w, http.StatusInternalServerError, "failed to update log")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.DB.DeleteLog(r.Context(), id); err != nil {
		h.Logger.Error("delete log failed", "err", err)
		writeErr(w, http.StatusInternalServerError, "failed to delete log")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
