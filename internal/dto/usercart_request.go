package dto

import (
	"time"

	"github.com/nmluci/sumber-sari-garden/internal/models"
)

type UpsertItemRequest struct {
	ProductID uint64 `json:"product_id"`
	Qty       uint64 `json:"qty"`
}

type HistoryParams struct {
	UserID    uint64
	Limit     uint64
	Offset    uint64
	DateStart time.Time
	DateEnd   time.Time
}

func (dto *UpsertItemRequest) ToEntity() (res *models.OrderDetail) {
	res = &models.OrderDetail{
		ProductID: dto.ProductID,
		Qty:       dto.Qty,
	}

	return
}
