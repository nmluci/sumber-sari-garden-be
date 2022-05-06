package dto

import "github.com/nmluci/sumber-sari-garden/internal/entity"

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

func (dto *NewProductRequest) ToEntity() (product *entity.Product) {
	product = &entity.Product{
		Name:        dto.Name,
		PictureURL:  dto.PictureURL,
		Description: dto.Description,
		CategoryID:  dto.CategoryID,
		Price:       dto.Price,
		Qty:         dto.Qty,
	}

	return
}

func (dto *NewCategoryRequest) ToEntity() (category *entity.ProductCategory) {
	category = &entity.ProductCategory{
		Name: dto.Name,
	}

	return
}

func (dto *UpdateProductRequest) ToEntity() (product *entity.Product) {
	product = &entity.Product{
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

func (dto *UpdateCategoryRequest) ToEntity() (category *entity.ProductCategory) {
	category = &entity.ProductCategory{
		ID: dto.CategoryID,
		Name: dto.Name,
	}

	return
}
