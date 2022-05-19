package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/auth"
	"github.com/nmluci/sumber-sari-garden/internal/entity"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/responseutil"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

func AuthMiddleware(service auth.AuthService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[AccessControl-Middleware] Incoming Request to %v\n", r.URL)

			token := r.Header.Get("Authorization")
			if token == "" {
				responseutil.WriteErrorResponse(w, errors.ErrUserPriv)
				return
			}

			splittedToken := strings.Split(token, " ")
			if len(splittedToken) != 2 {
				responseutil.WriteErrorResponse(w, errors.ErrInvalidRequestBody)
				return
			}

			if splittedToken[0] != "Bearer" {
				responseutil.WriteErrorResponse(w, errors.ErrInvalidRequestBody)
				return
			}

			accessToken := splittedToken[1]
			user, err := service.FindUserByAccessToken(r.Context(), accessToken)
			if err != nil {
				responseutil.WriteErrorResponse(w, err)
				return
			} else if user == nil {
				responseutil.WriteErrorResponse(w, errors.ErrUserCredential)
				return
			}

			var auth entity.AuthCtxKey = "user-ctx"
			authCtx := context.WithValue(r.Context(), auth, &entity.UserContext{
				UserID: user.UserID,
				Priv:   user.UserRole,
			})
			r = r.WithContext(authCtx)

			next.ServeHTTP(w, r)
		})
	}
}
