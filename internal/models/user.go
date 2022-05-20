package models

type UserInfo struct {
	FirstName string
	LastName  string
	Phone     string
	Address   string
	RoleName  string
	UserID    int64
	RoleID    int64
}

type UserCred struct {
	Email    string
	Password string
	UserID   int64
	UserRole int64
}

type UserContext struct {
	UserID int64
	Priv   int64
}

type AuthCtxKey string
