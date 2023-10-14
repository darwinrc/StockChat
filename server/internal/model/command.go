package model

type CommmandService interface {
	ProcessCommand(command string, broadcast chan []byte)
}
