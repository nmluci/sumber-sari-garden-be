package impl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/models"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type UsercartRepository interface {
	NewCart(ctx context.Context, usr int64) (orderID int64, err error)
	GetCartByUserID(ctx context.Context, usr int64) (res *models.OrderData, err error)
	GetItemsByOrderID(ctx context.Context, orderID uint64) (res models.OrderDetails, err error)
	GetCartMetadataByOrderID(ctx context.Context, orderID uint64) (res *models.OrderMetadata, err error)
	GetItem(ctx context.Context, orderID uint64, productID uint64) (res *models.OrderDetail, err error)
	InsertItem(ctx context.Context, orderID uint64, productID uint64, qty uint64) (err error)
	UpdateItem(ctx context.Context, orderID uint64, productID uint64, qty uint64) (err error)
	RemoveItem(ctx context.Context, orderID uint64, productID uint64) (err error)
	GetHistoryMetadata(ctx context.Context, params dto.HistoryParams) (meta []*models.OrderHistoryMetadata, err error)
	GetHistoryMetadataByID(ctx context.Context, productID uint64) (meta *models.OrderHistoryMetadata, err error)
	GetCouponByCode(ctx context.Context, code string) (res *models.Coupon, err error)
	Checkout(ctx context.Context, userID int64, orderID uint64, couponID *uint64) (err error)
	VerifyOrder(ctx context.Context, orderID uint64) (err error)
	GetUnpaidOrder(ctx context.Context) (res []*models.OrderMetadata, err error)
	GetHistoryMetadataAll(ctx context.Context, params dto.HistoryParams) (meta []*models.OrderHistoryMetadata, err error)
	GetStatistics(ctx context.Context, dateStart time.Time, dateEnd time.Time) (res *models.Statistics, err error)
}

type usercartRepositoryImpl struct {
	db *sql.DB
}

