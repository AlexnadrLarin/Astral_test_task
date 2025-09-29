package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	models "docs_storage/internal/models"
	utils "docs_storage/internal/utils"
	"docs_storage/pkg/logger"
)

type docsService interface {
	Create(ctx context.Context, meta *models.Document, fileName string, fileData []byte, jsonData []byte, token string) (*models.Document, error)
	List(ctx context.Context, token, login, key, value string, limit int) ([]models.Document, error)
	GetByID(ctx context.Context, id, token string) (*models.Document, error)
	Delete(ctx context.Context, id, token string) error
}

type DocsHandler struct {
	svc    docsService
	logger *logger.Logger
}

func NewDocsHandler(svc docsService, log *logger.Logger) *DocsHandler {
	return &DocsHandler{svc: svc, logger: log}
}

func (h *DocsHandler) HandleUploadDoc(w http.ResponseWriter, r *http.Request) {
	token := utils.ExtractToken(r)
	if token == "" {
		h.logger.Error.Print("upload attempt without token")
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("token required"))
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Error.Printf("failed to parse multipart form: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("cannot parse form"))
		return
	}

	metaPart := r.FormValue("meta")
	if metaPart == "" {
		h.logger.Error.Print("upload attempt without meta")
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("meta is required"))
		return
	}

	var meta models.Document
	if err := json.Unmarshal([]byte(metaPart), &meta); err != nil {
		h.logger.Error.Printf("invalid meta json: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("invalid meta json"))
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
			h.logger.Error.Print("file is required but missing")
			utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("file is required"))
			return
		}
		defer file.Close()

		fileName = header.Filename
		fileData, _ = io.ReadAll(file)
	}

	doc, err := h.svc.Create(r.Context(), &meta, fileName, fileData, jsonData, token)
	if err != nil {
		h.logger.Error.Printf("failed to create document: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResp("cannot create document"))
		return
	}

	h.logger.Info.Printf("document uploaded: %s by token %s", doc.ID, token)
	utils.WriteJSON(w, http.StatusOK, utils.DocDetail(*doc))
}

func (h *DocsHandler) HandleListDocs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := utils.ExtractToken(r)
	if token == "" {
		h.logger.Error.Print("list attempt without token")
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("token required"))
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
		h.logger.Error.Printf("failed to list documents: %v", err)
		utils.WriteJSON(w, http.StatusForbidden, utils.ErrorResp(err.Error()))
		return
	}

	h.logger.Info.Printf("documents listed by token %s, count: %d", token, len(docs))
	utils.WriteJSON(w, http.StatusOK, utils.DocsList(docs))
}

func (h *DocsHandler) HandleGetDoc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	token := utils.ExtractToken(r)
	if token == "" {
		h.logger.Error.Print("get document attempt without token")
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("token required"))
		return
	}

	doc, err := h.svc.GetByID(ctx, id, token)
	if err != nil {
		h.logger.Error.Printf("failed to get document %s: %v", id, err)
		utils.WriteJSON(w, http.StatusForbidden, utils.ErrorResp(err.Error()))
		return
	}

	h.logger.Info.Printf("document retrieved: %s by token %s", id, token)
	if doc.File {
		w.Header().Set("Content-Type", doc.Mime)
		http.ServeFile(w, r, doc.FilePath)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.DocDetail(*doc))
}

func (h *DocsHandler) HandleDeleteDoc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	token := utils.ExtractToken(r)
	if token == "" {
		h.logger.Error.Print("delete attempt without token")
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("token required"))
		return
	}

	if err := h.svc.Delete(ctx, id, token); err != nil {
		h.logger.Error.Printf("failed to delete document %s: %v", id, err)
		utils.WriteJSON(w, http.StatusForbidden, utils.ErrorResp(err.Error()))
		return
	}

	h.logger.Info.Printf("document deleted: %s by token %s", id, token)
	utils.WriteJSON(w, http.StatusOK, utils.DeleteResp(id))
}
