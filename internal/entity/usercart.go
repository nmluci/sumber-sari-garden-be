package entity

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

type OrderDetails []*OrderDetail