var (
	GET_ORDER_DETAIL_BY_ID = `SELECT od.id, od.order_id, od.product_id, p.name, p.price, od.qty, 
		(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end) disc, 
		((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) subtotal FROM order_detail od
		LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id
		LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.id=?`
	GET_ORDER_METADATA = `SELECT o.id, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
		COUNT(*) item_count FROM order_detail od LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id
		LEFT JOIN coupon c ON o.coupon_id=c.id GROUP BY o.id HAVING o.id=?`

	HISTORY_METADATA = `SELECT o.id, o.user_id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
		COUNT(*) item_count, c.code, s.name FROM order_detail od 
		LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id
		LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.status_id = ? GROUP BY o.id, o.user_id`

	GET_H_METADATA = fmt.Sprintf("SELECT * FROM (%s UNION %s UNION %s) t HAVING t.user_id=? ORDER BY t.created_at DESC LIMIT ? OFFSET ?",
		HISTORY_METADATA, HISTORY_METADATA, HISTORY_METADATA)

	GET_H_METADATA_BY_ID = fmt.Sprintf("SELECT * FROM (%s UNION %s UNION %s) t WHERE t.id=?",
		HISTORY_METADATA, HISTORY_METADATA, HISTORY_METADATA)

	GET_H_METADATA_ALL = fmt.Sprintf("SELECT * FROM (%s UNION %s UNION %s) t WHERE (t.created_at BETWEEN ? AND ?) ORDER BY t.created_at DESC LIMIT ? OFFSET ?",
		HISTORY_METADATA, HISTORY_METADATA, HISTORY_METADATA)
	// GET_H_METADATA = strings.Join([]string{HISTORY_METADATA_A, HISTORY_METADATA_B, HISTORY_METADATA_C}, " UNION ") + " WHERE (o.created_at BETWEEN ? AND ?) LIMIT ? OFFSET ?"

	GET_ORDER = `SELECT o.id, o.user_id, o.status_id, s.name, o.coupon_id FROM order_data o 
			LEFT JOIN order_status s ON o.status_id=s.id WHERE o.user_id=? AND o.status_id=1 ORDER BY o.created_at DESC LIMIT 1`
	NEW_ORDER = `INSERT INTO order_data(user_id) VALUES (?)`

	GET_ITEM    = GET_ORDER_DETAIL_BY_ID + " AND p.id = ?"
	INSERT_ITEM = `INSERT INTO order_detail(order_id, product_id, qty) VALUES (?, ?, ?)`
	UPDATE_ITEM = `UPDATE order_detail SET qty=? WHERE product_id=? AND order_id=?`
	DELETE_ITEM = `DELETE FROM order_detail WHERE product_id=? AND order_id=?`

	GET_COUPON = `SELECT c.id, c.code, (c.amount/100), c.expired_at FROM coupon c WHERE c.code=? AND c.expired_at > NOW()`

	CHECKOUT_ORDER = `UPDATE order_data SET created_at=NOW(), status_id=2, coupon_id=? WHERE order_data.user_id=? AND order_data.id=?`

	GET_UNPAID_ORDER = `SELECT o.id, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
		COUNT(*) item_count FROM order_detail od LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id
		LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.status_id=2 GROUP BY o.id`
	VERIFY_ORDER = `UPDATE order_data SET status_id=3 WHERE order_data.id=?`

	COUNT_TRX = fmt.Sprintf("SELECT COUNT(*) as trx_total, SUM(t.grand_total) FROM (%s) t WHERE(t.created_at BETWEEN ? AND ?)", HISTORY_METADATA)
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

func (repo *usercartRepositoryImpl) GetCartByUserID(ctx context.Context, usr int64) (res *models.OrderData, err error) {
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

func (repo *usercartRepositoryImpl) GetItemsByOrderID(ctx context.Context, orderID uint64) (res models.OrderDetails, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_ORDER_DETAIL_BY_ID)
	if err != nil {
		log.Printf("[GetItemsByOrderID] failed to prepare query, err => %+v\n", err)
		return
	}

	// log.Println(query)

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

	// log.Println(res)

	return
}

func (repo *usercartRepositoryImpl) GetCartMetadataByOrderID(ctx context.Context, orderID uint64) (res *models.OrderMetadata, err error) {
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

func (repo *usercartRepositoryImpl) GetItem(ctx context.Context, orderID uint64, productID uint64) (res *models.OrderDetail, err error) {
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
		log.Printf("[DeleteItem] failed to remove item, orderID => %d, pid => %d, err => %+v\n", orderID, productID, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[DeleteItem] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) GetHistoryMetadataByID(ctx context.Context, productID uint64) (*models.OrderHistoryMetadata, error) {
	query, err := repo.db.PrepareContext(ctx, GET_H_METADATA_BY_ID)
	if err != nil {
		log.Printf("[GetHistoryMetadataByID] failed to prepare query, err => %+v\n", err)
		return nil, err
	}

	rows, err := query.QueryContext(ctx,
		1,         // Cart
		2,         // Unpaid
		3,         // Paid
		productID, // Filter
	)
	if err != nil {
		log.Printf("[GetHistoryMetadataByID] failed to fetch order metadatas, err => %+v\n", err)
		return nil, err
	}

	res := &models.OrderHistoryMetadata{}
	rows.Next()
	err = rows.Scan(&res.OrderID, &res.UserID, &res.OrderDate, &res.GrandTotal, &res.ItemCount, &res.CouponName, &res.StatusName)
	if err != nil {
		log.Printf("[mapHistory] an error occured while parsing query result, err => %+v\n", err)
		return nil, err
	}

	return res, nil
}

func (repo *usercartRepositoryImpl) GetHistoryMetadataAll(ctx context.Context, params dto.HistoryParams) (meta []*models.OrderHistoryMetadata, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_H_METADATA_ALL)
	if err != nil {
		log.Printf("[GetHistoryMetadata] failed to prepare query, err => %+v\n", err)
		return
	}

	rows, err := query.QueryContext(ctx,
		1, 2, 3,
		params.DateStart, params.DateEnd, params.Limit, params.Offset*params.Limit, // Filter
	)

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

func (repo *usercartRepositoryImpl) GetHistoryMetadata(ctx context.Context, params dto.HistoryParams) (meta []*models.OrderHistoryMetadata, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_H_METADATA)
	if err != nil {
		log.Printf("[GetHistoryMetadata] failed to prepare query, err => %+v\n", err)
		return
	}

	rows, err := query.QueryContext(ctx,
		1, 2, 3,
		params.UserID, params.Limit, params.Offset*params.Limit,
	)

	if err != nil {
		log.Printf("[GetHistoryMetadata] failed to fetch order metadatas, err => %+v\n", err)
		return
	}

	meta, err = mapHistory(rows)
	if err != nil {
		log.Printf("[GetHistoryMetadata] failed to parse order metadatas, err => %+v\n", err)
		return
	}
	log.Println(meta, err)

	return
}

func (repo *usercartRepositoryImpl) GetCouponByCode(ctx context.Context, code string) (res *models.Coupon, err error) {
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

func (repo *usercartRepositoryImpl) VerifyOrder(ctx context.Context, orderID uint64) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[VerifyOrder] failed to start new transaction, err => %v\n", err)
		return
	}

	defer tx.Rollback()

	query, err := tx.PrepareContext(ctx, VERIFY_ORDER)
	if err != nil {
		log.Printf("[VerifyOrder] failed to prepare query, err => %+v\n", err)
		return
	}

	_, err = query.ExecContext(ctx, orderID)
	if err != nil {
		log.Printf("[VerifyOrder] failed to checkout an order, orderID => %d, err => %+v\n", orderID, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[VerifyOrder] failed to commit transaction, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) GetUnpaidOrder(ctx context.Context) (res []*models.OrderMetadata, err error) {
	query, err := repo.db.PrepareContext(ctx, GET_UNPAID_ORDER)
	if err != nil {
		log.Printf("[GetUnpaidOrder] failed to prepare query, err => %+v\n", err)
		return
	}

	rows, err := query.QueryContext(ctx)
	if err != nil {
		log.Printf("[GetUnpaidOrder] failed to fetch usercart, err => %+v\n", err)
		return
	}

	res, err = mapOrderMetadatas(rows)
	if err != nil {
		log.Printf("[GetUnpaidOrder] failed to fetch usercart, err => %+v\n", err)
		return
	}

	return
}

func (repo *usercartRepositoryImpl) GetStatistics(ctx context.Context, dateStart time.Time, dateEnd time.Time) (res *models.Statistics, err error) {
	row := repo.db.QueryRowContext(ctx, COUNT_TRX, 3, dateStart, dateEnd)

	res, err = mapStatistic(row)
	if err != nil {
		log.Printf("[GetStatistics] failed to fetch statistics, err => %+v\n", err)
	}

	return
}

func mapOrderData(r *sql.Row) (res *models.OrderData, err error) {
	res = &models.OrderData{}
	err = r.Scan(&res.ID, &res.UserID, &res.StatusID, &res.StatusName, &res.CouponID)
	return
}

func mapOrderDetails(r *sql.Rows) (res models.OrderDetails, err error) {
	res = models.OrderDetails{}

	for r.Next() {
		temp := &models.OrderDetail{}
		err = r.Scan(&temp.ID, &temp.OrderID, &temp.ProductID, &temp.ProductName, &temp.Price, &temp.Qty, &temp.Disc, &temp.SubTotal)
		if err != nil {
			log.Printf("[mapOrderDetails] an error occured while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
	}

	return
}

func mapOrderMetadata(r *sql.Row) (res *models.OrderMetadata, err error) {
	res = &models.OrderMetadata{}
	err = r.Scan(&res.OrderID, &res.GrandTotal, &res.ItemCount)
	return
}

func mapOrderMetadatas(r *sql.Rows) (res []*models.OrderMetadata, err error) {
	res = []*models.OrderMetadata{}

	for r.Next() {
		temp := &models.OrderMetadata{}
		err = r.Scan(&temp.OrderID, &temp.GrandTotal, &temp.ItemCount)
		if err != nil {
			log.Printf("[mapOrderMetadatas] an error occured while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
	}

	return
}

func mapOrderDetail(r *sql.Row) (res *models.OrderDetail, err error) {
	res = &models.OrderDetail{}
	err = r.Scan(&res.ID, &res.OrderID, &res.ProductID, &res.ProductName, &res.Price, &res.Qty, &res.Disc, &res.SubTotal)
	return
}

func mapHistory(r *sql.Rows) (res []*models.OrderHistoryMetadata, err error) {
	res = []*models.OrderHistoryMetadata{}

	count := 0
	for r.Next() {
		temp := &models.OrderHistoryMetadata{}
		err = r.Scan(&temp.OrderID, &temp.UserID, &temp.OrderDate, &temp.GrandTotal, &temp.ItemCount, &temp.CouponName, &temp.StatusName)
		if err != nil {
			log.Printf("[mapHistory] an error occured while parsing query result, err => %+v\n", err)
			return nil, err
		}
		res = append(res, temp)
		count++
	}

	if count == 0 {
		return nil, errors.ErrInvalidResources
	}

	return
}

func mapCoupon(r *sql.Row) (res *models.Coupon, err error) {
	res = &models.Coupon{}
	err = r.Scan(&res.ID, &res.Code, &res.Amount, &res.ExpiredAt)
	return
}

func mapStatistic(r *sql.Row) (res *models.Statistics, err error) {
	res = &models.Statistics{}
	err = r.Scan(&res.Count, &res.Income)
	return
}
