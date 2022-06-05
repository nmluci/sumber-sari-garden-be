package usercart

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/constant"
	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/responseutil"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/timeutil"
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
	routes.HandleFunc("/{id}", handler.RemoveItem()).Methods(http.MethodDelete, http.MethodOptions)
	routes.HandleFunc("/checkout", handler.Checkout()).Methods(http.MethodPost, http.MethodOptions)
	routes.HandleFunc("/verify", handler.GetUnpaidOrder()).Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("/verify/{id}", handler.VerifyOrder()).Methods(http.MethodPatch, http.MethodOptions)
	routes.HandleFunc("/history", handler.OrderHistory()).Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("/history/all", handler.OrderHistoryAll()).Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("/history/{id}", handler.SpecificOrderHistoryById()).Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("/statistics", handler.GetStatistics()).Methods(http.MethodGet, http.MethodOptions)
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
			panic(err)
		}

		err = handler.us.UpsertItem(r.Context(), data)
		if err != nil {
			panic(err)
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
			panic(errors.ErrInvalidRequestBody)
		}

		parsedID, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[RemoveItem] failed to parse productID, err => %+v\n\n", err)
			panic(errors.ErrInvalidRequestBody)
		}

		err = handler.us.RemoveItem(r.Context(), parsedID)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (handler *UsercartHandler) GetCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := handler.us.GetCart(r.Context())
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, res)
	}
}

func (handler *UsercartHandler) Checkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &dto.OrderCheckoutRequest{}
		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			log.Printf("[Checkout] failed to parse JSON data, err => %+v\n", err)
			panic(err)
		}

		err = handler.us.Checkout(r.Context(), data)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (handler *UsercartHandler) SpecificOrderHistoryById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productID, ok := vars["id"]
		if !ok {
			log.Printf("[SpecificOrderHistoryById] invalid productID\n")
			panic(errors.ErrInvalidRequestBody)
		}

		parsedProductId, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[SpecificOrderHistoryById] error while parsed\n")
			panic(errors.ErrInvalidRequestBody)
		}

		data, err := handler.us.SpecificOrderHistoryById(r.Context(), parsedProductId)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, data)
	}
}

func (handler *UsercartHandler) OrderHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		limit := query.Get("limit")
		limitParsed, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			log.Printf("[OrderHistory] failed to parse limit data, err => %+v\n", err)
			if limit == "" {
				limitParsed = 10
			}
		}

		offset := query.Get("offset")
		offsetParsed, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			log.Printf("[OrderHistory] failed to parse offset data, err => %+v\n", err)
			if offset == "" {
				offsetParsed = 0
			}
		}

		dStart := query.Get("dateStart")
		dateStart, err := timeutil.ParseLocalTime(dStart, "2006-01-02")
		if err != nil {
			log.Printf("[OrderHistory] failed to parse dateStart data, err => %+v\n", err)
			if dStart == "" {
				dateStart = timeutil.Localize(time.Now().Add(-(24 * 31 * time.Hour)))
			}
		}

		dEnd := query.Get("dateEnd")
		dateEnd, err := timeutil.ParseLocalTime(dEnd, "2006-01-02")
		if err != nil {
			log.Printf("[OrderHistory] failed to parse dateEnd data, err => %+v\n", err)
			if dEnd == "" {
				dateEnd = timeutil.Localize(time.Now())
			}
		}

		if dateEnd.Before(dateStart) {
			log.Printf("[OrderHistory] invalid date range\n")
			panic(errors.ErrInvalidRequestBody)
		}

		params := dto.HistoryParams{
			Limit:     limitParsed,
			Offset:    offsetParsed,
			DateStart: dateStart,
			DateEnd:   dateEnd,
		}

		data, err := handler.us.OrderHistory(r.Context(), params)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, data)
	}
}

func (handler *UsercartHandler) VerifyOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		order := mux.Vars(r)["id"]
		orderParsed, err := strconv.ParseUint(order, 10, 64)
		if err != nil {
			log.Printf("[VerifyOrder] failed to parsed orderID, err => %+v\n", err)
			panic(err)
		}

		err = handler.us.VerifyOrder(r.Context(), orderParsed)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (handler *UsercartHandler) GetUnpaidOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := handler.us.GetUnpaidOrder(r.Context())
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, data)
	}
}

func (handler *UsercartHandler) OrderHistoryAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		limit := query.Get("limit")
		limitParsed, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			log.Printf("[OrderHistory] failed to parse limit data, err => %+v\n", err)
			if limit == "" {
				limitParsed = 10
			}
		}

		offset := query.Get("offset")
		offsetParsed, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			log.Printf("[OrderHistory] failed to parse offset data, err => %+v\n", err)
			if offset == "" {
				offsetParsed = 0
			}
		}

		dStart := query.Get("dateStart")
		dateStart, err := timeutil.ParseLocalTime(dStart, "2006-01-02 15:04:05")
		if err != nil {
			log.Printf("[OrderHistory] failed to parse dateStart data, err => %+v\n", err)
			if dStart == "" {
				dateStart = timeutil.Localize(time.Now().Add(-(24 * 31 * time.Hour)))
			}
		}

		dEnd := query.Get("dateEnd")
		dateEnd, err := timeutil.ParseLocalTime(dEnd, "2006-01-02 15:04:05")
		if err != nil {
			log.Printf("[OrderHistory] failed to parse dateEnd data, err => %+v\n", err)
			if dEnd == "" {
				dateEnd = timeutil.Localize(time.Now())
			}
		}

		if dateEnd.Before(dateStart) {
			log.Printf("[OrderHistory] invalid date range\n")
			panic(errors.ErrInvalidRequestBody)
		}

		params := dto.HistoryParams{
			Limit:     limitParsed,
			Offset:    offsetParsed,
			DateStart: dateStart,
			DateEnd:   dateEnd,
		}

		data, err := handler.us.OrderHistoryAll(r.Context(), params)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, data)
	}
}

func (handler *UsercartHandler) GetStatistics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := handler.us.GetStatistics(r.Context())
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, data)
	}
}
