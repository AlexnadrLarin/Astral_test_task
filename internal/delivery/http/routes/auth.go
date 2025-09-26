package routes

import (
	"github.com/gorilla/mux"

	handlers "docs_storage/internal/delivery/http/handlers"
)

func SetupAuthRoutes(r *mux.Router, docsHandler *handlers.AuthHandler) {
	r.HandleFunc("/api/register", docsHandler.Register).Methods("POST")
    r.HandleFunc("/api/auth", docsHandler.Auth).Methods("POST")
    r.HandleFunc("/api/auth/{token}", docsHandler.Logout).Methods("DELETE")
}
