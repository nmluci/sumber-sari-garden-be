package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/constant"
	"github.com/nmluci/sumber-sari-garden/internal/ping/service/api"
	"github.com/nmluci/sumber-sari-garden/pkg/entity/response"
)

type PingController struct {
	r *mux.Router
	ps api.PingService
}

func (pc PingController) PingHandler(w http.ResponseWriter, r *http.Request) {
	res := pc.ps.Ping(r.Context())
	response.NewBaseResponse(
		200,
		response.RESPONSE_SUCCESS_MESSAGE,
		nil,
		res).SendResponse(&w)
}

func (pc *PingController) InitController() {
	routes := pc.r.PathPrefix(constant.PING_API_PATH).Subrouter()
	routes.HandleFunc("", pc.PingHandler).Methods(http.MethodGet)
}

func ProvideMsibController(r *mux.Router, ps api.PingService) *PingController {
	return &PingController{
		r: r,
		ps: ps,
	}
}