package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"server/internal/model"
	"strings"
)

type PostHandler struct {
	Service        model.PostService
	CommandService model.CommmandService
}

const (
	stockBotMessage = "/stock="
)

var (
	broadcast = make(chan []byte)
	clients   = make(map[*websocket.Conn]bool)

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// NewPostHandler builds a handler and injects its dependencies
func NewPostHandler(s model.PostService, cs model.CommmandService) *PostHandler {
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("error upgrading connection to support websockets: %s", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	h.readMessages(r.Context(), conn)
}

// readMessages watches for messages coming through the websocket connection and queues them in the broadcast channel
func (h *PostHandler) readMessages(ctx context.Context, conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("error getting reader: %s", err)
			return
		}

		post := &model.Post{}
		if err := json.Unmarshal(msg, &post); err != nil {
			log.Fatalf("error getting post from json: %s", err)
			return
		}

		// if the messages is a command to query a stock, process the command asynchronously
		// and share the broadcast channel, so it can send the message back to the chatroom
		if strings.Index(post.Message, stockBotMessage) == 0 {
			go h.CommandService.ProcessCommand(post.Message, broadcast)

			continue
		}

		posts, err := h.Service.GetRecentPosts(ctx)
		if err != nil {
			log.Fatalf("error getting recents posts: %s", err)
			return
		}

		newPost, err := h.Service.CreatePost(ctx, post)
		if err != nil {
			log.Fatalf("error creating post: %s", err)
			return
		}

		posts = append([]*model.Post{newPost}, posts...)

		jsonPosts, err := json.Marshal(posts)
		if err != nil {
			log.Fatalf("error getting posts to json: %s", err)
			return
		}

		broadcast <- jsonPosts
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
