package impl

import (
	"context"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
)

type ProductServiceImpl struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{repo: repo}
}

func (prd *ProductServiceImpl) GetAllProduct(ctx context.Context, limit int64, offset int64) (res dto.ProductsResponse, err error) {

	return
}

func (prd *ProductServiceImpl) GetProductByID(ctx context.Context, id uint64) (res *dto.ProductResponse, err error) {

	return
}

func (prd *ProductServiceImpl) StoreNewProduct(ctx context.Context, data *dto.NewProductRequest) (err error) {

	return
}

func (prd *ProductServiceImpl) UpdateProduct(ctx context.Context, data *dto.UpdateProductRequest) (err error) {

	return
}

func (prd *ProductServiceImpl) DeleteProduct(ctx context.Context, id uint64) (err error) {

	return
}

func (prd *ProductServiceImpl) GetAllCategory(ctx context.Context) (res dto.CategoriesResponse, err error) {

	return
}

func (prd *ProductServiceImpl) StoreNewCategory(ctx context.Context, data *dto.NewCategoryRequest) (err error) {

	return
}

func (prd *ProductServiceImpl) UpdateCategory(ctx context.Context, data *dto.UpdateCategoryRequest) (err error) {

	return
}

func (prd *ProductServiceImpl) DeleteCategory(ctx context.Context, id uint64) (err error) {

	return
}
