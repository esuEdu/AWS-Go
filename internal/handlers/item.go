package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/esuEdu/AWS-Go/internal/models"
)

type Handler struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Handler { return &Handler{DB: db} }

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), `SELECT id, name, description, created_at, updated_at FROM items ORDER BY id`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := make([]models.Item, 0)
	for rows.Next() {
		var it models.Item
		if err := rows.Scan(&it.ID, &it.Name, &it.Description, &it.CreatedAt, &it.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, it)
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}
	var it models.Item
	err := h.DB.QueryRow(r.Context(), `SELECT id, name, description, created_at, updated_at FROM items WHERE id=$1`, id).
		Scan(&it.ID, &it.Name, &it.Description, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if in.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	var it models.Item
	err := h.DB.QueryRow(r.Context(), `INSERT INTO items(name, description) VALUES ($1,$2) RETURNING id, name, description, created_at, updated_at`, in.Name, in.Description).
		Scan(&it.ID, &it.Name, &it.Description, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, it)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}
	var in struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var it models.Item
	err := h.DB.QueryRow(r.Context(), `UPDATE items SET name=COALESCE($1,name), description=COALESCE($2,description) WHERE id=$3 RETURNING id, name, description, created_at, updated_at`, in.Name, in.Description, id).
		Scan(&it.ID, &it.Name, &it.Description, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}
	ct, err := h.DB.Exec(r.Context(), `DELETE FROM items WHERE id=$1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ct.RowsAffected() == 0 {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseIDParam(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
