package impl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/nmluci/sumber-sari-garden/internal/models"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type ProductRepository interface {
	CountProduct(ctx context.Context) (sum uint64, err error)
	GetAllProduct(ctx context.Context, params *models.ProductParameter) (res models.ProductDetails, err error)
	GetProductByID(ctx context.Context, id uint64) (res *models.Product, err error)
	StoreProduct(ctx context.Context, res *models.Product) (id int64, err error)
	UpdateProduct(ctx context.Context, res *models.Product) (err error)
	DeleteProduct(ctx context.Context, id uint64) (err error)

	GetAllCategory(ctx context.Context) (res models.ProductCategories, err error)
	GetCategoryByID(ctx context.Context, id uint64) (res *models.ProductCategory, err error)
	StoreCategory(ctx context.Context, res *models.ProductCategory) (id int64, err error)
	UpdateCategory(ctx context.Context, res *models.ProductCategory) (err error)
	DeleteCategory(ctx context.Context, id uint64) (err error)

	GetActiveCoupon(ctx context.Context, limit int64, offset int64) (res models.Coupons, err error)
	GetCouponByCode(ctx context.Context, code string) (res *models.Coupon, err error)
	GetCouponByID(ctx context.Context, id int64) (res *models.Coupon, err error)
	GetAllCoupon(ctx context.Context, limit int64, offset int64) (res models.Coupons, err error)
	StoreCoupon(ctx context.Context, data *models.Coupon) (err error)
	UpdateCoupon(ctx context.Context, id int64, data *models.Coupon) (err error)
	DeleteCoupon(ctx context.Context, id int64) (err error)
}

type productRepositoryImpl struct {
	db *sql.DB
}

const (
	COUNT_PRODUCT   = `SELECT COUNT(*) FROM product`
	GET_ALL_PRODUCT = `SELECT p.id, p.category_id, p.name, pc.name, p.price, p.qty, p.url, p.description FROM product p
		LEFT JOIN product_category pc ON p.category_id=pc.id WHERE p.name LIKE "%%%s%%" AND pc.id IN (%s) ORDER BY %s LIMIT ? OFFSET ?`

	GET_PRODUCT_BY_ID = `SELECT p.id, p.category_id, p.name, p.price, p.qty, p.url, p.description FROM product p WHERE p.id=?`
	STORE_NEW_PRODUCT = `INSERT INTO product(category_id, name, price, qty, url, description) VALUES (?, ?, ?, ?, ?, ?)`
	UPDATE_PRODUCT    = `UPDATE product SET category_id = ?, name=?, price=?, qty=?, url=?, description=? WHERE id=?`
	DELETE_PRODUCT    = `DELETE FROM product WHERE id=?`

	GET_ALL_CATEGORY   = `SELECT c.id, c.name FROM product_category c ORDER BY c.id ASC`
	GET_CATEGORY_BY_ID = `SELECT c.id, c.name FROM product_category c WHERE c.id=?`
	STORE_NEW_CATEGORY = `INSERT INTO product_category(name) VALUES (?)`
	UPDATE_CATEGORY    = `UPDATE product_category SET name=? WHERE id=?`
	DELETE_CATEGORY    = `DELETE FROM product_category WHERE id=?`

	GET_ALL_ACTIVE_COUPON = `SELECT c.id, c.code, c.amount, c.description, c.expired_at FROM coupon c WHERE c.expired_at > NOW() LIMIT ? OFFSET ?`
	GET_ALL_COUPON        = `SELECT c.id, c.code, c.amount, c.description, c.expired_at FROM coupon c LIMIT ? OFFSET ?`
	GET_COUPON_BY_CODE    = `SELECT c.id, c.code, c.amount, c.description, c.expired_at FROM coupon c WHERE c.code LIKE "%%%s%%" LIMIT 1`
	GET_COUPON_BY_ID      = `SELECT c.id, c.code, c.amount, c.description, c.expired_at FROM coupon c WHERE c.id=? LIMIT 1`
	STORE_NEW_COUPON      = `INSERT INTO coupon(code, amount, description, expired_at) VALUES (?, ?, ?, ?)`
	UPDATE_COUPON         = `UPDATE coupon SET code=?, amount=?, description=?, expired_at=? WHERE id=?`
	DELETE_COUPON         = `DELETE FROM coupon WHERE id=?`
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

func (repo productRepositoryImpl) GetAllProduct(ctx context.Context, params *models.ProductParameter) (res models.ProductDetails, err error) {
	query, err := repo.db.PrepareContext(ctx, fmt.Sprintf(GET_ALL_PRODUCT, params.Keyword, strings.Join(params.Categories, ", "), params.Order))
	if err != nil {
		log.Printf("[GetAllProduct] failed to prepare query, err => %+v\n", err)
		return
	}

	rows, err := query.QueryContext(ctx, params.Limit, params.Limit*params.Offset)
	if err != nil {
		log.Printf("[GetAllProduct] failed to fetch data from db, err => %+v\n", err)
		return
	}

	return mapProductDetails(rows)
}

func (repo productRepositoryImpl) GetProductByID(ctx context.Context, id uint64) (res *models.Product, err error) {
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

func (repo productRepositoryImpl) StoreProduct(ctx context.Context, res *models.Product) (id int64, err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[StoreProduct] failed to start new transaction, err => %+v\n", err)
		return
	}

	defer tx.Rollback()

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
		return
	}

	id, _ = queryRes.LastInsertId()
	return
}

