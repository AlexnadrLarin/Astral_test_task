package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"docs_storage/internal/utils"
	"docs_storage/pkg/logger"
)

type authService interface {
	Register(ctx context.Context, token, login, pswd string) error
	Auth(ctx context.Context, login, pswd string) (string, error)
	Logout(ctx context.Context, token string) error
}

type AuthHandler struct {
	svc    authService
	logger *logger.Logger
}

func NewAuthHandler(svc authService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		svc:    svc,
		logger: log,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token string `json:"token"`
		Login string `json:"login"`
		Pswd  string `json:"pswd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error.Printf("failed to decode register input: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp(err.Error()))
		return
	}

	if err := h.svc.Register(r.Context(), input.Token, input.Login, input.Pswd); err != nil {
		h.logger.Error.Printf("register failed for login %s: %v", input.Login, err)
		utils.WriteJSON(w, http.StatusForbidden, utils.ErrorResp(err.Error()))
		return
	}

	h.logger.Info.Printf("user registered: %s", input.Login)
	utils.WriteJSON(w, http.StatusOK, utils.RegisterResp(input.Login))
}

func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Login string `json:"login"`
		Pswd  string `json:"pswd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error.Printf("failed to decode auth input: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp(err.Error()))
		return
	}

	token, err := h.svc.Auth(r.Context(), input.Login, input.Pswd)
	if err != nil {
		h.logger.Error.Printf("auth failed for login %s: %v", input.Login, err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.ErrorResp(err.Error()))
		return
	}

	h.logger.Info.Printf("user authenticated: %s", input.Login)
	utils.WriteJSON(w, http.StatusOK, utils.AuthResp(token))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/api/auth/")
	if token == "" {
		h.logger.Error.Print("logout attempt with missing token")
		utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResp("missing token"))
		return
	}

	if err := h.svc.Logout(r.Context(), token); err != nil {
		h.logger.Error.Printf("logout failed for token %s: %v", token, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResp(err.Error()))
		return
	}

	h.logger.Info.Printf("user logged out, token: %s", token)
	utils.WriteJSON(w, http.StatusOK, utils.LogoutResp(token))
}
