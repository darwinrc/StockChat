package service

import (
	"bufio"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	mock_infra "server/internal/infra/mocks"
	mock_repo "server/internal/repo/mocks"
	"testing"
)

func TestProcessCommandSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	payload := []byte("{\"stockCode\":\"aapl.us\"}")

	mockAMQP := mock_infra.NewMockAMQPClient(ctrl)
	mockAMQP.EXPECT().SetupAMQExchange().Return(nil)
	mockAMQP.EXPECT().PublishAMQMessage(payload).Return(nil)

	mockPostRepo := &mock_repo.MockPostRepo{}

	service := NewCommandService(mockPostRepo, mockAMQP)

	broadcast := make(chan []byte)
	scanner, reader, writer := mockLogger(t)
	defer resetLogger(reader, writer)

	go service.ProcessCommand("/stock=aapl.us", broadcast)

	scanner.Scan() // first log: Processing command ...
	scanner.Scan() // last log: Stock sent ...
	got := scanner.Text()
	msg := "Stock sent: {\"stockCode\":\"aapl.us\"}"
	assert.Contains(t, got, msg)
}

func TestProcessCommandFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := &mock_repo.MockPostRepo{}
	mockAMQP := mock_infra.NewMockAMQPClient(ctrl)

	service := NewCommandService(mockPostRepo, mockAMQP)

	broadcast := make(chan []byte)
	scanner, reader, writer := mockLogger(t)
	defer resetLogger(reader, writer)

	go service.ProcessCommand("invalid-command", broadcast)

	scanner.Scan() // first log: Processing command ...
	scanner.Scan() // last log: invalid command ...
	got := scanner.Text()
	msg := "invalid command: invalid-command. It should be something like /stock=aapl.us"
	assert.Contains(t, got, msg)
}

func TestBroadcastCommand(t *testing.T) {
	
}

// Util functions
func mockLogger(t *testing.T) (*bufio.Scanner, *os.File, *os.File) {
	reader, writer, err := os.Pipe()
	if err != nil {
		assert.Fail(t, "couldn't get os Pipe: %v", err)
	}
	log.SetOutput(writer)

	return bufio.NewScanner(reader), reader, writer
}

func resetLogger(reader *os.File, writer *os.File) {
	err := reader.Close()
	if err != nil {
		fmt.Println("error closing reader was ", err)
	}
	if err = writer.Close(); err != nil {
		fmt.Println("error closing writer was ", err)
	}
	log.SetOutput(os.Stderr)
}
