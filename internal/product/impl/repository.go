package impl

import (
	"database/sql"

	"github.com/nmluci/sumber-sari-garden/pkg/database"
)

type ProductRepository interface {
}

type productRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *database.DatabaseClient) *productRepositoryImpl {
	return &productRepositoryImpl{db: db.DB}
}
