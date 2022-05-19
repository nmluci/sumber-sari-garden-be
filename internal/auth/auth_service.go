package auth

import (
	"context"

	"github.com/nmluci/sumber-sari-garden/internal/auth/impl"
	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/entity"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
)

type AuthService interface {
	RegisterNewUser(ctx context.Context, res *dto.UserRegistrationRequest) (err error)
	LoginUser(ctx context.Context, res *dto.UserSignIn) (data *dto.UserSignInResponse, err error)
	FindUserByAccessToken(ctx context.Context, accessToken string) (*entity.UserCred, error)
}

func NewAuthService(db *database.DatabaseClient) AuthService {
	repo := impl.NewAuthRepository(db)
	return impl.NewAuthService(repo)
}
