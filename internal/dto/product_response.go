package dto

import (
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/entity"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type ProductResponse struct {
	ID           uint64 `json:"id"`
	CategoryID   uint64 `json:"category_id"`
	Name         string `json:"product_name"`
	CategoryName string `json:"category_name"`
	Qty          uint64 `json:"qty"`
	Price        uint64 `json:"price"`
	PictureURL   string `json:"picture_url"`
	Description  string `json:"description"`
}

type ProductsResponse struct {
	Products []*ProductResponse `json:"products"`
	Limit    int64              `json:"limit"`
	Offset   int64              `json:"offset"`
	Total    int64              `json:"total"`
}

type CategoryResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type CategoriesResponse []*CategoryResponse

func NewProductResponse(product *entity.Product, cat *entity.ProductCategory) (res *ProductResponse, err error) {
	if product == nil || cat == nil {
		log.Printf("[NewProductResponse] failed to encode response data due to incomplete data")
		err = errors.ErrInvalidResources
		return
	}

	res = &ProductResponse{
		ID:           product.ID,
		CategoryID:   cat.ID,
		Name:         product.Name,
		CategoryName: cat.Name,
		Qty:          product.Qty,
		Price:        product.Price,
		PictureURL:   product.PictureURL,
		Description:  product.Description,
	}

	return
}

func NewProductsResponse(items entity.ProductDetails, limit int64, offset int64, total int64) (res *ProductsResponse, err error) {
	if items == nil {
		log.Printf("[NewProductsResponse] failed to encode response data due to incomplete data")
		err = errors.ErrInvalidResources
		return
	}

	res = &ProductsResponse{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}

	for _, itm := range items {
		temp := &ProductResponse{
			ID:           itm.ID,
			CategoryID:   itm.CategoryID,
			Name:         itm.Name,
			CategoryName: itm.CategoryName,
			Qty:          itm.Qty,
			Price:        itm.Price,
			PictureURL:   itm.PictureURL,
			Description:  itm.Description,
		}
		res.Products = append(res.Products, temp)
	}

	return
}
