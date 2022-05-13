package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/entity"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type ProductRepository interface {
	CountProduct(ctx context.Context) (sum uint64, err error)
	GetAllProduct(ctx context.Context, limit uint64, offset uint64) (res entity.ProductDetails, err error)
	GetProductByID(ctx context.Context, id uint64) (res *entity.Product, err error)
	StoreProduct(ctx context.Context, res *entity.Product) (id int64, err error)
	UpdateProduct(ctx context.Context, res *entity.Product) (err error)
	DeleteProduct(ctx context.Context, id uint64) (err error)

	GetAllCategory(ctx context.Context) (res entity.ProductCategories, err error)
	GetCategoryByID(ctx context.Context, id uint64) (res *entity.ProductCategory, err error)
	StoreCategory(ctx context.Context, res *entity.ProductCategory) (id int64, err error)
	UpdateCategory(ctx context.Context, res *entity.ProductCategory) (err error)
	DeleteCategory(ctx context.Context, id uint64) (err error)
}

type productRepositoryImpl struct {
	db *sql.DB
}

const (
	COUNT_PRODUCT   = `SELECT COUNT(*) FROM product`
	GET_ALL_PRODUCT = `SELECT p.id, p.category_id, p.name, pc.name, p.price, p.qty, p.url, p.description FROM product p
			 LEFT JOIN product_category pc ON p.category_id=pc.id LIMIT ? OFFSET ?`
	GET_PRODUCT_BY_ID = `SELECT p.id, p.category_id, p.name, p.price, p.qty, p.url, p.description FROM product p WHERE p.id=?`
	STORE_NEW_PRODUCT = `INSERT INTO product(category_id, name, price, qty, url, description) VALUES (?, ?, ?, ?, ?, ?)`
	UPDATE_PRODUCT    = `UPDATE product SET category_id = ?, name=?, price=?, qty=?, url=?, description=? WHERE id=?`
	DELETE_PRODUCT    = `DELETE FROM product WHERE id=?`

	GET_ALL_CATEGORY   = `SELECT c.id, c.name FROM product_category c`
	GET_CATEGORY_BY_ID = `SELECT c.id, c.name FROM product_category c WHERE c.id=?`
	STORE_NEW_CATEGORY = `INSERT INTO product_category(name) VALUES (?)`
	UPDATE_CATEGORY    = `UPDATE product_category SET name=? WHERE id=?`
	DELETE_CATEGORY    = `DELETE FROM product_category WHERE id=?`
)

func NewProductRepository(db *database.DatabaseClient) *productRepositoryImpl {
	return &productRepositoryImpl{db: db.DB}
}

func (repo productRepositoryImpl) CountProduct(ctx context.Context) (sum uint64, err error) {
	query, err := repo.db.PrepareContext(ctx, COUNT_PRODUCT)
	if err != nil {
		log.Printf("[CountProduct] failed to prepare query, err => %+v\n", err)
		return
	}

	err = query.QueryRowContext(ctx).Scan(&sum)
	if err != nil {
		log.Printf("[CountProduct] failed to count products, err => %+v\n", err)
		return
	}

	return
}

func (repo productRepositoryImpl) GetAllProduct(ctx context.Context, limit uint64, offset uint64) (res entity.ProductDetails, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ALL_PRODUCT)
	if err != nil {
		log.Printf("[GetAllProduct] failed to prepare query, err => %+v\n", err)
		return
	}

	rows, err := query.QueryContext(ctx, limit, offset*limit)
	if err != nil {
		log.Printf("[GetAllProduct] failed to fetch data from db, err => %+v\n", err)
		return
	}

	return mapProductDetails(rows)
}

func (repo productRepositoryImpl) GetProductByID(ctx context.Context, id uint64) (res *entity.Product, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_PRODUCT_BY_ID)
	if err != nil {
		log.Printf("[GetProductByID] failed to prepare query, err => %+v\n", err)
		return
	}

	row := query.QueryRowContext(ctx, id)
	res, err = mapProduct(row)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[GetProductByID] failed to fetch product, err => %+v\n", err)
		return
	} else if err == sql.ErrNoRows {
		log.Printf("[GetProductByID] product not found, id => %d\n", id)
		return nil, errors.ErrInvalidResources
	}

	return
}

func (repo productRepositoryImpl) StoreProduct(ctx context.Context, res *entity.Product) (id int64, err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[StoreProduct] failed to start new transaction, err => %+v\n", err)
		return
	}

	query, err := tx.PrepareContext(ctx, STORE_NEW_PRODUCT)
	if err != nil {
		log.Printf("[StoreProduct] failed to prepare query, err => %+v\n", err)
		return
	}

	queryRes, err := query.ExecContext(ctx, res.CategoryID, res.Name, res.Price, res.Qty, res.PictureURL, res.Description)
	if err != nil {
		log.Printf("[StoreProduct] failed to insert new product, err => %+v\n", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[StoreProduct] transaction failed, err => %+v\n", err)
		tx.Rollback()
		return
	}

	id, _ = queryRes.LastInsertId()
	return
}

