package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vlladoff/micro-learn/internal/handler"
	"github.com/vlladoff/micro-learn/internal/middleware"
	"go.uber.org/fx"
)

type SmplServer struct {
	jobHandler *handler.JobHandler
}

func NewSmplServer(jobHandler *handler.JobHandler) *SmplServer {
	return &SmplServer{
		jobHandler: jobHandler,
	}
}

func (s *SmplServer) GetPing(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log.Printf("[%s] Ping request", requestID)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
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
