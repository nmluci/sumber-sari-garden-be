package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/entity"
)

type AuthRepository interface {
	StoreUserInfo(ctx context.Context, data *entity.UserInfo) (userID int64, err error)
	StoreUserCred(ctx context.Context, data *entity.UserCred) (err error)
	GetCredByEmail(ctx context.Context, email string) (usr *entity.UserCred, err error)
	GetUserInfoByID(ctx context.Context, userID int64) (usr *entity.UserInfo, err error)
}

type authRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *authRepositoryImpl {
	return &authRepositoryImpl{db: db}
}

const (
	STORE_USER_INFO    = `INSERT INTO user_info(first_name, last_name, phone, address) VALUES (?, ?, ?, ?);`
	STORE_USER_CRED    = `INSERT INTO user(user_id, email, password) VALUES (?, ?, ?);`
	FIND_CRED_BY_EMAIL = `SELECT user_id, email, password, role_id FROM user WHERE email=?;`
	FIND_USER_BY_ID    = `SELECT first_name, last_name, phone, address FROM user WHERE id=?;`
)

func (auth *authRepositoryImpl) StoreUserInfo(ctx context.Context, data *entity.UserInfo) (userID int64, err error) {
	res, err := auth.db.ExecContext(ctx, STORE_USER_INFO, data.FirstName, data.LastName, data.Phone, data.Address)
	if err != nil {
		log.Printf("[StoreUserInfo] err: %v\n", err)
		return
	}

	userID, err = res.LastInsertId()
	if err != nil {
		log.Printf("[StoreUserInfo] err: %v\n", err)
		return 0, err
	}

	return
}

func (auth *authRepositoryImpl) StoreUserCred(ctx context.Context, data *entity.UserCred) (err error) {
	res, err := auth.db.ExecContext(ctx, STORE_USER_CRED, data.UserID, data.Email, data.Password)
	if err != nil {
		log.Printf("[StoreUserCred] err: %v\n", err)
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Printf("[StoreUserCred] err: %v\n", err)
		return
	}

	return
}

func (auth *authRepositoryImpl) GetCredByEmail(ctx context.Context, email string) (usr *entity.UserCred, err error) {
	res := auth.db.QueryRowContext(ctx, FIND_CRED_BY_EMAIL, email)
	usr, err = mapCredToEntity(res)
	if err != nil {
		log.Printf("[GetCredByEmail] err: %v\n", err)
		return
	}

	return
}

func (auth *authRepositoryImpl) GetUserInfoByID(ctx context.Context, userID int64) (usr *entity.UserInfo, err error) {
	res := auth.db.QueryRowContext(ctx, FIND_USER_BY_ID, userID)
	usr, err = mapUserInfoToEntity(res)
	if err != nil {
		log.Printf("[GetUserInfoByID] err: %v\n", err)
		return
	}

	return
}

func mapCredToEntity(row *sql.Row) (usr *entity.UserCred, err error) {
	err = row.Scan(&usr.UserID, &usr.Email, &usr.Password, &usr.UserRole)
	return
}

func mapUserInfoToEntity(row *sql.Row) (usr *entity.UserInfo, err error) {
	err = row.Scan(&usr.FirstName, &usr.LastName, &usr.Phone, &usr.Address)
	return
}
