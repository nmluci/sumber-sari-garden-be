package dto

import "github.com/nmluci/sumber-sari-garden/internal/entity"

type CheckoutItem struct {
	ProductID uint64 `json:"product_id"`
	Qty       uint64 `json:"qty"`
}

type OrderCheckoutRequest struct {
	CouponCode string          `json:"coupon_code"`
	Items      []*CheckoutItem `json:"items"`
}

func (dto *OrderCheckoutRequest) ToEntity() (res entity.OrderDetails, coupon string) {
	res = entity.OrderDetails{}

	for _, itm := range dto.Items {
		temp := &entity.OrderDetail{
			ProductID: itm.ProductID,
			Qty:       itm.Qty,
		}
		res = append(res, temp)
	}

	return res, dto.CouponCode
}
