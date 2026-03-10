package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-backend-blueprint/internal/entity"
	"go-backend-blueprint/internal/store"
)

// ItemsHandler menangani CRUD untuk resource Item (GET/POST /items, GET/PUT/DELETE /items/:id).
type ItemsHandler struct {
	Store store.ItemStore
}

// createItemRequest body untuk POST /items.
type createItemRequest struct {
	Name string `json:"name"`
}

// updateItemRequest body untuk PUT /items/:id.
type updateItemRequest struct {
	Name string `json:"name"`
}

// HandleItems menangani GET (list) dan POST (create) /items.
func (h *ItemsHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/items" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleItemByID menangani GET, PUT, DELETE /items/:id.
func (h *ItemsHandler) HandleItemByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ItemsHandler) list(w http.ResponseWriter, _ *http.Request) {
	items, err := h.Store.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if items == nil {
		items = []*entity.Item{}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}

func (h *ItemsHandler) getByID(w http.ResponseWriter, _ *http.Request, id string) {
	item, err := h.Store.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if item == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *ItemsHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := h.Store.Create(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *ItemsHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req updateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := h.Store.Update(id, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if item == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *ItemsHandler) delete(w http.ResponseWriter, _ *http.Request, id string) {
	err := h.Store.Delete(id)
	if err != nil {
		if err == store.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
