package product

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nmluci/sumber-sari-garden/internal/constant"
	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/responseutil"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type ProductHandler struct {
	r  *mux.Router
	p  *mux.Router
	ps ProductService
}

func (ps *ProductHandler) InitHandler() {
	routes := ps.r.PathPrefix(constant.INVENTORY_API_PATH).Subrouter()
	protected := ps.p.PathPrefix(constant.INVENTORY_API_PATH).Subrouter()
	// Products
	routes.HandleFunc("/products", ps.GetAllProduct()).Methods(http.MethodPut, http.MethodOptions)
	routes.HandleFunc("/products/{id}", ps.GetProductByID()).Methods(http.MethodGet, http.MethodOptions)
	protected.HandleFunc("/products", ps.StoreNewProduct()).Methods(http.MethodPost, http.MethodOptions)
	protected.HandleFunc("/products/{id}", ps.UpdateProduct()).Methods(http.MethodPatch, http.MethodOptions)
	protected.HandleFunc("/products/{id}", ps.DeleteProduct()).Methods(http.MethodDelete, http.MethodOptions)

	// Product Categories
	routes.HandleFunc("/category", ps.GetAllCategory()).Methods(http.MethodGet, http.MethodOptions)
	protected.HandleFunc("/category", ps.StoreNewCategory()).Methods(http.MethodPost, http.MethodOptions)
	protected.HandleFunc("/category/{id}", ps.UpdateCategory()).Methods(http.MethodPatch, http.MethodOptions)
	protected.HandleFunc("/category/{id}", ps.DeleteCategory()).Methods(http.MethodDelete, http.MethodOptions)

	// Coupon
	routes.HandleFunc("/coupons", ps.GetActiveCoupons()).Methods(http.MethodGet, http.MethodOptions)
	protected.HandleFunc("/coupons/all", ps.GetAllCoupon()).Methods(http.MethodGet, http.MethodOptions)
	protected.HandleFunc("/coupons", ps.StoreNewCoupon()).Methods(http.MethodPost, http.MethodOptions)
	protected.HandleFunc("/coupons/{id}", ps.UpdateCoupon()).Methods(http.MethodPatch, http.MethodOptions)
	protected.HandleFunc("/coupons/{id}", ps.DeleteCoupon()).Methods(http.MethodDelete, http.MethodOptions)
}

func NewProductHandler(r *mux.Router, p *mux.Router, ps ProductService) *ProductHandler {
	return &ProductHandler{r: r, p: p, ps: ps}
}

func (prd *ProductHandler) GetAllProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &dto.ProductSearchParams{}
		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil && err != io.EOF {
			log.Printf("[GetAllProduct] failed to parse JSON data, err => %+v", err)
			panic(err)
		}

		res, err := prd.ps.GetAllProduct(r.Context(), data)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, res)
	}
}

func (prd *ProductHandler) GetProductByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		productID, ok := routeVar["id"]
		if !ok {
			log.Printf("[GetProductByID] failed to parsed productID data\n")
			panic(errors.ErrInvalidRequestBody)
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[GetProductByID] failed to convert productID, err => %+v\n", err)
			panic(errors.ErrInvalidRequestBody)
		}

		res, err := prd.ps.GetProductByID(r.Context(), pidParsed)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, res)
	}
}

func (prd *ProductHandler) StoreNewProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &dto.NewProductRequest{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[StoreNewProduct] failed to parse JSON data, err => %+v", err)
			panic(err)
		}

		err = prd.ps.StoreNewProduct(r.Context(), data)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (prd *ProductHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		productID, ok := routeVar["id"]
		if !ok {
			log.Printf("[UpdateProduct] failed to parsed productID data\n")
			panic(errors.ErrInvalidRequestBody)
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[UpdateProduct] failed to convert productID, err => %+v\n", err)
			panic(errors.ErrInvalidRequestBody)
		}

		data := &dto.UpdateProductRequest{}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[UpdateProduct] failed to parse JSON data, err => %+v", err)
			panic(errors.ErrInvalidRequestBody)
		}

		data.ID = pidParsed
		err = prd.ps.UpdateProduct(r.Context(), data)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (prd *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		productID, ok := routeVar["id"]
		if !ok {
			log.Printf("[DeleteProduct] failed to parsed productID data\n")
			panic(errors.ErrInvalidRequestBody)
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[DeleteProduct] failed to convert productID, err => %+v\n", err)
			panic(errors.ErrInvalidRequestBody)
		}

		err = prd.ps.DeleteProduct(r.Context(), pidParsed)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (prd *ProductHandler) GetAllCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := prd.ps.GetAllCategory(r.Context())
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, res)
	}
}

