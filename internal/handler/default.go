package handler

import (
	"net/http"

	"encoding/json"

	"github.com/vlladoff/micro-learn/internal/service"
	"go.uber.org/fx"
)

type DefaultHandler struct {
	defaultService *service.DefaultService
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{
		defaultService: service.NewDefaultService(),
	}
}

func (h *DefaultHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	message := h.defaultService.Ping()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

var HandlerModule = fx.Module("handlers",
	fx.Provide(NewDefaultHandler),
)
