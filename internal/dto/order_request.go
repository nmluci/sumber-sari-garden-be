package dto

import (
	"github.com/nmluci/sumber-sari-garden/internal/global/util/timeutil"
	"github.com/nmluci/sumber-sari-garden/internal/models"
)

type CheckoutItem struct {
	ProductID uint64 `json:"product_id"`
	Qty       uint64 `json:"qty"`
}

type OrderCheckoutRequest struct {
	CouponCode string          `json:"coupon_code"`
	Items      []*CheckoutItem `json:"items"`
}

func (dto *OrderCheckoutRequest) ToEntity() (res models.OrderDetails, coupon string) {
	res = models.OrderDetails{}

	for _, itm := range dto.Items {
		temp := &models.OrderDetail{
			ProductID: itm.ProductID,
			Qty:       itm.Qty,
		}
		res = append(res, temp)
	}

	return res, dto.CouponCode
}

func (dto *Coupon) ToEntity() (res *models.Coupon) {
	res = &models.Coupon{
		Code:        dto.Code,
		Amount:      dto.Amount,
		Description: dto.Description,
	}

	res.ExpiredAt, _ = timeutil.ParseLocalTime(dto.ExpiredAt, "2006-01-02 15:04:05")
	return
}
