package main

import (
	"log"
	"net/http"

	"backend/delivery"
	"backend/infrastructure"
	"backend/usecase"
)

func main() {
	// Initialize Infrastructure (Database & Redis)
	db := infrastructure.NewMySQLConnection()
	rdb := infrastructure.NewRedisClient()

	// Initialize Repositories (Infrastructure Layer)
	messageRepo := infrastructure.NewMysqlMessageRepository(db)
	pubsubRepo := infrastructure.NewRedisPubSub(rdb)

	// Initialize Usecases (Application Layer)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, pubsubRepo)

	// Initialize Handlers (Delivery Layer)
	httpHandler := delivery.NewHTTPHandler(messageUsecase)
	wsHandler := delivery.NewWebSocketHandler(messageUsecase)

	// Setup Router
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)
	wsHandler.RegisterRoutes(mux)

	// Start Server
	log.Println("Server starting on :8080 (Clean Architecture)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
