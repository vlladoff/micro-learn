package app

import (
	"net/http"

	"github.com/vlladoff/micro-learn/internal/handler"
)

type SmplServer struct {
	defaultHandler *handler.DefaultHandler
}

func NewSmplServer() *SmplServer {
	defaultHandler := handler.NewDefaultHandler()

	return &SmplServer{
		defaultHandler: defaultHandler,
	}
}

func (s *SmplServer) GetPing(w http.ResponseWriter, r *http.Request) {
	s.defaultHandler.GetPing(w, r)
}
