package service

import "go.uber.org/fx"

type DefaultService struct {
}

func NewDefaultService() *DefaultService {
	return &DefaultService{}
}

func (s *DefaultService) Ping() string {
	return "pong"
}

var ServiceModule = fx.Module("services",
	fx.Provide(NewDefaultService),
)
