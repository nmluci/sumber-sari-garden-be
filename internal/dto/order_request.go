package dto

import "github.com/nmluci/sumber-sari-garden/internal/models"

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
