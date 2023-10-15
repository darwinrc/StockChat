package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"server/db"
	"server/internal/handler"
	"server/internal/infra"
	"server/internal/repo"
	"server/internal/service"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := db.NewDatabase()
	if err != nil {
		fmt.Errorf("error getting db connection: %s", err.Error())
	}

	router := mux.NewRouter()

	userRepo := repo.NewUserRepository(conn.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userHandler.Attach(router)

	amqpClient := infra.NewAMQPClient()

	postRepo := repo.NewPostRepository(conn.GetDB())
	postService := service.NewPostService(postRepo)
	commandService := service.NewCommandService(postRepo, amqpClient)
	postHandler := handler.NewPostHandler(postService, commandService)
	postHandler.Attach(router)

	// Separate goroutine for listening to new messages
	go postHandler.WriteMessages()

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})

	err = http.ListenAndServe(":5000", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router))
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}

	fmt.Println("Server listening on port 5000")
}
