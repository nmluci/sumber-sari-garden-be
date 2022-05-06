package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/constant"
	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/responseutil"
)

type AuthHandler struct {
	r  *mux.Router
	as AuthService
}

func (ac *AuthHandler) InitHandler() {
	routes := ac.r.PathPrefix(constant.AUTH_API_PATH).Subrouter()
	routes.HandleFunc("/register", ac.RegisterNewUser()).Methods(http.MethodPost, http.MethodOptions)
	routes.HandleFunc("/login", ac.LoginUser()).Methods(http.MethodPost, http.MethodOptions)
}

func NewAuthHandler(r *mux.Router, as AuthService) *AuthHandler {
	return &AuthHandler{r: r, as: as}
}

func (auth *AuthHandler) RegisterNewUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &dto.UserRegistrationRequest{}

		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			log.Printf("[RegisterNewUser] failed to parsed JSON data, err => %+v", err)
			responseutil.WriteErrorResponse(w, err)
			return
		}

		err = auth.as.RegisterNewUser(r.Context(), data)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (auth *AuthHandler) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &dto.UserSignIn{}

		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			log.Printf("[LoginUser] failed to parsed JSON data, err => %+v", err)
			responseutil.WriteErrorResponse(w, err)
			return
		}

		res, err := auth.as.LoginUser(r.Context(), data)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, res)
	}
}
