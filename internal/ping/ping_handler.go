package ping

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/constant"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/responseutil"
)

type PingHandler struct {
	r  *mux.Router
	ps PingService
}

func (pc *PingHandler) InitHandler() {
	routes := pc.r.PathPrefix(constant.PING_API_PATH).Subrouter()
	routes.HandleFunc("", pc.PingHandler).Methods(http.MethodGet)
}

func NewPingHandler(r *mux.Router, ps PingService) *PingHandler {
	return &PingHandler{r: r, ps: ps}
}

func (pc PingHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	res := pc.ps.Ping(r.Context())
	responseutil.WriteSuccessResponse(w, http.StatusOK, res)
}
