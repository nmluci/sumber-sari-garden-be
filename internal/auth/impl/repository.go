package impl

import (
	"context"
	"database/sql"

	"github.com/nmluci/sumber-sari-garden/internal/entity"
)

type AuthRepository interface {
	StoreUserInfo(ctx context.Context, data entity.UserInfo) (err error)
	StoreUserCred(ctx context.Context, data entity.UserCred) (err error)
	GetCredByUsername(ctx context.Context, username string) (usr entity.UserCred, err error)
	GetUserInfoByID(ctx context.Context, userID int64) (usr entity.UserInfo, err error)
}

type authRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *authRepositoryImpl {
	return &authRepositoryImpl{db: db}
}
