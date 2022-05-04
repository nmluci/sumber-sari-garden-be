package impl

import "context"

type pingServiceimpl struct{
	
}

func (ps pingServiceimpl) Ping(ctx context.Context) string {
	return "meow"
}

func NewPingService() *pingServiceimpl {
	return &pingServiceimpl{}
}