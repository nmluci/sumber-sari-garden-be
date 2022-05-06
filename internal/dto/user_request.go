package dto

import (
	"github.com/nmluci/sumber-sari-garden/internal/entity"
)

type UserRegistrationRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *UserRegistrationRequest) ToEntity() (usr *entity.UserInfo, cred *entity.UserCred) {
	usr = &entity.UserInfo{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Phone:     dto.Phone,
		Address:   dto.Address,
	}

	cred = &entity.UserCred{
		Email:    dto.Email,
		Password: dto.Password,
	}

	return
}

func (dto *UserSignIn) ToEntity() (cred *entity.UserCred) {
	return &entity.UserCred{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
