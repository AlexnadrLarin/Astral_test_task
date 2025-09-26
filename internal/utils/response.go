package utils

import (
	"encoding/json"
	"net/http"

	"docs_storage/internal/models"
)

type DocsListResponse struct {
	Data struct {
		Docs []DocResponse `json:"docs"`
	} `json:"data"`
}

type DocDetailResponse struct {
	Data DocResponse `json:"data"`
}

type DeleteResponse struct {
	Response map[string]bool `json:"response"`
}

type DocResponse struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Mime    string          `json:"mime"`
	File    bool            `json:"file"`
	Public  bool            `json:"public"`
	Grant   []string        `json:"grant"`
	Created string          `json:"created"`
	JSON    json.RawMessage `json:"json_data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func ToDocResponse(d models.Document, withJSON bool) DocResponse {
	resp := DocResponse{
		ID:      d.ID,
		Name:    d.Name,
		Mime:    d.Mime,
		File:    d.File,
		Public:  d.Public,
		Grant:   d.Grant,
		Created: d.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if withJSON && len(d.JSONData) > 0 {
		resp.JSON = json.RawMessage(d.JSONData)
	}

	return resp
}

func ToDocsListResponse(docs []models.Document) DocsListResponse {
	resp := DocsListResponse{}
	resp.Data.Docs = make([]DocResponse, 0, len(docs))
	for _, d := range docs {
		resp.Data.Docs = append(resp.Data.Docs, ToDocResponse(d, false))
	}
	return resp
}

func ToDocDetailResponse(d models.Document) DocDetailResponse {
	return DocDetailResponse{
		Data: ToDocResponse(d, true),
	}
}

func ToDeleteResponse(id string) DeleteResponse {
	return DeleteResponse{
		Response: map[string]bool{id: true},
	}
}
