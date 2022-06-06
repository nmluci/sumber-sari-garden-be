package impl

import (
	"context"
	"log"
	"time"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/authutil"
	"github.com/nmluci/sumber-sari-garden/internal/global/util/timeutil"
	"github.com/nmluci/sumber-sari-garden/internal/models"
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

func (us *UsercartServiceImpl) Checkout(ctx context.Context, req *dto.OrderCheckoutRequest) (res *dto.CheckoutResponse, err error) {
	var (
		couponID        = new(uint64)
		orderID  uint64 = 0
	)

	usr := authutil.GetUserIDFromCtx(ctx)
	data, coupon := req.ToEntity()

	orderInfo, err := us.repo.GetCartByUserID(ctx, usr)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[Checkout] an error occured while validating user's cart, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources {
		log.Printf("[Checkout] user doesn't hawve any order right now\n")
		return
	} else {
		orderID = orderInfo.ID
	}

	items, err := us.repo.GetItemsByOrderID(ctx, orderID)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[Checkout] an error occured while fetching order's items, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources {
		log.Printf("[Checkout] user doesn't hawve any order right now\n")
		return
	}

	couponInfo, err := us.repo.GetCouponByCode(ctx, coupon)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[Checkout] an error occured while fetching coupon's info, err => %+v\n", err)
		return
	} else if err == errors.ErrInvalidResources {
		log.Printf("[Checkout] coupon aren't active or not existed\n")
	} else {
		*couponID = couponInfo.ID
	}
	log.Println(err)

	validItem := models.OrderDetails{}
	for _, itm := range data {
		for _, c := range items {
			if itm.ProductID == c.ProductID && itm.Qty == c.Qty {
				validItem = append(validItem, itm)
				break
			}
		}
	}

	if len(validItem) != len(items) {
		newID, err2 := us.repo.NewCart(ctx, usr)
		if err2 != nil {
			log.Printf("[Checkout] an error occured while making new cart, err => %+v\n", err2)
			return
		}

		orderID = uint64(newID)
		for _, itm := range items {
			err = us.repo.InsertItem(ctx, orderID, itm.ProductID, itm.Qty)
			if err != nil {
				log.Printf("[Checkout] an error occured while inserting items, productID=%d, err => %+v\n", itm.ProductID, err)
				return
			}
		}
	}

	err = us.repo.Checkout(ctx, usr, orderID, couponID)
	if err != nil {
		log.Printf("[Checkout] an error occured while checkouting order, orderID=%d, err => %+v\n", orderInfo.ID, err)
		return
	}

	return dto.NewCheckoutResponse(orderID, time.Now())
}

func (us *UsercartServiceImpl) SpecificOrderHistoryById(ctx context.Context, productID uint64) (res *dto.TrxMetadata, err error) {
	usrID := authutil.GetUserIDFromCtx(ctx)
	if priv := authutil.GetUserPrivFromCtx(ctx); priv != 1 {
		log.Printf("[SpecificOrderHistoryById] user doesn't have enough permission, user_id => %d\n", usrID)
		err = errors.ErrUserPriv
		return
	}

	meta, err := us.repo.GetHistoryMetadataByID(ctx, productID)
	if err != nil {
		log.Printf("[OrderHistory] an error occured while fetching histories' metadata, err => %+v\n", err)
		return
	}

	orderInfo, err := us.repo.GetItemsByOrderID(ctx, productID)
	if err != nil {
		log.Printf("[SpecificOrderHistoryById] an error occured while fetching order's item, orderID => %d, err => %+v\n", productID, err)
		return nil, err
	}

	return dto.NewOrderDetailsById(meta, orderInfo)
}

func (us *UsercartServiceImpl) OrderHistory(ctx context.Context, params dto.HistoryParams) (res *dto.OrderHistoryResponse, err error) {
	usrID := authutil.GetUserIDFromCtx(ctx)
	params.UserID = uint64(usrID)

	meta, err := us.repo.GetHistoryMetadata(ctx, params)
	if err != nil {
		log.Printf("[OrderHistory] an error occured while fetching histories' metadata, err => %+v\n", err)
		return
	}

	log.Println(meta, err)

	items := models.OrderDetails{}
	for _, itm := range meta {
		orderInfo, err2 := us.repo.GetItemsByOrderID(ctx, itm.OrderID)
		if err2 != nil {
			log.Printf("[OrderHistory] an error occured while fetching order's item, orderID => %d, err => %+v\n", itm.OrderID, err2)
			return res, err2
		}

		items = append(items, orderInfo...)
	}

	return dto.NewOrderHistoryResponse(params.UserID, meta, items)
}

