package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/responseutil"
	"github.com/nmluci/sumber-sari-garden/pkg/dto"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
}

func CorsMiddleware(whitelistedUrls map[string]bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE, PATCH")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				requestOriginUrl := r.Header.Get("Origin")
				log.Printf("INFO CorsMiddleware: received request from %s %v", requestOriginUrl, whitelistedUrls[requestOriginUrl])
				if whitelistedUrls[requestOriginUrl] {
					w.Header().Set("Access-Control-Allow-Origin", requestOriginUrl)
				}

				if r.Method != http.MethodOptions {
					next.ServeHTTP(w, r)
					return
				}

				w.Write([]byte("safe flight packet"))
			})
	}
}

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					switch v := r.(type) {
					case *dto.BaseResponse:
						responseutil.BaseResponseWriter(
							w,
							500,
							v.Data,
							v.Error,
						)
					case error:
						responseutil.WriteErrorResponse(w, v)
					default:
						responseutil.WriteErrorResponse(w, errors.ErrUnknown)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
}
