package usercart

import (
	"context"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/usercart/impl"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
)

type UsercartService interface {
	UpsertItem(ctx context.Context, res *dto.UpsertItemRequest) (err error)
	RemoveItem(ctx context.Context, pid uint64) (err error)
	GetCart(ctx context.Context) (cart dto.UsercartResponse, err error)
	Checkout(ctx context.Context, res *dto.OrderCheckoutRequest) (err error)
	OrderHistory(ctx context.Context, params dto.HistoryParams) (data *dto.OrderHistoryResponse, err error)
	VerifyOrder(ctx context.Context, orderID uint64) (err error)
	GetUnpaidOrder(ctx context.Context) (res []*dto.TrxBrief, err error)
	OrderHistoryAll(ctx context.Context, params dto.HistoryParams) (res map[uint64][]*dto.TrxMetadata, err error)
}

func NewUsercartService(db *database.DatabaseClient) UsercartService {
	repo := impl.NewUsercartRepository(db)
	return impl.NewUsercartService(repo)
}
