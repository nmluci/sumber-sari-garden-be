package impl

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/entity"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type UsercartRepository interface {
	NewCart(ctx context.Context, usr int64) (orderID int64, err error)
	GetCartByUserID(ctx context.Context, usr int64) (res *entity.OrderData, err error)
	GetItemsByOrderID(ctx context.Context, orderID uint64) (res entity.OrderDetails, err error)
	GetCartMetadataByOrderID(ctx context.Context, orderID uint64) (res *entity.OrderMetadata, err error)
	GetItem(ctx context.Context, orderID uint64, productID uint64) (res *entity.OrderDetail, err error)
	InsertItem(ctx context.Context, orderID uint64, productID uint64, qty uint64) (err error)
	UpdateItem(ctx context.Context, orderID uint64, productID uint64, qty uint64) (err error)
	RemoveItem(ctx context.Context, orderID uint64, productID uint64) (err error)
	GetHistoryMetadata(ctx context.Context, params dto.HistoryParams) (meta []*entity.OrderHistoryMetadata, err error)
	GetCouponByCode(ctx context.Context, code string) (res *entity.ActiveCoupon, err error)
	Checkout(ctx context.Context, userID int64, orderID uint64, couponID *uint64) (err error)
}

type usercartRepositoryImpl struct {
	db *sql.DB
}

var (
	GET_ORDER_DETAIL_BY_ID = `SELECT od.id, od.order_id, od.product_id, p.name, p.price, od.qty, 
			(case when o.coupon_id is not null then (c.amount*(p.price*od.qty)) else 0 end) disc, 
			((p.price*od.qty)-(case when o.coupon_id is not null then (c.amount*(p.price*od.qty)) else 0 end)) subtotal FROM order_detail od
			LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id
			LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.id=?`
	GET_ORDER_METADATA = `SELECT o.id, SUM((p.price*od.qty)-(case when o.coupon_id is not null then (c.amount*(p.price*od.qty)) else 0 end)) grand_total, 
			COUNT(*) item_count FROM order_detail od LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id
			LEFT JOIN coupon c ON o.coupon_id=c.id GROUP BY o.id HAVING o.id=?`

	HISTORY_METADATA_A = `SELECT o.id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then (c.amount*(p.price*od.qty)) else 0 end)) grand_total, 
			COUNT(*) item_count, c.code, s.name FROM order_detail od 
			LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id
			LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.status_id = 1 GROUP BY o.id, o.user_id HAVING o.user_id=?`

	HISTORY_METADATA_B = `SELECT o.id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then (c.amount*(p.price*od.qty)) else 0 end)) grand_total, 
			COUNT(*) item_count, c.code, s.name FROM order_detail od 
			LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id
			LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.status_id = 2 GROUP BY o.id, o.user_id HAVING o.user_id=?`

	HISTORY_METADATA_C = `SELECT o.id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then (c.amount*(p.price*od.qty)) else 0 end)) grand_total, 
			COUNT(*) item_count, c.code, s.name FROM order_detail od 
			LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id
			LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.status_id = 3 GROUP BY o.id, o.user_id HAVING o.user_id=?`

	GET_H_METADATA = fmt.Sprintf("SELECT * FROM (%s UNION %s UNION %s) t WHERE (t.created_at BETWEEN ? AND ?) LIMIT ? OFFSET ?", HISTORY_METADATA_A, HISTORY_METADATA_B, HISTORY_METADATA_C)
	// GET_H_METADATA = strings.Join([]string{HISTORY_METADATA_A, HISTORY_METADATA_B, HISTORY_METADATA_C}, " UNION ") + " WHERE (o.created_at BETWEEN ? AND ?) LIMIT ? OFFSET ?"

	GET_ORDER = `SELECT o.id, o.user_id, o.status_id, s.name, o.coupon_id FROM order_data o 
			LEFT JOIN order_status s ON o.status_id=s.id WHERE o.user_id=? AND o.status_id=1 ORDER BY o.created_at DESC LIMIT 1`
	NEW_ORDER = `INSERT INTO order_data(user_id) VALUES (?)`

	GET_ITEM    = GET_ORDER_DETAIL_BY_ID + " AND p.id = ?"
	INSERT_ITEM = `INSERT INTO order_detail(order_id, product_id, qty) VALUES (?, ?, ?)`
	UPDATE_ITEM = `UPDATE order_detail SET qty=? WHERE product_id=? AND order_id=?`
	DELETE_ITEM = `DELETE FROM order_detail WHERE product_id=? AND order_id=?`

	GET_COUPON = `SELECT c.id, c.code, c.amount FROM coupon c WHERE c.code=? AND c.expired_at > NOW()`

	CHECKOUT_ORDER = `UPDATE order_data SET created_at=NOW(), status_id=2, coupon_id=? WHERE order_data.user_id=? AND order_data.id=?`
)

