package health

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	// TODO add check to connection database, ....
}

func NewHandler() *Handler {
	return &Handler{}
}

// Health check dependencies like database connections, ... if are working fine
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{
		Status: "Ok",
	})
}
