package entity

type UserInfo struct {
	FirstName string
	LastName  string
	Phone     string
	Address   string
	RoleName  string
	RoleID    *int64
}

type UserCred struct {
	Email    string
	Password string
	UserID   *int64
}
