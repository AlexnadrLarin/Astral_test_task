package routes

import (
	"github.com/gorilla/mux"

	handlers "docs_storage/internal/delivery/http/handlers"
)

func SetupDocsRoutes(r *mux.Router, docsHandler *handlers.DocsHandler) {
	r.HandleFunc("/api/docs", docsHandler.HandleUploadDoc).Methods("POST")
    r.HandleFunc("/api/docs", docsHandler.HandleListDocs).Methods("GET", "HEAD")
    r.HandleFunc("/api/docs/{id}", docsHandler.HandleGetDoc).Methods("GET", "HEAD")
    r.HandleFunc("/api/docs/{id}", docsHandler.HandleDeleteDoc).Methods("DELETE")
}
