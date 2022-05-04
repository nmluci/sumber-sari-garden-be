package ping

import (
	"context"

	"github.com/nmluci/sumber-sari-garden/internal/ping/impl"
)

type PingService interface {
	Ping(ctx context.Context) string
}

func NewPingService() PingService {
	return impl.NewPingService()
}
