package models

import "time"

type OrderData struct {
	ID         uint64
	UserID     uint64
	StatusID   uint64
	CouponID   *uint64
	StatusName string
}

type OrderDetail struct {
	ID          uint64
	OrderID     uint64
	ProductID   uint64
	ProductName string
	Price       uint64
	Qty         uint64
	Disc        float32
	SubTotal    float32
}

type OrderMetadata struct {
	OrderID    uint64
	GrandTotal float32
	ItemCount  uint64
}

type OrderHistoryMetadata struct {
	OrderID    uint64
	UserID     uint64
	OrderDate  time.Time
	GrandTotal float32
	ItemCount  uint64
	CouponName *string
	StatusName string
}

type ActiveCoupon struct {
	ID          uint64
	Code        string
	Amount      float32
	Description *string
	ExpiredAt   time.Time
}

type OrderDetails []*OrderDetail

type ActiveCoupons []*ActiveCoupon
