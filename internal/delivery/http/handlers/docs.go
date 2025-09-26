package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"context"

	"github.com/gorilla/mux"

	models "docs_storage/internal/models"
	"docs_storage/internal/utils"
)

type docsService interface {
	Create(ctx context.Context, meta *models.Document, fileName string, fileData []byte, jsonData []byte, token string) (*models.Document, error)
	List(ctx context.Context, token, login, key, value string, limit int) ([]models.Document, error)
	GetByID(ctx context.Context, id, token string) (*models.Document, error)
	Delete(ctx context.Context, id, token string) error
}

type DocsHandler struct {
	svc docsService
}

func NewDocsHandler(svc docsService) *DocsHandler {
	return &DocsHandler{svc: svc}
}

func (h *DocsHandler) HandleUploadDoc(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "cannot parse form", http.StatusBadRequest)
		return
	}

	metaPart := r.FormValue("meta")
	if metaPart == "" {
		http.Error(w, "meta is required", http.StatusBadRequest)
		return
	}

	var meta models.Document
	if err := json.Unmarshal([]byte(metaPart), &meta); err != nil {
		http.Error(w, "invalid meta json", http.StatusBadRequest)
		return
	}

	var jsonData []byte
	if jsonPart := r.FormValue("json"); jsonPart != "" {
		jsonData = []byte(jsonPart)
	}

	var fileName string
	var fileData []byte
	if meta.File {
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileName = header.Filename
		fileData, _ = io.ReadAll(file)
	}

	doc, err := h.svc.Create(r.Context(), &meta, fileName, fileData, jsonData, token)
	if err != nil {
		http.Error(w, "cannot create document", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.ToDocDetailResponse(*doc))
}

func (h *DocsHandler) HandleListDocs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	login := r.URL.Query().Get("login")
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	limit := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil {
			limit = n
		}
	}

	docs, err := h.svc.List(ctx, token, login, key, value, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.ToDocsListResponse(docs))
}

func (h *DocsHandler) HandleGetDoc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	doc, err := h.svc.GetByID(ctx, id, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if doc.File {
		w.Header().Set("Content-Type", doc.Mime)
		http.ServeFile(w, r, doc.FilePath)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.ToDocDetailResponse(*doc))
}

func (h *DocsHandler) HandleDeleteDoc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(ctx, id, token); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.ToDeleteResponse(id))
}