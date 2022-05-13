package usercart

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/constant"
	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/responseutil"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type UsercartHandler struct {
	r  *mux.Router
	us UsercartService
}

func (handler *UsercartHandler) InitHandler() {
	routes := handler.r.PathPrefix(constant.USERCART_API_PATH).Subrouter()

	routes.HandleFunc("", handler.GetCart()).Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("", handler.UpsertItem()).Methods(http.MethodPost, http.MethodOptions)
	routes.HandleFunc("/:id", handler.RemoveItem()).Methods(http.MethodDelete, http.MethodOptions)
}

func NewUsercartHandler(r *mux.Router, us UsercartService) *UsercartHandler {
	return &UsercartHandler{r: r, us: us}
}

func (handler *UsercartHandler) UpsertItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &dto.UpsertItemRequest{}
		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			log.Printf("[UpsertItem] failed to parse JSON data, err => %+v", err)
			responseutil.WriteErrorResponse(w, err)
			return
		}

		err = handler.us.UpsertItem(r.Context(), data)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (handler *UsercartHandler) RemoveItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productID, ok := vars["id"]
		if !ok {
			log.Printf("[RemoveItem] invalid product_id\n")
			responseutil.WriteErrorResponse(w, errors.ErrInvalidRequestBody)
			return
		}

		parsedID, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[RemoveItem] failed to parse productID, err => %+v\n\n", err)
			responseutil.WriteErrorResponse(w, errors.ErrInvalidRequestBody)
			return
		}

		err = handler.us.RemoveItem(r.Context(), parsedID)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (handler *UsercartHandler) GetCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := handler.us.GetCart(r.Context())
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, res)
	}
}
