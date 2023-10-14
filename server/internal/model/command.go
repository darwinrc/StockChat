package model

type CommmandService interface {
	ProcessCommand(command string, posts []*Post, broadcast chan []byte)
}