func (us *UsercartServiceImpl) VerifyOrder(ctx context.Context, orderID uint64) (err error) {
	usrID := authutil.GetUserIDFromCtx(ctx)
	if priv := authutil.GetUserPrivFromCtx(ctx); priv != 1 {
		log.Printf("[VerifyOrder] user doesn't have enough permission, user_id => %d\n", usrID)
		err = errors.ErrUserPriv
		return
	}

	if exists, err := us.repo.GetCartMetadataByOrderID(ctx, orderID); err != nil {
		log.Printf("[VerifyOrder] an error occured while verifying order, orderID => %d, err => %+v\n", orderID, err)
		return err
	} else if exists == nil {
		log.Printf("[VerifyOrder] order doesn't existed, orderID => %d, err => %+v\n", orderID, err)
		return errors.ErrInvalidResources
	}

	err = us.repo.VerifyOrder(ctx, orderID)
	if err != nil {
		log.Printf("[VerifyOrder] an error occured while verifying order, orderID => %d, err => %+v\n", orderID, err)
		return
	}

	return
}

func (us *UsercartServiceImpl) GetUnpaidOrder(ctx context.Context) (res []*dto.TrxBrief, err error) {
	usrID := authutil.GetUserIDFromCtx(ctx)
	if priv := authutil.GetUserPrivFromCtx(ctx); priv != 1 {
		log.Printf("[VerifyOrder] user doesn't have enough permission, user_id => %d\n", usrID)
		err = errors.ErrUserPriv
		return
	}

	orders, err := us.repo.GetUnpaidOrder(ctx)
	if err != nil {
		log.Printf("[VerifyOrder] an error occured while fetching unpaid orders, err => %+v\n", err)
		return
	}

	return dto.NewTrxBrief(orders)
}

func (us *UsercartServiceImpl) OrderHistoryAll(ctx context.Context, params dto.HistoryParams) (res dto.OrderHistoriesResponse, err error) {
	usrID := authutil.GetUserIDFromCtx(ctx)
	if priv := authutil.GetUserPrivFromCtx(ctx); priv != 1 {
		log.Printf("[VerifyOrder] user doesn't have enough permission, user_id => %d\n", usrID)
		err = errors.ErrUserPriv
		return
	}

	meta, err := us.repo.GetHistoryMetadataAll(ctx, params)
	if err != nil {
		log.Printf("[OrderHistory] an error occured while fetching histories' metadata, err => %+v\n", err)
		return
	}

	items := models.OrderDetails{}
	for _, itm := range meta {
		orderInfo, err := us.repo.GetItemsByOrderID(ctx, itm.OrderID)
		if err != nil {
			log.Printf("[OrderHistory] an error occured while fetching order's item, orderID => %d, err => %+v\n", itm.OrderID, err)
			return res, err
		}

		items = append(items, orderInfo...)
	}

	return dto.NewOrderHistoriesResponse(meta, items)
}

func (us *UsercartServiceImpl) GetStatistics(ctx context.Context) (res *dto.StatisticsResponse, err error) {
	usrID := authutil.GetUserIDFromCtx(ctx)
	if priv := authutil.GetUserPrivFromCtx(ctx); priv != 1 {
		log.Printf("[VerifyOrder] user doesn't have enough permission, user_id => %d\n", usrID)
		err = errors.ErrUserPriv
		return
	}

	now := time.Now()
	dailyStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dailyEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, time.UTC)

	dailyData, err := us.repo.GetStatistics(ctx, dailyStart, dailyEnd)
	if err != nil {
		log.Printf("[GetStatistics] an error occured while fetching daily insight, err => %+v\n", err)
		return
	} else {
		dailyData.DateStart = dailyStart
		dailyData.DateEnd = dailyEnd
	}

	weeklyData, err := us.repo.GetStatistics(ctx, (dailyStart.Add(-7 * 24 * time.Hour)), (dailyEnd))
	if err != nil {
		log.Printf("[GetStatistics] an error occured while fetching weekly insight, err => %+v\n", err)
		return
	} else {
		weeklyData.DateStart = dailyStart.Add(-7 * 24 * time.Hour)
		weeklyData.DateEnd = dailyEnd
	}

	monthData, err := us.repo.GetStatistics(ctx, timeutil.FirstDayOfMonth(dailyStart), timeutil.LastDayOfMonth(dailyStart))
	if err != nil {
		log.Printf("[GetStatistics] an error occured while fetching monthly insight, err => %+v\n", err)
		return
	} else {
		monthData.DateStart = timeutil.FirstDayOfMonth(dailyStart)
		monthData.DateEnd = timeutil.LastDayOfMonth(dailyStart)
	}

	annualStart := time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	annualEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, time.UTC)
	annualData, err := us.repo.GetStatistics(ctx, annualStart, annualEnd)
	if err != nil {
		log.Printf("[GetStatistics] an error occured while fetching annual insight, err => %+v\n", err)
		return
	} else {
		annualData.DateStart = annualStart
		annualData.DateEnd = annualEnd
	}

	return dto.NewOrderStatistics(dailyData, weeklyData, monthData, annualData)
}
