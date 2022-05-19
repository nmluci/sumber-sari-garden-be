package entity

type ProductParameter struct {
	Keyword    string
	Order      string
	Categories []string
	Limit      uint64
	Offset     uint64
}

type Product struct {
	Name        string
	PictureURL  string
	Description string
	ID          uint64
	CategoryID  uint64
	Price       uint64
	Qty         uint64
}

type ProductCategory struct {
	ID   uint64
	Name string
}

type ProductDetail struct {
	Name         string
	CategoryName string
	PictureURL   string
	Description  string
	ID           uint64
	CategoryID   uint64
	Price        uint64
	Qty          uint64
}

type ProductDetails []*ProductDetail

type Products []*Product

type ProductCategories []*ProductCategory