func (repo productRepositoryImpl) UpdateProduct(ctx context.Context, res *models.Product) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[UpdateProduct] failed to start new transaction, err => %+v\n", err)
		return
	}

	defer tx.Rollback()

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

	defer tx.Rollback()

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
		log.Printf("[DeleteProduct] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo productRepositoryImpl) GetAllCategory(ctx context.Context) (res models.ProductCategories, err error) {
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

func (repo productRepositoryImpl) GetCategoryByID(ctx context.Context, id uint64) (res *models.ProductCategory, err error) {
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

func (repo productRepositoryImpl) StoreCategory(ctx context.Context, res *models.ProductCategory) (id int64, err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[StoreCategory] failed to start new transaction, err => %+v\n", err)
		return
	}

	defer tx.Rollback()

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

func (repo productRepositoryImpl) UpdateCategory(ctx context.Context, res *models.ProductCategory) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[UpdateCateogry] failed to start new transaction, err => %+v\n", err)
		return
	}

	defer tx.Rollback()

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

	defer tx.Rollback()

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
		log.Printf("[DeleteCategory] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo productRepositoryImpl) GetActiveCoupon(ctx context.Context, limit int64, offset int64) (res models.Coupons, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ALL_ACTIVE_COUPON)
	if err != nil {
		log.Printf("[GetActiveCoupon] failed to prepare query, err => %+v", err)
		return
	}

	rows, err := query.QueryContext(ctx, limit, offset*limit)
	if err != nil {
		log.Printf("[GetActiveCoupon] failed to fetch coupon, err => %+v\n", err)
		return
	}

	res, err = mapCoupons(rows)
	if err != nil {
		log.Printf("[GetActiveCoupon] failed to fetch coupon, err => %+v\n", err)
		return
	}

	return
}

func (repo productRepositoryImpl) GetAllCoupon(ctx context.Context, limit int64, offset int64) (res models.Coupons, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ALL_COUPON)
	if err != nil {
		log.Printf("[GetAllCoupon] failed to prepare query, err => %+v", err)
		return
	}

	rows, err := query.QueryContext(ctx, limit, offset*limit)
	if err != nil {
		log.Printf("[GetAllCoupon] failed to fetch coupon, err => %+v\n", err)
		return
	}

	res, err = mapCoupons(rows)
	if err != nil {
		log.Printf("[GetAllCoupon] failed to fetch coupon, err => %+v\n", err)
		return
	}

	return
}

