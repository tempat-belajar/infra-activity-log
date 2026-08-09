package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akr/infra-activity-log/internal/models"
	"github.com/go-chi/chi/v5"
)

// === Master Job Titles ===
func (h *Handler) ListMasterJobTitles(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, is_active, created_at, updated_at FROM master_job_titles WHERE is_active = true ORDER BY name`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateMasterJobTitle(w http.ResponseWriter, r *http.Request) {
	var item models.MasterJobTitle
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.DB.QueryRow(`INSERT INTO master_job_titles (name, is_active) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		item.Name, true).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.IsActive = true
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMasterJobTitle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var item models.MasterJobTitle
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.DB.Exec(`UPDATE master_job_titles SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, item.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMasterJobTitle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.DB.Exec(`UPDATE master_job_titles SET is_active = false WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// === Master PICs ===
func (h *Handler) ListMasterPICs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, email, is_active, created_at, updated_at FROM master_pics WHERE is_active = true ORDER BY name`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.MasterPIC
	for rows.Next() {
		var item models.MasterPIC
		var email sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &email, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		item.Email = email.String
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateMasterPIC(w http.ResponseWriter, r *http.Request) {
	var item models.MasterPIC
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.DB.QueryRow(`INSERT INTO master_pics (name, email, is_active) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		item.Name, item.Email, true).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.IsActive = true
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMasterPIC(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var item models.MasterPIC
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.DB.Exec(`UPDATE master_pics SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`, item.Name, item.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMasterPIC(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.DB.Exec(`UPDATE master_pics SET is_active = false WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// === Master Statuses ===
func (h *Handler) ListMasterStatuses(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, color, is_active, created_at, updated_at FROM master_statuses WHERE is_active = true ORDER BY name`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.MasterStatus
	for rows.Next() {
		var item models.MasterStatus
		var color sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &color, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		item.Color = color.String
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateMasterStatus(w http.ResponseWriter, r *http.Request) {
	var item models.MasterStatus
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.DB.QueryRow(`INSERT INTO master_statuses (name, color, is_active) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		item.Name, item.Color, true).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.IsActive = true
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMasterStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var item models.MasterStatus
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.DB.Exec(`UPDATE master_statuses SET name = $1, color = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`, item.Name, item.Color, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMasterStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.DB.Exec(`UPDATE master_statuses SET is_active = false WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// === Master Categories ===
func (h *Handler) ListMasterCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, description, is_active, created_at, updated_at FROM master_categories WHERE is_active = true ORDER BY name`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.MasterCategory
	for rows.Next() {
		var item models.MasterCategory
		var desc sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &desc, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		item.Description = desc.String
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateMasterCategory(w http.ResponseWriter, r *http.Request) {
	var item models.MasterCategory
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.DB.QueryRow(`INSERT INTO master_categories (name, description, is_active) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		item.Name, item.Description, true).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.IsActive = true
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateMasterCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var item models.MasterCategory
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.DB.Exec(`UPDATE master_categories SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`, item.Name, item.Description, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMasterCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.DB.Exec(`UPDATE master_categories SET is_active = false WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
