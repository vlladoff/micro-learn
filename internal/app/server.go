package app

import (
	"net/http"

	"github.com/vlladoff/micro-learn/internal/handler"
	"go.uber.org/fx"
)

type SmplServer struct {
	defaultHandler *handler.DefaultHandler
	jobHandler     *handler.JobHandler
}

func NewSmplServer(defaultHandler *handler.DefaultHandler, jobHandler *handler.JobHandler) *SmplServer {
	return &SmplServer{
		defaultHandler: defaultHandler,
		jobHandler:     jobHandler,
	}
}

func (s *SmplServer) GetPing(w http.ResponseWriter, r *http.Request) {
	s.defaultHandler.GetPing(w, r)
}

func (s *SmplServer) CreateJob(w http.ResponseWriter, r *http.Request) {
	s.jobHandler.CreateJob(w, r)
}

func (s *SmplServer) GetJob(w http.ResponseWriter, r *http.Request, jobID string) {
	s.jobHandler.GetJob(w, r)
}

func (s *SmplServer) DeleteJob(w http.ResponseWriter, r *http.Request, jobID string) {
	s.jobHandler.DeleteJob(w, r)
}

var AppModule = fx.Module("app",
	fx.Provide(NewSmplServer),
)
