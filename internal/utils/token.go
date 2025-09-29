package utils

import (
	"net/http"
	"strings"
)

func ExtractToken(r *http.Request) string {
	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}

	authHeader := r.Header.Get("Authorization")
	if after, ok :=strings.CutPrefix(authHeader, "Bearer "); ok  {
		return after
	}

	return ""
}
