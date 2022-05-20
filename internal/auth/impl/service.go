package impl

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nmluci/sumber-sari-garden/internal/dto"
	"github.com/nmluci/sumber-sari-garden/internal/global/config"
	"github.com/nmluci/sumber-sari-garden/internal/models"
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

	token, err := auth.newAccessToken(userCred)
	if err != nil {
		log.Printf("[LoginUser] failed to generate access token, err => %+v\n", err)
		return
	}

	return dto.NewUserSignInResponse(userCred, userInfo, token)
}

func (auth AuthServiceImpl) FindUserByAccessToken(ctx context.Context, accessToken string) (*models.UserCred, error) {
	config := config.GetConfig()

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != config.JWT_SIGNING_METHOD {
			return nil, fmt.Errorf("signin method invalid")
		}
		return config.JWT_SIGNATURE_KEY, nil
	})

	if err != nil {
		log.Printf("[FindUserByAccessToken] error: %v\n", err)
		return nil, errors.ErrTokenExpired
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Printf("[FindUserByAccessToken] error: %v\n", err)
		return nil, err
	}

	data := claims["data"].(map[string]interface{})
	userId := int64(data["id"].(float64))

	user, err := auth.repo.GetCredByID(ctx, userId)
	if err != nil {
		log.Printf("[FindUserByAccessToken] error: %v\n", err)
		return nil, errors.ErrUserCredential
	}

	return user, nil
}

func (auth AuthServiceImpl) newAccessToken(user *models.UserCred) (string, error) {
	config := config.GetConfig()

	claims := auth.newUserClaim(user.UserID, user.Email, config.JWT_AT_EXPIRATION)
	accessToken := jwt.NewWithClaims(config.JWT_SIGNING_METHOD, claims)
	signed, err := accessToken.SignedString(config.JWT_SIGNATURE_KEY)
	if err != nil {
		log.Printf("[NewAccessToken] error: %v\n", err)
		return "", err
	}

	return signed, nil
}

func (auth AuthServiceImpl) newUserClaim(id int64, username string, exp time.Duration) *jwt.MapClaims {
	return &jwt.MapClaims{
		"iss": config.GetConfig().JWT_ISSUER,
		"exp": time.Now().Add(exp).Unix(),
		"data": map[string]interface{}{
			"id":       id,
			"username": username,
		},
	}
}
