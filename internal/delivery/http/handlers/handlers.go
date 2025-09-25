package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	models "docs_storage/internal/models"
)
type docsService interface {
	Create(ctx context.Context, meta *models.Document, fileName string, fileData []byte, jsonData []byte) (*models.Document, error)
}

type DocsHanlder struct {
	svc docsService
}

func NewDocsHandler(svc docsService) *DocsHanlder {
	return &DocsHanlder{svc: svc}
}

func (h *DocsHanlder) HandleUploadDoc(w http.ResponseWriter, r *http.Request)  {
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

	doc, err := h.svc.Create(r.Context(), &meta, fileName, fileData, jsonData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "cannot create document", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"data": map[string]interface{}{
			"json": json.RawMessage(doc.JSONData),
			"file": doc.Name,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *DocsHanlder) HandleListDocs(w http.ResponseWriter, r *http.Request) {

}

func (h *DocsHanlder) HandleGetDoc(w http.ResponseWriter, r *http.Request) {

}

func (h *DocsHanlder) HandleDeleteDoc(w http.ResponseWriter, r *http.Request) {

}
