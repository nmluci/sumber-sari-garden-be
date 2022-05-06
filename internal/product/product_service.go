package product

import (
	"context"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/product/impl"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
)

type ProductService interface {
	GetAllProduct(ctx context.Context, limit int64, offset int64) (res dto.ProductsResponse, err error)
	GetProductByID(ctx context.Context, id uint64) (res *dto.ProductResponse, err error)
	StoreNewProduct(ctx context.Context, data *dto.NewProductRequest) (err error)
	UpdateProduct(ctx context.Context, data *dto.UpdateProductRequest) (err error)
	DeleteProduct(ctx context.Context, id uint64) (err error)

	GetAllCategory(ctx context.Context) (res dto.CategoriesResponse, err error)
	StoreNewCategory(ctx context.Context, data *dto.NewCategoryRequest) (err error)
	UpdateCategory(ctx context.Context, data *dto.UpdateCategoryRequest) (err error)
	DeleteCategory(ctx context.Context, id uint64) (err error)
}

func NewProductService(db *database.DatabaseClient) ProductService {
	repo := impl.NewProductRepository(db)
	return impl.NewProductService(repo)
}
