package product

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

type ProductHandler struct {
	r  *mux.Router
	ps ProductService
}

func (ps *ProductHandler) InitHandler() {
	routes := ps.r.PathPrefix(constant.INVENTORY_API_PATH).Subrouter()
	// Products
	routes.HandleFunc("/products", ps.GetAllCategory()).Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("/products", ps.StoreNewProduct()).Methods(http.MethodPost, http.MethodOptions)
	routes.HandleFunc("/products/{id}", ps.GetProductByID()).Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("/products/{id}", ps.UpdateProduct()).Methods(http.MethodPatch, http.MethodOptions)
	routes.HandleFunc("/products/{id}", ps.DeleteProduct()).Methods(http.MethodDelete, http.MethodOptions)

	// Product Categories
	routes.HandleFunc("/category", ps.GetAllCategory()).Methods(http.MethodGet, http.MethodOptions)
	routes.HandleFunc("/category", ps.StoreNewCategory()).Methods(http.MethodPost, http.MethodOptions)
	routes.HandleFunc("/category/{id}", ps.UpdateCategory()).Methods(http.MethodPatch, http.MethodOptions)
	routes.HandleFunc("/category/{id}", ps.DeleteCategory()).Methods(http.MethodDelete, http.MethodOptions)
}

func (prd *ProductHandler) GetAllProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		limit := params.Get("limit")
		if limit == "" {
			limit = "10"
		}

		limitParsed, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			log.Printf("[GetAllProduct] failed to parsed limit data, err => %+v\n", err)
			responseutil.WriteErrorResponse(w, err)
			return
		}

		offset := params.Get("limit")
		if limit == "" {
			limit = "0"
		}

		offsetParsed, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			log.Printf("[GetAllProduct] failed to parsed limit data, err => %+v\n", err)
			responseutil.WriteErrorResponse(w, err)
			return
		}

		res, err := prd.ps.GetAllProduct(r.Context(), limitParsed, offsetParsed)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
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
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[GetProductByID] failed to convert productID, err => %+v\n", err)
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		res, err := prd.ps.GetProductByID(r.Context(), pidParsed)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
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
			responseutil.WriteErrorResponse(w, err)
			return
		}

		err = prd.ps.StoreNewProduct(r.Context(), data)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
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
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[UpdateProduct] failed to convert productID, err => %+v\n", err)
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		data := &dto.UpdateProductRequest{}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[UpdateProduct] failed to parse JSON data, err => %+v", err)
			responseutil.WriteErrorResponse(w, err)
			return
		}

		data.ID = pidParsed
		err = prd.ps.UpdateProduct(r.Context(), data)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
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
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[DeleteProduct] failed to convert productID, err => %+v\n", err)
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		err = prd.ps.DeleteProduct(r.Context(), pidParsed)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}

func (prd *ProductHandler) GetAllCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := prd.ps.GetAllCategory(r.Context())
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
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
			responseutil.WriteErrorResponse(w, err)
			return
		}

		err = prd.ps.StoreNewCategory(r.Context(), data)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
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
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[UpdateCategory] failed to convert categoryID, err => %+v\n", err)
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		data := &dto.UpdateCategoryRequest{}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[UpdateCategory] failed to parse JSON data, err => %+v", err)
			responseutil.WriteErrorResponse(w, err)
			return
		}

		data.CategoryID = pidParsed
		err = prd.ps.UpdateCategory(r.Context(), data)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
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
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		pidParsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			log.Printf("[DeleteCategory] failed to convert categoryID, err => %+v\n", err)
			responseutil.WriteErrorResponse(w, errors.ErrUnknown)
			return
		}

		err = prd.ps.DeleteCategory(r.Context(), pidParsed)
		if err != nil {
			responseutil.WriteErrorResponse(w, err)
			return
		}

		responseutil.WriteSuccessResponse(w, http.StatusOK, nil)
	}
}
