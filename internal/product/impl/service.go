package impl

import (
	"context"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type ProductServiceImpl struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{repo: repo}
}

func (prd *ProductServiceImpl) GetAllProduct(ctx context.Context, limit uint64, offset uint64) (res dto.ProductsResponse, err error) {
	itemCount, err := prd.repo.CountProduct(ctx)
	if err != nil {
		err = errors.ErrInvalidResources
		return
	}

	data, err := prd.repo.GetAllProduct(ctx, limit, offset)
	if err != nil {
		log.Printf("[GetAllProduct] an error occured while fetching data, err => %+v\n", err)
		return
	}

	return dto.NewProductsResponse(data, limit, offset, itemCount)
}

func (prd *ProductServiceImpl) GetProductByID(ctx context.Context, id uint64) (res *dto.ProductResponse, err error) {
	data, err := prd.repo.GetProductByID(ctx, id)
	if data == nil {
		err = errors.ErrNotFound
		log.Printf("[GetProductByID] data not found, id => %d\n", id)
		return
	} else if err != nil {
		log.Printf("[GetProductByID] an error occured while fetching data, err => %+v\n", err)
		return
	}

	cat, err := prd.repo.GetCategoryByID(ctx, data.CategoryID)
	if err != nil {
		log.Printf("[GetProductByID] an error occured while fetching product category, err => %+v\n", err)
		return
	}

	return dto.NewProductResponse(data, cat)
}

func (prd *ProductServiceImpl) StoreNewProduct(ctx context.Context, data *dto.NewProductRequest) (err error) {
	res := data.ToEntity()
	cat, err := prd.repo.GetCategoryByID(ctx, res.CategoryID)
	if err != nil {
		if cat == nil {
			log.Printf("[StoreNewProduct] category isn't valid, cat_id => %d\n", res.CategoryID)
			err = errors.ErrInvalidRequestBody
		} else {
			log.Printf("[StoreNewProduct] an error occured while validating product category, cat_id => %d, err => %+v\n", res.CategoryID, err)
		}
		return
	}

	// TODO: log new product by ID
	_, err = prd.repo.StoreProduct(ctx, res)
	if err != nil {
		log.Printf("[StoreNewProduct] an error occured while storing new product, err => %+v\n", err)
		return
	}

	return
}

func (prd *ProductServiceImpl) UpdateProduct(ctx context.Context, data *dto.UpdateProductRequest) (err error) {
	res := data.ToEntity()

	err = prd.repo.UpdateProduct(ctx, res)
	if err != nil {
		log.Printf("[UpdateProduct] an error occured while updating product, id => %d. err => %+v\n", res.ID, err)
		return
	}

	return
}

func (prd *ProductServiceImpl) DeleteProduct(ctx context.Context, id uint64) (err error) {
	err = prd.repo.DeleteProduct(ctx, id)
	if err != nil {
		log.Printf("[DeleteProduct] an error occured while deleting product, id => %v, err => %+v\n", id, err)
		return
	}

	return
}

func (prd *ProductServiceImpl) GetAllCategory(ctx context.Context) (res dto.CategoriesResponse, err error) {
	data, err := prd.repo.GetAllCategory(ctx)
	if err != nil {
		log.Printf("[GetAllCategory] an error occured while fetching categories, err => %+v\n", err)
		return
	}

	return dto.NewCategoriesResponse(data)
}

func (prd *ProductServiceImpl) StoreNewCategory(ctx context.Context, data *dto.NewCategoryRequest) (err error) {
	res := data.ToEntity()
	_, err = prd.repo.StoreCategory(ctx, res)
	if err != nil {
		log.Printf("[StoreNewCategory] an error occured while storing new category, err => %+V\n", err)
	}

	return
}

func (prd *ProductServiceImpl) UpdateCategory(ctx context.Context, data *dto.UpdateCategoryRequest) (err error) {
	res := data.ToEntity()
	err = prd.repo.UpdateCategory(ctx, res)
	if err != nil {
		log.Printf("[UpdateCategory] an error occured while storing new category, err => %+V\n", err)
	}

	return
}

func (prd *ProductServiceImpl) DeleteCategory(ctx context.Context, id uint64) (err error) {
	err = prd.repo.DeleteCategory(ctx, id)
	if err != nil {
		log.Printf("[UpdateCategory] an error occured while storing new category, err => %+V\n", err)
	}

	return
}
