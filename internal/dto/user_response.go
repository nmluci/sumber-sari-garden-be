package dto

import (
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/models"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type UserSignInResponse struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	RoleID      int64  `json:"role_id"`
	AccessToken string `json:"access_token"`
}

type UserProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	RoleID    int64  `json:"role_id"`
}

func NewUserSignInResponse(cred *models.UserCred, info *models.UserInfo, ac string) (res *UserSignInResponse, err error) {
	if cred == nil || info == nil {
		log.Printf("[NewUserSignInResponse] failed to encode response data due to inconsisted data")
		err = errors.ErrInvalidResources
		return
	}

	res = &UserSignInResponse{
		FirstName:   info.FirstName,
		LastName:    info.LastName,
		Email:       cred.Email,
		RoleID:      cred.UserRole,
		AccessToken: ac,
	}

	return
}

func NewUserProfileResponse(cred *models.UserCred, info *models.UserInfo) (res *UserProfileResponse, err error) {
	if cred == nil || info == nil {
		log.Printf("[NewUserProfileResponse] failed to encode response data due to inconsistent data")
		err = errors.ErrInvalidResources
		return
	}

	res = &UserProfileResponse{
		FirstName: info.FirstName,
		LastName:  info.LastName,
		Phone:     info.Phone,
		Address:   info.Address,
		Email:     cred.Email,
		Password:  "Password Hidden",
		RoleID:    cred.UserRole,
	}

	return
}
