package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"server/internal/model"
	"server/internal/service"
)

type PostHandler struct {
	Service        service.PostService
	CommandService service.CommmandService
}

var (
	broadcast = make(chan []byte)
	clients   = make(map[*websocket.Conn]bool)
)

// NewPostHandler builds a handler and injects its dependencies
func NewPostHandler(s service.PostService, cs service.CommmandService) *PostHandler {
	return &PostHandler{
		Service:        s,
		CommandService: cs,
	}
}

// Attach attaches the web socket endpoints to the router
func (h *PostHandler) Attach(r *mux.Router) {
	r.HandleFunc("/ws", h.HandleWebSocketConnection)
}

// HandleWebSocketConnection establishes a web socket connection and reads messages coming through it
func (h *PostHandler) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection to support websockets: %s", err)
		return
	}

	h.readMessages(r.Context(), conn)
}

// readMessages watches for messages coming through the websocket connection and queues them in the broadcast channel
func (h *PostHandler) readMessages(ctx context.Context, conn *websocket.Conn) {
	defer conn.Close()

	clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error getting reader: %s", err)
		}

		var post *model.Post
		if err := json.Unmarshal(msg, &post); err != nil {
			log.Printf("error getting post from json: %s", err)
		}

		stockCode, err := h.CommandService.ParseCommand(post.Message)
		if err := json.Unmarshal(msg, &post); err != nil {
			log.Printf("error parsing the command: %s", err)
		}

		if stockCode != "" {
			// if the message is a command to query a stock, process the command asynchronously
			// share the broadcast channel, so it can send the message back to the chatroom
			go h.CommandService.ProcessCommand(stockCode)
			go h.CommandService.BroadcastCommand(broadcast)

			continue
		}

		if err := h.Service.CreatePost(ctx, post, broadcast); err != nil {
			log.Printf("error creating post: %s", err)
		}
	}
}

// WriteMessages watches for messages in the broadcast channel and send them to all connected clients
func (h *PostHandler) WriteMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				delete(clients, client)
				client.Close()
			}
		}
	}
}