func (repo productRepositoryImpl) UpdateProduct(ctx context.Context, res *entity.Product) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[UpdateProduct] failed to start new transaction, err => %+v\n", err)
		return
	}

	query, err := tx.PrepareContext(ctx, UPDATE_PRODUCT)
	if err != nil {
		log.Printf("[UpdateProduct] failed to prepare query, err => %+v\n", err)
		return
	}

	_, err = query.ExecContext(ctx, res.CategoryID, res.Name, res.Price, res.Qty, res.PictureURL, res.Description, res.ID)
	if err != nil {
		log.Printf("[UpdateProduct] failed to update data, id => %d, err => %+v\n", res.ID, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[UpdateProduct] failed to commit transaction, err => %+v\n", err)
		tx.Rollback()
		return
	}

	return
}

func (repo *productRepositoryImpl) DeleteProduct(ctx context.Context, id uint64) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[DeleteProduct] failed to start new transaction, err => %+v\n", err)
		return
	}

	query, err := tx.PrepareContext(ctx, DELETE_PRODUCT)
	if err != nil {
		log.Printf("[DeleteProduct] failed to prepare query, err => %+v\n", err)
		return
	}

	_, err = query.ExecContext(ctx, id)
	if err != nil {
		log.Printf("[DeleteProduct] failed to delete product, id => %d, err => %+v\n", id, err)
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("[DeleteProduct] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo productRepositoryImpl) GetAllCategory(ctx context.Context) (res entity.ProductCategories, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ALL_CATEGORY)
	if err != nil {
		log.Printf("[GetAllCategory] failed to prepare query, err => %+v\n", err)
		return
	}

	rows, err := query.QueryContext(ctx)
	if err != nil {
		log.Printf("[GetAllCategory] failed to fetch data from db, err => %+v\n", err)
		return
	}

	return mapCategories(rows)
}

func (repo productRepositoryImpl) GetCategoryByID(ctx context.Context, id uint64) (res *entity.ProductCategory, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_CATEGORY_BY_ID)
	if err != nil {
		log.Printf("[GetCategoryByID] failed to prepare query, err => %+v\n", err)
		return
	}

	rows := query.QueryRowContext(ctx, id)
	res, err = mapCategory(rows)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[GetCategoryByID] failed to fetch product, err => %+v\n", err)
		return
	} else if err == sql.ErrNoRows {
		log.Printf("[GetCategoryByID] category not found, id => %d\n", id)
		return nil, errors.ErrInvalidResources
	}

	return
}

func (repo productRepositoryImpl) StoreCategory(ctx context.Context, res *entity.ProductCategory) (id int64, err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[StoreCategory] failed to start new transaction, err => %+v\n", err)
		return
	}

	query, err := tx.PrepareContext(ctx, STORE_NEW_CATEGORY)
	if err != nil {
		log.Printf("[StoreCategory] failed to prepare query, err => %+v\n", err)
		return
	}

	queryRes, err := query.ExecContext(ctx, res.Name)
	if err != nil {
		log.Printf("[StoreCategory] failed to store new category, err => %+v\n", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[StoreCategory] failed to commit transaction, err => %+v\n", err)
		return
	}

	id, _ = queryRes.LastInsertId()
	return
}

func (repo productRepositoryImpl) UpdateCategory(ctx context.Context, res *entity.ProductCategory) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[UpdateCateogry] failed to start new transaction, err => %+v\n", err)
		return
	}

	query, err := tx.PrepareContext(ctx, UPDATE_CATEGORY)
	if err != nil {
		log.Printf("[UpdateCategory] failed to prepare query, err => %+v\n", err)
		return
	}

	_, err = query.ExecContext(ctx, res.Name, res.ID)
	if err != nil {
		log.Printf("[UpdateCategory] failed to update cateogry, id => %d, err => %+v\n", res.ID, err)
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("[UpdateCategory] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo productRepositoryImpl) DeleteCategory(ctx context.Context, id uint64) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[DeleteCategory] failed to start new transaction, err => %+v\n", err)
		return
	}

	query, err := tx.PrepareContext(ctx, DELETE_CATEGORY)
	if err != nil {
		log.Printf("[DeleteCategory] failed to prepare query, err => %+v\n", err)
		return
	}

	_, err = query.ExecContext(ctx, id)
	if err != nil {
		log.Printf("[DeleteCategory] failed to delete category, id => %d, err => %+v\n", id, err)
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("[DeleteCategory] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func mapCategory(r *sql.Row) (res *entity.ProductCategory, err error) {
	res = &entity.ProductCategory{}
	err = r.Scan(&res.ID, &res.Name)
	return
}

func mapCategories(r *sql.Rows) (res entity.ProductCategories, err error) {
	res = entity.ProductCategories{}

	for r.Next() {
		temp := &entity.ProductCategory{}
		err = r.Scan(&temp.ID, &temp.Name)
		if err != nil {
			log.Printf("[mapCategories] error while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
	}

	return
}

func mapProduct(r *sql.Row) (res *entity.Product, err error) {
	res = &entity.Product{}
	err = r.Scan(&res.ID, &res.CategoryID, &res.Name, &res.Price, &res.Qty, &res.PictureURL, &res.Description)

	return
}

func mapProductDetails(r *sql.Rows) (res entity.ProductDetails, err error) {
	res = entity.ProductDetails{}

	for r.Next() {
		temp := &entity.ProductDetail{}
		err = r.Scan(&temp.ID, &temp.CategoryID, &temp.Name, &temp.CategoryName, &temp.Price, &temp.Qty, &temp.PictureURL, &temp.Description)
		if err != nil {
			log.Printf("[mapProductDetails] error while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
	}

	return
}
