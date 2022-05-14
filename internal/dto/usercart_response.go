package dto

import (
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/entity"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type CartItemResponse struct {
	ProductID   uint64  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       uint64  `json:"price"`
	Qty         uint64  `json:"qty"`
	Disc        float32 `json:"disc"`
	Subtotal    float32 `json:"sub_total"`
}

type UsercartResponse struct {
	OrderID    uint64              `json:"order_id"`
	UserID     uint64              `json:"user_id"`
	Status     string              `json:"order_status"`
	ItemCount  uint64              `json:"item_count"`
	GrandTotal float32             `json:"grand_total"`
	Items      []*CartItemResponse `json:"items"`
}

func NewUsercartResponse(meta *entity.OrderMetadata, data *entity.OrderData, items entity.OrderDetails) (res UsercartResponse, err error) {
	if meta == nil || data == nil || len(items) == 0 {
		log.Printf("[NewUsercartResponse] failed to encode response data due to inconsisted data")
		err = errors.ErrInvalidResources
		return
	}

	res = UsercartResponse{
		OrderID:    data.ID,
		UserID:     data.UserID,
		Status:     data.StatusName,
		ItemCount:  meta.ItemCount,
		GrandTotal: meta.GrandTotal,
	}

	for _, itm := range items {
		temp := &CartItemResponse{
			ProductID:   itm.ProductID,
			ProductName: itm.ProductName,
			Price:       itm.Price,
			Qty:         itm.Qty,
			Disc:        itm.Disc,
			Subtotal:    itm.SubTotal,
		}

		res.Items = append(res.Items, temp)
	}

	return
}
