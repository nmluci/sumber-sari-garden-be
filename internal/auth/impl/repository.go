package impl

import (
	"context"
	"database/sql"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/models"
	"github.com/nmluci/sumber-sari-garden/pkg/database"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

type AuthRepository interface {
	StoreUserInfo(ctx context.Context, data models.UserInfo) (err error)
	StoreUserCred(ctx context.Context, data models.UserCred) (userID int64, err error)
	GetCredByEmail(ctx context.Context, email string) (usr *models.UserCred, err error)
	GetUserInfoByID(ctx context.Context, userID int64) (usr *models.UserInfo, err error)
	GetCredByID(ctx context.Context, userID int64) (usr *models.UserCred, err error)
}

type authRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepository(db *database.DatabaseClient) *authRepositoryImpl {
	return &authRepositoryImpl{db: db.DB}
}

var (
	STORE_USER_INFO    string = `INSERT INTO user_info(user_id, first_name, last_name, phone, address) VALUES (?, ?, ?, ?, ?)`
	STORE_USER_CRED    string = `INSERT INTO user (role_id, email, password) VALUES (?, ?, ?);`
	FIND_CRED_BY_EMAIL string = `SELECT id, email, password, role_id FROM user WHERE email=?;`
	FIND_CRED_BY_ID    string = `SELECT id, email, password, role_id FROM user WHERE id=?;`
	FIND_USER_BY_ID    string = `SELECT first_name, last_name, phone, address FROM user_info WHERE id=?;`
)

func (auth authRepositoryImpl) StoreUserInfo(ctx context.Context, data models.UserInfo) (err error) {
	_, err = auth.db.ExecContext(ctx, STORE_USER_INFO, data.UserID, data.FirstName, data.LastName, data.Phone, data.Address)
	if err != nil {
		log.Printf("[StoreUserInfo] err: %v\n", err)
		return
	}

	return
}

func (auth authRepositoryImpl) StoreUserCred(ctx context.Context, data models.UserCred) (userID int64, err error) {
	stmt, err := auth.db.PrepareContext(ctx, STORE_USER_CRED)
	if err != nil {
		log.Printf("[StoreUserCred] err: %v\n", err)
		return
	}

	res, err := stmt.ExecContext(ctx, data.UserRole, data.Email, data.Password)
	if err != nil {
		log.Printf("[StoreUserCred] err: %v\n", err)
		return
	}

	userID, err = res.LastInsertId()
	if err != nil {
		log.Printf("[StoreUserCred] err: %v\n", err)
		return 0, err
	}

	return
}

func (auth authRepositoryImpl) GetCredByEmail(ctx context.Context, email string) (usr *models.UserCred, err error) {
	res := auth.db.QueryRowContext(ctx, FIND_CRED_BY_EMAIL, email)
	usr, err = mapCredToEntity(res)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[GetCredByEmail] err: %v\n", err)
		return
	} else if err == sql.ErrNoRows {
		log.Printf("[GetCredByEmail] no data existed in user table\n")
		return nil, errors.ErrInvalidResources
	}

	return
}

func (auth authRepositoryImpl) GetUserInfoByID(ctx context.Context, userID int64) (usr *models.UserInfo, err error) {
	res := auth.db.QueryRowContext(ctx, FIND_USER_BY_ID, userID)
	usr, err = mapUserInfoToEntity(res)
	if err != nil {
		log.Printf("[GetUserInfoByID] err: %v\n", err)
		return
	}

	return
}

func (auth authRepositoryImpl) GetCredByID(ctx context.Context, userID int64) (usr *models.UserCred, err error) {
	res := auth.db.QueryRowContext(ctx, FIND_CRED_BY_ID, userID)
	usr, err = mapCredToEntity(res)
	if err != nil {
		log.Printf("[GetCredByID] err: %v\n", err)
		return
	}

	return
}

func mapCredToEntity(row *sql.Row) (usr *models.UserCred, err error) {
	usr = &models.UserCred{}
	err = row.Scan(&usr.UserID, &usr.Email, &usr.Password, &usr.UserRole)
	return
}

func mapUserInfoToEntity(row *sql.Row) (usr *models.UserInfo, err error) {
	usr = &models.UserInfo{}
	err = row.Scan(&usr.FirstName, &usr.LastName, &usr.Phone, &usr.Address)
	return
}
