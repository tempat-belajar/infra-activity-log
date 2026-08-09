package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akr/infra-activity-log/internal/models"
	"github.com/go-chi/chi/v5"
)

// === Master Job Titles ===
func (h *Handler) ListMasterJobTitles(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Pool.Query(r.Context(),
		`SELECT id, name, is_active, created_at, updated_at FROM master_job_titles WHERE is_active = true ORDER BY name`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []models.MasterJobTitle
	for rows.Next() {
		var item models.MasterJobTitle
		if err := rows.Scan(&item.ID, &item.Name, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.MasterJobTitle{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateMasterJobTitle(w http.ResponseWriter, r *http.Request) {
	var item models.MasterJobTitle
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.DB.Pool.QueryRow(r.Context(),
		`INSERT INTO master_job_titles (name, is_active) VALUES ($1, true) RETURNING id, created_at, updated_at`,
		item.Name).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	item.IsActive = true
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMasterJobTitle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var item models.MasterJobTitle
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err := h.DB.Pool.Exec(r.Context(),
		`UPDATE master_job_titles SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, item.Name, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMasterJobTitle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.DB.Pool.Exec(r.Context(),
		`UPDATE master_job_titles SET is_active = false WHERE id = $1`, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// === Master PICs ===
func (h *Handler) ListMasterPICs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Pool.Query(r.Context(),
		`SELECT id, name, COALESCE(email, ''), is_active, created_at, updated_at FROM master_pics WHERE is_active = true ORDER BY name`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []models.MasterPIC
	for rows.Next() {
		var item models.MasterPIC
		if err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.MasterPIC{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateMasterPIC(w http.ResponseWriter, r *http.Request) {
	var item models.MasterPIC
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.DB.Pool.QueryRow(r.Context(),
		`INSERT INTO master_pics (name, email, is_active) VALUES ($1, $2, true) RETURNING id, created_at, updated_at`,
		item.Name, item.Email).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	item.IsActive = true
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMasterPIC(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var item models.MasterPIC
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err := h.DB.Pool.Exec(r.Context(),
		`UPDATE master_pics SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`,
		item.Name, item.Email, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMasterPIC(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.DB.Pool.Exec(r.Context(),
		`UPDATE master_pics SET is_active = false WHERE id = $1`, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// === Master Statuses ===
func (h *Handler) ListMasterStatuses(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Pool.Query(r.Context(),
		`SELECT id, name, COALESCE(color, ''), is_active, created_at, updated_at FROM master_statuses WHERE is_active = true ORDER BY name`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []models.MasterStatus
	for rows.Next() {
		var item models.MasterStatus
		if err := rows.Scan(&item.ID, &item.Name, &item.Color, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.MasterStatus{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateMasterStatus(w http.ResponseWriter, r *http.Request) {
	var item models.MasterStatus
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.DB.Pool.QueryRow(r.Context(),
		`INSERT INTO master_statuses (name, color, is_active) VALUES ($1, $2, true) RETURNING id, created_at, updated_at`,
		item.Name, item.Color).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	item.IsActive = true
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMasterStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var item models.MasterStatus
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err := h.DB.Pool.Exec(r.Context(),
		`UPDATE master_statuses SET name = $1, color = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`,
		item.Name, item.Color, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMasterStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.DB.Pool.Exec(r.Context(),
		`UPDATE master_statuses SET is_active = false WHERE id = $1`, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// === Master Categories ===
func (h *Handler) ListMasterCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Pool.Query(r.Context(),
		`SELECT id, name, COALESCE(description, ''), is_active, created_at, updated_at FROM master_categories WHERE is_active = true ORDER BY name`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []models.MasterCategory
	for rows.Next() {
		var item models.MasterCategory
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.MasterCategory{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateMasterCategory(w http.ResponseWriter, r *http.Request) {
	var item models.MasterCategory
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.DB.Pool.QueryRow(r.Context(),
		`INSERT INTO master_categories (name, description, is_active) VALUES ($1, $2, true) RETURNING id, created_at, updated_at`,
		item.Name, item.Description).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	item.IsActive = true
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMasterCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var item models.MasterCategory
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err := h.DB.Pool.Exec(r.Context(),
		`UPDATE master_categories SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`,
		item.Name, item.Description, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMasterCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.DB.Pool.Exec(r.Context(),
		`UPDATE master_categories SET is_active = false WHERE id = $1`, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
