package service

import (
	"fmt"

	"github.com/cyberhawk12121/p2pchat/internal/models"
	"github.com/cyberhawk12121/p2pchat/internal/repository"
	"github.com/cyberhawk12121/p2pchat/internal/transport"
)

type ChatService interface {
	SendMessage(toPeerID string, content string) error
	ReceiveMessage(msg models.Message)
}

type chatService struct {
	peerRepo     repository.PeerRepository
	p2pTransport transport.P2PTransport
}

func NewChatService(peerRepo repository.PeerRepository, p2pTransport transport.P2PTransport) ChatService {
	return &chatService{
		peerRepo:     peerRepo,
		p2pTransport: p2pTransport,
	}
}

func (c *chatService) SendMessage(toPeerID, content string) error {
	// Create msg object and send it using transport service
	msg := models.Message{
		From:    c.p2pTransport.SelfID(),
		Content: content,
		To:      toPeerID,
	}

	return c.p2pTransport.Send(msg)
}

func (c *chatService) ReceiveMessage(msg models.Message) {
	// Here you could handle displaying the message to the user
	// For now we’ll just print it
	fmt.Printf("Received message from %s: %s\n", msg.From, msg.Content)
}