func (repo *productRepositoryImpl) GetCouponByCode(ctx context.Context, code string) (res *models.Coupon, err error) {
	data := repo.db.QueryRowContext(ctx, fmt.Sprintf(GET_COUPON_BY_CODE, code))
	res, err = mapCoupon(data)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[GetCouponByCode] failed to fetch coupon by code, code => %s, err => %+v\n", code, err)
		return
	} else if err == sql.ErrNoRows {
		log.Printf("[GetCouponByCode] coupon not existed\n")
		return nil, errors.ErrInvalidResources
	}

	return
}

func (repo *productRepositoryImpl) GetCouponByID(ctx context.Context, id int64) (res *models.Coupon, err error) {
	data := repo.db.QueryRowContext(ctx, GET_COUPON_BY_ID, id)
	res, err = mapCoupon(data)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[GetCouponByID] failed to fetch coupon by code, code => %d, err => %+v\n", id, err)
		return
	} else if err == sql.ErrNoRows {
		log.Printf("[GetCouponByID] coupon not existed\n")
		return nil, errors.ErrInvalidResources
	}

	return
}

func (repo productRepositoryImpl) StoreCoupon(ctx context.Context, data *models.Coupon) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[StoreCoupon] failed to start new transaction, err => %+v\n", err)
		return
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, STORE_NEW_COUPON, data.Code, data.Amount, data.Description, data.ExpiredAt)
	if err != nil {
		log.Printf("[StoreCoupon] failed to store new coupon, err => %+v\n", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[StoreCoupon] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo productRepositoryImpl) UpdateCoupon(ctx context.Context, id int64, data *models.Coupon) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[UpdateCoupon] failed to start new transaction, err => %+v\n", err)
		return
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, UPDATE_COUPON, data.Code, data.Amount, data.Description, data.ExpiredAt, id)
	if err != nil {
		log.Printf("[UpdateCoupon] failed to update coupon data, err => %+v\n", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[UpdateCoupon] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo productRepositoryImpl) DeleteCoupon(ctx context.Context, id int64) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[DeleteCoupon] failed to start new transaction, err => %+v\n", err)
		return
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, DELETE_COUPON, id)
	if err != nil {
		log.Printf("[DeleteCoupon] failed to delete coupon, id => %d, err => %+v\n", id, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[DeleteCoupon] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func mapCategory(r *sql.Row) (res *models.ProductCategory, err error) {
	res = &models.ProductCategory{}
	err = r.Scan(&res.ID, &res.Name)
	return
}

func mapCategories(r *sql.Rows) (res models.ProductCategories, err error) {
	res = models.ProductCategories{}

	for r.Next() {
		temp := &models.ProductCategory{}
		err = r.Scan(&temp.ID, &temp.Name)
		if err != nil {
			log.Printf("[mapCategories] error while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
	}

	return
}

func mapProduct(r *sql.Row) (res *models.Product, err error) {
	res = &models.Product{}
	err = r.Scan(&res.ID, &res.CategoryID, &res.Name, &res.Price, &res.Qty, &res.PictureURL, &res.Description)

	return
}

func mapProductDetails(r *sql.Rows) (res models.ProductDetails, err error) {
	res = models.ProductDetails{}

	for r.Next() {
		temp := &models.ProductDetail{}
		err = r.Scan(&temp.ID, &temp.CategoryID, &temp.Name, &temp.CategoryName, &temp.Price, &temp.Qty, &temp.PictureURL, &temp.Description)
		if err != nil {
			log.Printf("[mapProductDetails] error while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
	}

	return
}

func mapCoupons(r *sql.Rows) (res models.Coupons, err error) {
	res = models.Coupons{}

	for r.Next() {
		temp := &models.Coupon{}
		err = r.Scan(&temp.ID, &temp.Code, &temp.Amount, &temp.Description, &temp.ExpiredAt)
		if err != nil {
			log.Printf("[mapCoupons] an error occured while parsing query result, err => %+v\n", err)
			return nil, err
		}

		res = append(res, temp)
	}

	return
}

func mapCoupon(r *sql.Row) (res *models.Coupon, err error) {
	res = &models.Coupon{}
	err = r.Scan(&res.ID, &res.Code, &res.Amount, &res.Description, &res.ExpiredAt)
	return
}
