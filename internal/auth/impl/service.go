package impl

type AuthServiceImpl struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{repo: repo}
}
