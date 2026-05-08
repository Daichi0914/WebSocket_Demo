package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"backend/domain"
	"backend/usecase"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	usecase   usecase.MessageUsecase
	clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
	upgrader  websocket.Upgrader
}

func NewWebSocketHandler(u usecase.MessageUsecase) *WebSocketHandler {
	handler := &WebSocketHandler{
		usecase: u,
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo
			},
		},
	}

	// Start listening to messages from Usecase (PubSub)
	go handler.listenToMessages()

	return handler
}

func (h *WebSocketHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws", h.handleWebSocket)
}

func (h *WebSocketHandler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	h.clientsMu.Lock()
	h.clients[conn] = true
	h.clientsMu.Unlock()

	defer func() {
		h.clientsMu.Lock()
		delete(h.clients, conn)
		h.clientsMu.Unlock()
		conn.Close()
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var input struct {
			Sender  string `json:"sender"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(p, &input); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		// Use the usecase to save and broadcast
		if err := h.usecase.SaveAndPublishMessage(input.Sender, input.Content); err != nil {
			log.Println("Failed to save and publish message:", err)
		}
	}
}

func (h *WebSocketHandler) listenToMessages() {
	err := h.usecase.ListenToMessages(func(msg domain.Message) {
		payload, _ := json.Marshal(msg)
		h.broadcast(payload)
	})
	if err != nil {
		log.Println("Error listening to messages:", err)
	}
}

func (h *WebSocketHandler) broadcast(payload []byte) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()

	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Println("Write error:", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}
