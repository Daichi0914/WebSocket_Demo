package delivery

import (
	"encoding/json"
	"net/http"
	"backend/usecase"
)

type HTTPHandler struct {
	usecase usecase.MessageUsecase
}

func NewHTTPHandler(u usecase.MessageUsecase) *HTTPHandler {
	return &HTTPHandler{usecase: u}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/messages", h.handleGetMessages)
}

func (h *HTTPHandler) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for demo
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	messages, err := h.usecase.GetRecentMessages(50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
