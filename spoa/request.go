package main

import (
	"context"
	"github.com/negasus/haproxy-spoe-go/request"
)

type SpopRequest struct {
	*request.Request
	ctxCancel context.CancelFunc
}

func (sr *SpopRequest) Done() {
	sr.ctxCancel()
}

func NewRequestWithCancel(req *request.Request) (*SpopRequest, context.Context) {
	var ctx context.Context
	var newReq *SpopRequest = new(SpopRequest)
	newReq.Request = req
	ctx, newReq.ctxCancel = context.WithCancel(context.Background())
	return newReq, ctx
}
