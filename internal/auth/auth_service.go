package auth

import (
	"database/sql"

	"github.com/nmluci/sumber-sari-garden/internal/auth/impl"
)

type AuthService interface {
}

func NewAuthService(db *sql.DB) AuthService {
	repo := impl.NewAuthRepository(db)
	return impl.NewAuthService(repo)
}
