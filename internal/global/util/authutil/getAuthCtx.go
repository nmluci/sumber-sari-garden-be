package authutil

import (
	"context"
	"log"

	"github.com/nmluci/sumber-sari-garden/internal/models"
)

func GetUserIDFromCtx(ctx context.Context) int64 {
	usr, ok := ctx.Value(models.AuthCtxKey("user-ctx")).(*models.UserContext)
	if !ok {
		log.Printf("[GetUserIDFromCtx] invalid user-id\n")
		return 0
	}

	return usr.UserID
}

func GetUserPrivFromCtx(ctx context.Context) int64 {
	usr, ok := ctx.Value(models.AuthCtxKey("user-ctx")).(*models.UserContext)
	if !ok {
		log.Printf("[GetUserPrivFromCtx] invalid user-level\n")
		return 0
	}

	return usr.Priv
}
