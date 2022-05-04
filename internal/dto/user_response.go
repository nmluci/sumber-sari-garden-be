package dto

import (
	"errors"

	"github.com/nmluci/sumber-sari-garden/internal/entity"
)

type UserSignInResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	RoleID    int64  `json:"role_id"`
}

func NewUserSignInResponse(cred *entity.UserCred, info *entity.UserInfo) (res *UserSignInResponse, err error) {
	if cred != nil || info != nil {
		err = errors.New("user has incompleted data")
		return
	}

	res = &UserSignInResponse{
		FirstName: info.FirstName,
		LastName:  info.LastName,
		Email:     cred.Email,
		RoleID:    *cred.UserRole,
	}

	return
}
