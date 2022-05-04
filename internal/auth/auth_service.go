package auth

import (
	"context"
	"database/sql"

	"github.com/nmluci/sumber-sari-garden/internal/auth/impl"
	"github.com/nmluci/sumber-sari-garden/internal/dto"
)

type AuthService interface {
	RegisterNewUser(ctx context.Context, res dto.UserRegistrationRequest) (err error)
	LoginUser(ctx context.Context, res dto.UserSignIn) (data *dto.UserSignInResponse, err error)
}

func NewAuthService(db *sql.DB) AuthService {
	repo := impl.NewAuthRepository(db)
	return impl.NewAuthService(repo)
}
