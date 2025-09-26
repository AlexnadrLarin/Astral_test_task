package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"docs_storage/internal/service"
)

type AuthHandler struct {
    svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
    return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Token string `json:"token"`
        Login string `json:"login"`
        Pswd  string `json:"pswd"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.svc.Register(r.Context(), input.Token, input.Login, input.Pswd); err != nil {
        http.Error(w, err.Error(), http.StatusForbidden)
        return
    }

    json.NewEncoder(w).Encode(map[string]any{
        "response": map[string]string{"login": input.Login},
    })
}

func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var input struct {
        Login string `json:"login"`
        Pswd  string `json:"pswd"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    token, err := h.svc.Auth(r.Context(), input.Login, input.Pswd)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(map[string]any{
        "response": map[string]string{"token": token},
    })
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    token := strings.TrimPrefix(r.URL.Path, "/api/auth/")
    if token == "" {
        http.Error(w, "missing token", http.StatusBadRequest)
        return
    }

    if err := h.svc.Logout(r.Context(), token); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]any{
        "response": map[string]bool{token: true},
    })
}