func NewUsercartRepository(db *database.DatabaseClient) *usercartRepositoryImpl {
	return &usercartRepositoryImpl{db: db.DB}
}

func (repo *usercartRepositoryImpl) NewCart(ctx context.Context, user int64) (orderID int64, err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[NewCart] failed to start new transaction, err => %+v\n", err)
		return 0, err
	}

	defer tx.Rollback()

	query, err := tx.PrepareContext(ctx, NEW_ORDER)
	if err != nil {
		log.Printf("[NewCart] failed to prepare query, err => %+v\n", err)
		return 0, err
	}

	res, err := query.ExecContext(ctx, user)
	if err != nil {
		log.Printf("[NewCart] failed to create new order, err => %+v\n", err)
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("[NewCart] failed to commit transaction, err => %+v\n", err)
		return
	}

	lid, _ := res.LastInsertId()

	return lid, nil
}

func (repo *usercartRepositoryImpl) GetCartByUserID(ctx context.Context, usr int64) (res *entity.OrderData, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ORDER)
	if err != nil {
		log.Printf("[GetCartByUserID] failed to prepare query, err => %+v\n", err)
		return
	}

	row := query.QueryRowContext(ctx, usr)
	res, err = mapOrderData(row)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[GetCartByUserID] failed to fetch usercart, err => %+v\n", err)
		return
	} else if err == sql.ErrNoRows {
		log.Printf("[GetCartByUserID] cart not found, user_id => %d\n", usr)
		return nil, errors.ErrInvalidResources
	}

	return
}

