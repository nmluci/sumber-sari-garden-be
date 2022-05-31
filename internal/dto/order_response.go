package dto

import (
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/global/util/timeutil"
	"github.com/nmluci/sumber-sari-garden/internal/models"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type TrxItem struct {
	ProductID   uint64  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       uint64  `json:"price"`
	Qty         uint64  `json:"qty"`
	Subtotal    float32 `json:"sub_total"`
}

type TrxMetadata struct {
	UserId     uint64     `json:"user_id"`
	OrderID    uint64     `json:"order_id"`
	OrderDate  string     `json:"order_date"`
	Status     string     `json:"order_status"`
	ItemCount  uint64     `json:"item_count"`
	GrandTotal float32    `json:"grand_total"`
	Coupon     *string    `json:"coupon_code"`
	Items      []*TrxItem `json:"items"`
}

type TrxBrief struct {
	OrderID    uint64  `json:"order_id"`
	GrandTotal float32 `json:"grand_total"`
	ItemCount  uint64  `json:"item_count"`
}

type OrderHistoryResponse struct {
	UserID uint64         `json:"user_id"`
	Trx    []*TrxMetadata `json:"trx"`
}

type Coupon struct {
	ID          uint64  `json:"coupon_id"`
	Code        string  `json:"code"`
	Amount      float32 `json:"amount"`
	ExpiredAt   string  `json:"expired_at"`
	Description *string `json:"description"`
}

type OrderHistoriesResponse []*TrxMetadata

type Coupons []*Coupon

func NewOrderHistoriesResponse(meta []*models.OrderHistoryMetadata, items models.OrderDetails) (res OrderHistoriesResponse, err error) {
	res = []*TrxMetadata{}
	for mi, m := range meta {
		trx := &TrxMetadata{
			OrderID:    m.OrderID,
			UserId:     m.UserID,
			OrderDate:  timeutil.FormatLocalTime(m.OrderDate, "2006-01-02 15:04:05"),
			Status:     m.StatusName,
			ItemCount:  m.ItemCount,
			GrandTotal: m.GrandTotal,
			Coupon:     m.CouponName,
		}

		for ti, t := range items {
			if t.OrderID != m.OrderID {
				continue
			}

			data := &TrxItem{
				ProductID:   t.ProductID,
				ProductName: t.ProductName,
				Price:       t.Price,
				Qty:         t.Qty,
				Subtotal:    t.SubTotal,
			}
			trx.Items = append(trx.Items, data)

			if len(items) >= 3 {
				items = append(items[ti:], items[:ti+1]...)
			}
		}

		if len(meta) >= 3 {
			meta = append(meta[mi:], meta[:mi+1]...)
		}

		res = append(res, trx)
	}

	return

}

func NewOrderHistoryResponse(userID uint64, meta []*models.OrderHistoryMetadata, items models.OrderDetails) (res *OrderHistoryResponse, err error) {
	if len(meta) == 0 || len(items) == 0 {
		log.Printf("[NewOrderHistoryResponse] failed to encode response data due to incomplete data\n")
		err = errors.ErrInvalidResources
		return
	}

	res = &OrderHistoryResponse{
		UserID: userID,
	}

	for mi, m := range meta {
		trx := &TrxMetadata{
			OrderID:    m.OrderID,
			OrderDate:  timeutil.FormatLocalTime(m.OrderDate, "2006-01-02 15:04:05"),
			Status:     m.StatusName,
			ItemCount:  m.ItemCount,
			GrandTotal: m.GrandTotal,
			Coupon:     m.CouponName,
		}

		for ti, t := range items {
			if t.OrderID != m.OrderID {
				continue
			}

			data := &TrxItem{
				ProductID:   t.ProductID,
				ProductName: t.ProductName,
				Price:       t.Price,
				Qty:         t.Qty,
				Subtotal:    t.SubTotal,
			}
			trx.Items = append(trx.Items, data)

			if len(items) >= 3 {
				items = append(items[ti:], items[:ti+1]...)
			}
		}

		if len(meta) >= 3 {
			meta = append(meta[mi:], meta[:mi+1]...)
		}

		res.Trx = append(res.Trx, trx)
	}

	return
}

func NewCouponResponse(coupons models.Coupons, isAdmin bool) (res Coupons, err error) {
	if len(coupons) == 0 {
		log.Printf("[NewActiveCouponResponse] failed to encode response data due to incomplete data\n")
		err = errors.ErrInvalidResources
		return
	}

	res = Coupons{}

	for _, c := range coupons {
		temp := &Coupon{
			Code:        c.Code,
			Amount:      c.Amount,
			ExpiredAt:   timeutil.FormatLocalTime(c.ExpiredAt, "2006-01-02 15:04:05"),
			Description: c.Description,
		}

		if isAdmin {
			temp.ID = c.ID
		}

		res = append(res, temp)
	}

	return
}

func NewTrxBrief(orders []*models.OrderMetadata) (res []*TrxBrief, err error) {
	if len(orders) == 0 {
		log.Printf("[NewTrxBrief] failed to encode response data due to incomplete data\n")
		err = errors.ErrInvalidResources
		return
	}

	res = []*TrxBrief{}

	for _, o := range orders {
		temp := &TrxBrief{
			OrderID:    o.OrderID,
			GrandTotal: o.GrandTotal,
			ItemCount:  o.ItemCount,
		}

		res = append(res, temp)
	}

	return
}
