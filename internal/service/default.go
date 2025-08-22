package service

import (
	"context"
	"log"

	"github.com/vlladoff/micro-learn/internal/middleware"
	"go.uber.org/fx"
)

type DefaultService struct {
}

func NewDefaultService() *DefaultService {
	return &DefaultService{}
}

func (s *DefaultService) Ping(ctx context.Context) string {
	requestID := middleware.GetRequestID(ctx)

	log.Printf("[%s] Processing ping request in service", requestID)

	return "pong"
}

var ServiceModule = fx.Module("services",
	fx.Provide(NewDefaultService),
)