func (repo *usercartRepositoryImpl) GetItemsByOrderID(ctx context.Context, orderID uint64) (res entity.OrderDetails, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ORDER_DETAIL_BY_ID)
	if err != nil {
		log.Printf("[GetItemsByOrderID] failed to prepare query, err => %+v\n", err)
		return
	}

	rows, err := query.QueryContext(ctx, orderID)
	if err != nil {
		log.Printf("[GetItemsByOrderID] failed to fetch cart items, err => %+v\n", err)
		return
	}

	res, err = mapOrderDetails(rows)
	if err != nil {
		log.Printf("[GetItemsByOrderID] failed to parse cart items, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) GetCartMetadataByOrderID(ctx context.Context, orderID uint64) (res *entity.OrderMetadata, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ORDER_METADATA)
	if err != nil {
		log.Printf("[GetCartMetadataByOrderID] failed to prepare query, err => %+v\n", err)
		return
	}

	row := query.QueryRowContext(ctx, orderID)
	res, err = mapOrderMetadata(row)
	if err != nil {
		log.Printf("[GetCartMetadataByOrderID] failed to parse query result, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) GetItem(ctx context.Context, orderID uint64, productID uint64) (res *entity.OrderDetail, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ITEM)
	if err != nil {
		log.Printf("[GetItem] failed to prepare query, err => %+v\n", err)
		return
	}

	row := query.QueryRowContext(ctx, orderID, productID)
	res, err = mapOrderDetail(row)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[GetItem] failed to fetch item, err => %+v\n", err)
		return
	} else if err == sql.ErrNoRows {
		log.Printf("[GetItem] item not found, order_id => %d, productID => %d\n", orderID, productID)
		return nil, errors.ErrInvalidResources
	}
	return
}

func (repo *usercartRepositoryImpl) InsertItem(ctx context.Context, orderID uint64, productID uint64, qty uint64) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[InsertItem] failed to start new transaction, err => %v\n", err)
		return
	}

	defer tx.Rollback()

	query, err := tx.PrepareContext(ctx, INSERT_ITEM)
	if err != nil {
		log.Printf("[InsertItem] failed to prepare query, err => %+v\n", err)
		return
	}

	_, err = query.ExecContext(ctx, orderID, productID, qty)
	if err != nil {
		log.Printf("[InsertItem] failed to insert new item, orderID => %d, pid => %d, err => %+v\n", orderID, productID, err)
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("[InsertItem] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) UpdateItem(ctx context.Context, orderID uint64, productID uint64, qty uint64) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[UpdateItem] failed to start new transaction, err => %v\n", err)
		return
	}

	defer tx.Rollback()

	query, err := tx.PrepareContext(ctx, UPDATE_ITEM)
	if err != nil {
		log.Printf("[UpdateItem] failed to prepare query, err => %+v\n", err)
		return
	}

	_, err = query.ExecContext(ctx, qty, productID, orderID)
	if err != nil {
		log.Printf("[UpdateItem] failed to insert new item, orderID => %d, pid => %d, err => %+v\n", orderID, productID, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[UpdateItem] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) RemoveItem(ctx context.Context, orderID uint64, productID uint64) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[DeleteItem] failed to start new transaction, err => %v\n", err)
		return
	}

	defer tx.Rollback()

	query, err := tx.PrepareContext(ctx, DELETE_ITEM)
	if err != nil {
		log.Printf("[DeleteItem] failed to prepare query, err => %+v\n", err)
		return
	}

	_, err = query.ExecContext(ctx, productID, orderID)
	if err != nil {
		log.Printf("[DeleteItem] failed to insert new item, orderID => %d, pid => %d, err => %+v\n", orderID, productID, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[DeleteItem] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) GetHistoryMetadata(ctx context.Context, params dto.HistoryParams) (meta []*entity.OrderHistoryMetadata, err error) {
	log.Println(params)

	query, err := repo.db.PrepareContext(ctx, GET_H_METADATA)
	if err != nil {
		log.Printf("[GetHistoryMetadata] failed to prepare query, err => %+v\n", err)
		return
	}

	rows, err := query.QueryContext(ctx, params.UserID, params.UserID, params.UserID, params.DateStart, params.DateEnd, params.Limit, params.Offset*params.Limit)
	if err != nil {
		log.Printf("[GetHistoryMetadata] failed to fetch order metadatas, err => %+v\n", err)
		return
	}

	meta, err = mapHistory(rows)
	if err != nil {
		log.Printf("[GetHistoryMetadata] failed to parse order metadatas, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) GetCouponByCode(ctx context.Context, code string) (res *entity.ActiveCoupon, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_COUPON)
	if err != nil {
		log.Printf("[GetCouponByCode] failed to prepare query, err => %+v", err)
		return
	}

	rows := query.QueryRowContext(ctx, code)

	res, err = mapCoupon(rows)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[GetCouponByCode] failed to fetch coupon, err => %+v\n", err)
		return
	} else if err == sql.ErrNoRows {
		log.Printf("[GetCouponByCode] coupon not found, code => %s\n", code)
		return nil, errors.ErrInvalidResources
	}

	return
}

func (repo *usercartRepositoryImpl) Checkout(ctx context.Context, userID int64, orderID uint64, couponID *uint64) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[Checkout] failed to start new transaction, err => %v\n", err)
		return
	}

	defer tx.Rollback()

	query, err := tx.PrepareContext(ctx, CHECKOUT_ORDER)
	if err != nil {
		log.Printf("[Checkout] failed to prepare query, err => %+v\n", err)
		return
	}

	log.Println(couponID, userID, orderID)
	_, err = query.ExecContext(ctx, couponID, userID, orderID)
	if err != nil {
		log.Printf("[Checkout] failed to checkout an order, orderID => %d, err => %+v\n", orderID, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[Checkout] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func mapOrderData(r *sql.Row) (res *entity.OrderData, err error) {
	res = &entity.OrderData{}
	err = r.Scan(&res.ID, &res.UserID, &res.StatusID, &res.StatusName, &res.CouponID)
	return
}

func mapOrderDetails(r *sql.Rows) (res entity.OrderDetails, err error) {
	res = entity.OrderDetails{}

	for r.Next() {
		temp := &entity.OrderDetail{}
		err = r.Scan(&temp.ID, &temp.OrderID, &temp.ProductID, &temp.ProductName, &temp.Price, &temp.Qty, &temp.Disc, &temp.SubTotal)
		if err != nil {
			log.Printf("[mapOrderDetails] an error occured while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
	}

	return
}

func mapOrderMetadata(r *sql.Row) (res *entity.OrderMetadata, err error) {
	res = &entity.OrderMetadata{}
	err = r.Scan(&res.OrderID, &res.GrandTotal, &res.ItemCount)
	return
}

func mapOrderDetail(r *sql.Row) (res *entity.OrderDetail, err error) {
	res = &entity.OrderDetail{}
	err = r.Scan(&res.ID, &res.OrderID, &res.ProductID, &res.ProductName, &res.Price, &res.Qty, &res.Disc, &res.SubTotal)
	return
}

func mapHistory(r *sql.Rows) (res []*entity.OrderHistoryMetadata, err error) {
	res = []*entity.OrderHistoryMetadata{}

	for r.Next() {
		temp := &entity.OrderHistoryMetadata{}
		err = r.Scan(&temp.OrderID, &temp.OrderDate, &temp.GrandTotal, &temp.ItemCount, &temp.CouponName, &temp.StatusName)
		if err != nil {
			log.Printf("[mapHistory] an error occured while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
	}

	return
}

func mapCoupon(r *sql.Row) (res *entity.ActiveCoupon, err error) {
	res = &entity.ActiveCoupon{}
	err = r.Scan(&res.ID, &res.Code, &res.Amount, &res.ExpiredAt)
	return
}
