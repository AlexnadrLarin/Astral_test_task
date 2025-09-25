package routes

import (
	"github.com/gorilla/mux"

	handlers "docs_storage/internal/delivery/http/handlers"
)

func SetupRoutes(r *mux.Router, docsHandler *handlers.DocsHanlder) {
	r.HandleFunc("/api/docs", docsHandler.HandleUploadDoc).Methods("POST")
    r.HandleFunc("/api/docs", docsHandler.HandleListDocs).Methods("GET", "HEAD")
    r.HandleFunc("/api/docs/{id}", docsHandler.HandleGetDoc).Methods("GET", "HEAD")
    r.HandleFunc("/api/docs/{id}", docsHandler.HandleDeleteDoc).Methods("DELETE")
}
