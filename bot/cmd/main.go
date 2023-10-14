package main

import (
	"bot/internal/service"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	stockService := service.NewStockService()
	stockService.ProcessMessages()
}
