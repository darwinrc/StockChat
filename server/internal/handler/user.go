package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"server/internal/model"
)

type UserHandler struct {
	Service model.UserService
}

// NewUserHandler builds a handler and injects its dependencies
func NewUserHandler(s model.UserService) *UserHandler {
	return &UserHandler{
		Service: s,
	}
}

// Attach attaches the user endpoints to the router
func (h *UserHandler) Attach(r *mux.Router) {
	r.HandleFunc("/login", h.HandleLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", h.HandleSignup).Methods("POST", "OPTIONS")
}

// HandleSignup signs a user up
func (h *UserHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := &model.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		log.Fatalf("error getting user from json: %s", err)
		return
	}

	user, err = h.Service.CreateUser(r.Context(), user)
	if err != nil {
		log.Fatalf("error creating user: %s", err)
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("error getting user to json: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

// HandleLogin logs a user in
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := &model.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		log.Fatalf("error getting user from json: %s", err)
		return
	}

	user, err = h.Service.LoginUser(r.Context(), user)
	if err != nil {
		log.Fatalf("error with login user: %s", err)
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("error getting user to json: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}
