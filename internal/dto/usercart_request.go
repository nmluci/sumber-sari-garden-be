package dto

import "github.com/nmluci/sumber-sari-garden/internal/entity"

type UpsertItemRequest struct {
	ProductID uint64 `json:"product_id"`
	Qty       uint64 `json:"qty"`
}

func (dto *UpsertItemRequest) ToEntity() (res *entity.OrderDetail) {
	res = &entity.OrderDetail{
		ProductID: dto.ProductID,
		Qty:       dto.Qty,
	}

	return
}