func (prd *ProductHandler) StoreNewCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &dto.NewCategoryRequest{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[StoreNewCategory] failed to parse JSON data, err => %+v", err)
			panic(err)
		}

		err = prd.ps.StoreNewCategory(r.Context(), data)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (prd *ProductHandler) UpdateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		productID, ok := routeVar["id"]
		if !ok {
			log.Printf("[UpdateCategory] failed to parsed categoryID data\n")
			panic(errors.ErrInvalidRequestBody)
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[UpdateCategory] failed to convert categoryID, err => %+v\n", err)
			panic(errors.ErrInvalidRequestBody)
		}

		data := &dto.UpdateCategoryRequest{}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[UpdateCategory] failed to parse JSON data, err => %+v", err)
			panic(errors.ErrInvalidRequestBody)
		}

		data.CategoryID = pidParsed
		err = prd.ps.UpdateCategory(r.Context(), data)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (prd *ProductHandler) DeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		productID, ok := routeVar["id"]
		if !ok {
			log.Printf("[DeleteCategory] failed to parsed categoryID data\n")
			panic(errors.ErrInvalidRequestBody)
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[DeleteCategory] failed to convert categoryID, err => %+v\n", err)
			panic(errors.ErrInvalidRequestBody)
		}

		err = prd.ps.DeleteCategory(r.Context(), pidParsed)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (handler *ProductHandler) GetActiveCoupons() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		limit := query.Get("limit")
		limitParsed, err := strconv.ParseInt(limit, 10, 64)
		if err != nil && limit != "" {
			panic(errors.ErrInvalidRequestBody)
		} else if limit == "" {
			limitParsed = 10
		}

		offset := query.Get("offset")
		offsetParsed, err := strconv.ParseInt(offset, 10, 64)
		if err != nil && offset != "" {
			panic(errors.ErrInvalidRequestBody)
		} else if offset == "" {
			offsetParsed = 0
		}

		data, err := handler.ps.GetActiveCoupons(r.Context(), limitParsed, offsetParsed)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, data)
	}
}

func (handler *ProductHandler) GetAllCoupon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		limit := query.Get("limit")
		limitParsed, err := strconv.ParseInt(limit, 10, 64)
		if err != nil && limit != "" {
			panic(errors.ErrInvalidRequestBody)
		} else if limit == "" {
			limitParsed = 10
		}

		offset := query.Get("offset")
		offsetParsed, err := strconv.ParseInt(offset, 10, 64)
		if err != nil && offset != "" {
			panic(errors.ErrInvalidRequestBody)
		} else if offset == "" {
			offsetParsed = 0
		}

		data, err := handler.ps.GetAllCoupon(r.Context(), limitParsed, offsetParsed)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, data)
	}
}

func (prd *ProductHandler) StoreNewCoupon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &dto.Coupon{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[StoreNewCoupon] failed to parse JSON data, err => %+v", err)
			panic(err)
		}

		err = prd.ps.StoreNewCoupon(r.Context(), data)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (prd *ProductHandler) UpdateCoupon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		couponID, ok := routeVar["id"]
		if !ok {
			log.Printf("[UpdateCoupon] failed to parsed couponID data\n")
			panic(errors.ErrInvalidRequestBody)
		}

		cidParsed, err := strconv.ParseInt(couponID, 10, 64)
		if err != nil {
			log.Printf("[UpdateCoupon] failed to convert couponID, err => %+v\n", err)
			panic(errors.ErrInvalidRequestBody)
		}

		data := &dto.Coupon{}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[UpdateCoupon] failed to parse JSON data, err => %+v", err)
			panic(errors.ErrInvalidRequestBody)
		}

		err = prd.ps.UpdateCoupon(r.Context(), cidParsed, data)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (prd *ProductHandler) DeleteCoupon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		couponID, ok := routeVar["id"]
		if !ok {
			log.Printf("[DeleteCoupon] failed to parsed couponID data\n")
			panic(errors.ErrInvalidRequestBody)
		}

		cidParsed, err := strconv.ParseInt(couponID, 10, 64)
		if err != nil {
			log.Printf("[DeleteCoupon] failed to convert couponID, err => %+v\n", err)
			panic(errors.ErrInvalidRequestBody)
		}

		err = prd.ps.DeleteCoupon(r.Context(), cidParsed)
		if err != nil {
			panic(err)
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}
