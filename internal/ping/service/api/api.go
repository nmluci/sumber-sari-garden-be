package api

import "context"

type PingService interface {
	Ping(ctx context.Context) string
}