package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/pkg/entity/response"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
}

func CorsMiddleware(whitelistedUrls map[string]bool) mux.MiddlewareFunc {
	return func (next http.Handler) http.Handler {
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
					case *response.BaseResponse:
						w.WriteHeader(v.Code)
						response.NewBaseResponse(
							v.Code,
							v.Message,
							v.Errors,
							v.Data,
						).ToJSON(w)
					case error:
						w.WriteHeader(500)
						response.NewBaseResponse(
							500,
							response.RESPONSE_ERROR_RUNTIME_MESSAGE,
							response.NewErrorResponseData(
								response.NewErrorResponseValue("msg", v.Error())),
							nil,
						).ToJSON(w)
					default:
						w.WriteHeader(500)
						response.NewBaseResponse(
							500,
							response.RESPONSE_ERROR_RUNTIME_MESSAGE,
							response.NewErrorResponseData(
								response.NewErrorResponseValue("msg", "runtime error")),
							nil,
						).ToJSON(w)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
}