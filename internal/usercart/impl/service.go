package impl

import (
	"context"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/authutil"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type UsercartServiceImpl struct {
	repo UsercartRepository
}

func NewUsercartService(repo UsercartRepository) *UsercartServiceImpl {
	return &UsercartServiceImpl{repo: repo}
}

func (us *UsercartServiceImpl) UpsertItem(ctx context.Context, res *dto.UpsertItemRequest) (err error) {
	var cartID int64
	usr := authutil.GetUserIDFromCtx(ctx)
	data := res.ToEntity()

	carts, err := us.repo.GetCartByUserID(ctx, usr)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[UpsertItem] an error occured while validating user's cart, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources || carts.StatusID != 1 {
		cartID, err = us.repo.NewCart(ctx, usr)
		if err != nil {
			log.Printf("[UpsertItem] an error occured while creating new order, err => %+v\n", err)
			return
		}
	} else {
		cartID = int64(carts.ID)
	}

	item, err := us.repo.GetItem(ctx, uint64(cartID), data.ProductID)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[UpsertItem] an error occured while validating product, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources {
		err = us.repo.InsertItem(ctx, uint64(cartID), data.ProductID, data.Qty)
		if err != nil {
			log.Printf("[UpsertItem] an error occured while inserting item, err => %+v\n", err)
			return
		}
	} else {
		err = us.repo.UpdateItem(ctx, uint64(cartID), data.ProductID, data.Qty)
		if err != nil {
			log.Printf("[UpsertItem] an error occured while updating item, item_id => %d, err => %+v\n", item.ID, err)
			return
		}
	}

	return
}

func (us *UsercartServiceImpl) RemoveItem(ctx context.Context, pid uint64) (err error) {
	usr := authutil.GetUserIDFromCtx(ctx)

	cart, err := us.repo.GetCartByUserID(ctx, usr)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[RemoveItem] an error occured while validating user's cart, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources || cart.StatusID != 1 {
		log.Printf("[RemoveItem] user doesn't have any order right now\n")
		return
	}

	_, err = us.repo.GetItem(ctx, cart.ID, pid)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[RemoveItem] an error occured while validating product, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources {
		log.Printf("[RemoveItem] user doesn't have requested product right now, item_id => %d\n", pid)
		return
	}

	err = us.repo.RemoveItem(ctx, cart.ID, pid)
	if err != nil {
		log.Printf("[RemoveItem] an error occured while removing product from order, err => %+v\n", err)
		return
	}

	return
}

func (us *UsercartServiceImpl) GetCart(ctx context.Context) (cart dto.UsercartResponse, err error) {
	usr := authutil.GetUserIDFromCtx(ctx)

	orderInfo, err := us.repo.GetCartByUserID(ctx, usr)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[GetCart] an error occured while validating user's cart, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources {
		log.Printf("[GetCart] user doesn't have any order right now\n")
		return
	}

	items, err := us.repo.GetItemsByOrderID(ctx, orderInfo.ID)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[GetCart] an error occured while fetching order's items, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources {
		log.Printf("[GetCart] user doesn't have any order right now\n")
		return
	}

	meta, err := us.repo.GetCartMetadataByOrderID(ctx, orderInfo.ID)
	if err != nil {
		log.Printf("[GetCart] an error occured while fetching order's metadata, err => %+v\n", err)
		return
	}

	return dto.NewUsercartResponse(meta, orderInfo, items)
}

func (us *UsercartServiceImpl) Checkout(ctx context.Context, dto *dto.OrderCheckoutRequest) (err error) {

	return
}

func (us *UsercartServiceImpl) OrderHistory(ctx context.Context) (err error) {

	return
}
