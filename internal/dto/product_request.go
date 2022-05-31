package dto

import (
	"strconv"

	"github.com/nmluci/sumber-sari-garden/internal/models"
)

type ProductSearchParams struct {
	Keyword    string   `json:"keyword"`
	Categories []uint64 `json:"categories"`
	Order      uint64   `json:"order_type"`
	Limit      uint64   `json:"limit"`
	Offset     uint64   `json:"offset"`
}

type NewProductRequest struct {
	Name        string `json:"name"`
	PictureURL  string `json:"picture_url"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Qty         uint64 `json:"qty"`
	CategoryID  uint64 `json:"category_id"`
}

type NewCategoryRequest struct {
	Name string `json:"category_name"`
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	PictureURL  string `json:"picture_url"`
	Description string `json:"description"`
	ID          uint64
	Price       uint64 `json:"price"`
	Qty         uint64 `json:"qty"`
	CategoryID  uint64 `json:"category_id"`
}

type UpdateCategoryRequest struct {
	Name       string `json:"category_name"`
	CategoryID uint64
}

func (dto *ProductSearchParams) ToEntity(cats models.ProductCategories) (params *models.ProductParameter) {
	params = &models.ProductParameter{
		Offset:  dto.Offset,
		Keyword: dto.Keyword,
	}

	if len(dto.Categories) == 0 {
		for _, cat := range cats {
			params.Categories = append(params.Categories, strconv.FormatInt(int64(cat.ID), 10))
		}
	} else {
		for _, cat := range dto.Categories {
			params.Categories = append(params.Categories, strconv.FormatUint(cat, 10))
		}
	}

	if dto.Order == 1 {
		params.Order = "p.id"
	} else if dto.Order == 2 {
		params.Order = "p.price"
	} else {
		params.Order = "p.name"
	}

	if dto.Limit == 0 {
		params.Limit = 10
	} else {
		params.Limit = dto.Limit
	}

	return
}

func (dto *NewProductRequest) ToEntity() (product *models.Product) {
	product = &models.Product{
		Name:        dto.Name,
		PictureURL:  dto.PictureURL,
		Description: dto.Description,
		CategoryID:  dto.CategoryID,
		Price:       dto.Price,
		Qty:         dto.Qty,
	}

	return
}

func (dto *NewCategoryRequest) ToEntity() (category *models.ProductCategory) {
	category = &models.ProductCategory{
		Name: dto.Name,
	}

	return
}

func (dto *UpdateProductRequest) ToEntity() (product *models.Product) {
	product = &models.Product{
		ID:          dto.ID,
		Name:        dto.Name,
		PictureURL:  dto.PictureURL,
		Description: dto.Description,
		CategoryID:  dto.CategoryID,
		Price:       dto.Price,
		Qty:         dto.Qty,
	}

	return
}

func (dto *UpdateCategoryRequest) ToEntity() (category *models.ProductCategory) {
	category = &models.ProductCategory{
		ID:   dto.CategoryID,
		Name: dto.Name,
	}

	return
}
