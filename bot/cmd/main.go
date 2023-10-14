package main

import (
	"bot/internal/service"
)

func main() {
	stockService := service.NewStockService()
	stockService.ProcessMessages()
}
