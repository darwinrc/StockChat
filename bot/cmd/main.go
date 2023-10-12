package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "This is the bot"}`)
}

func main() {
	server := http.NewServeMux()

	server.HandleFunc("/", handler)

	fmt.Println("Server listening on :5000")
	http.ListenAndServe(":5000", server)
}
