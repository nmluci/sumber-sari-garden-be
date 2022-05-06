package impl

import (
	"context"
	"errors"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{repo: repo}
}

func (auth AuthServiceImpl) RegisterNewUser(ctx context.Context, res *dto.UserRegistrationRequest) (err error) {
	userInfo, userCred := res.ToEntity()

	hashed, err := bcrypt.GenerateFromPassword([]byte(userCred.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[RegisterNewUser] failed to generate new password, err => %+v", err)
		return
	}

	userCred.Password = string(hashed)
	userID, err := auth.repo.StoreUserCred(ctx, *userCred)
	if err != nil {
		return
	}

	userInfo.UserID = userID
	err = auth.repo.StoreUserInfo(ctx, *userInfo)

	return
}

func (auth AuthServiceImpl) LoginUser(ctx context.Context, res *dto.UserSignIn) (data *dto.UserSignInResponse, err error) {
	cred := res.ToEntity()

	userCred, err := auth.repo.GetCredByEmail(ctx, cred.Email)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userCred.Password), []byte(cred.Password))
	if err != nil {
		err = errors.New("invalid user credentials")
		return
	}

	userInfo, err := auth.repo.GetUserInfoByID(ctx, userCred.UserID)
	if err != nil {
		return
	}

	return dto.NewUserSignInResponse(userCred, userInfo)
}
