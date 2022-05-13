package impl

import (
	"context"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
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

	existed, err := auth.repo.GetCredByEmail(ctx, userCred.Email)
	if err != nil && err != errors.ErrInvalidResources {
		log.Printf("[RegisterNewUser] failed to check email duplication, err => %+v\n", err)
		return
	}

	if existed != nil {
		err = errors.ErrUserAlreadyExist
		log.Printf("[RegisterNewUser] user with email: %s already existed\n", userCred.Email)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userCred.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[RegisterNewUser] failed to generate new password, err => %+v\n", err)
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
		log.Printf("[LoginUser] failed to fetch user with email: %s, err => %+v\n", cred.Email, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userCred.Password), []byte(cred.Password))
	if err != nil {
		log.Printf("[LoginUser] failed to validate user cred with email: %s, err => %+v\n", cred.Email, err)
		return
	}

	userInfo, err := auth.repo.GetUserInfoByID(ctx, userCred.UserID)
	if err != nil {
		log.Printf("[LoginUser] failed to fetch user info with email: %s, err => %+v\n", cred.Email, err)
		return
	}

	return dto.NewUserSignInResponse(userCred, userInfo)
}
