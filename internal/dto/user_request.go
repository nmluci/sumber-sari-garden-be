package dto

import "github.com/nmluci/sumber-sari-garden/internal/models"

type UserRegistrationRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	RoleID    int64  `json:"role_id"`
}

type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *UserRegistrationRequest) ToEntity() (usr *models.UserInfo, cred *models.UserCred) {
	usr = &models.UserInfo{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Phone:     dto.Phone,
		Address:   dto.Address,
	}

	if dto.RoleID == 0 || dto.RoleID >= 2 {
		dto.RoleID = 2
	}

	cred = &models.UserCred{
		Email:    dto.Email,
		Password: dto.Password,
		UserRole: dto.RoleID,
	}

	return
}

func (dto *UserSignIn) ToEntity() (cred *models.UserCred) {
	return &models.UserCred{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
