package service

type DefaultService struct {
}

func NewDefaultService() *DefaultService {
	return &DefaultService{}
}

func (s *DefaultService) Ping() string {
	return "pong"
}
