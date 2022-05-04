package auth

import (
	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/constant"
)

type AuthHandler struct {
	r  *mux.Router
	as AuthService
}

func (ac *AuthHandler) InitHandler() {
	_ = ac.r.PathPrefix(constant.AUTH_API_PATH).Subrouter()
	// routes.HandleFunc)
}

func NewAuthHandler(r *mux.Router, as AuthService) *AuthHandler {
	return &AuthHandler{r: r, as: as}
}
