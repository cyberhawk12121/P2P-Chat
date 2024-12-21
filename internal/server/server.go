package server

import (
	"context"

	"github.com/cyberhawk12121/p2pchat/internal/db"
	"github.com/cyberhawk12121/p2pchat/internal/repository"
	"github.com/cyberhawk12121/p2pchat/internal/service"
	"github.com/cyberhawk12121/p2pchat/internal/transport"
	"github.com/cyberhawk12121/p2pchat/pkg/logger"
)

type Server struct {
	ctx              context.Context
	cancel           context.CancelFunc
	log              logger.Logger
	transport        transport.P2PTransport
	peerRepo         repository.PeerRepository
	chatService      service.ChatService
	discoveryService service.DiscoveryService
}

func NewServer(ctx context.Context, log logger.Logger) (*Server, error) {
	ctx, cancel := context.WithCancel(ctx)

	// Setup DB and repository
	memdb := db.NewInMemoryDB()
	peerRepo := repository.NewPeerRepository(memdb)

	// Setup P2P transport (pass bootstrap addresses)
	bootstrapPeers := []string{
		// Add known bootstrap nodes multiaddrs here or leave empty if you have none
	}
	p2pTransport, err := transport.NewP2PTransport(bootstrapPeers)
	if err != nil {
		cancel()
		return nil, err
	}

	// Setup services
	chatService := service.NewChatService(peerRepo, p2pTransport)
	discoveryService := service.NewDiscoveryService(peerRepo, p2pTransport)

	// Set callback for received messages
	p2pTransport.SetOnMessageCallback(chatService.ReceiveMessage)

	return &Server{
		ctx:              ctx,
		cancel:           cancel,
		log:              log,
		transport:        p2pTransport,
		peerRepo:         peerRepo,
		chatService:      chatService,
		discoveryService: discoveryService,
	}, nil
}

func (s *Server) Start() error {
	s.log.Info("Starting server...")

	if err := s.transport.Start(s.ctx); err != nil {
		return err
	}

	s.log.Info("Server started with PeerID: ", s.transport.SelfID())
	return nil
}

func (s *Server) Stop() {
	s.log.Info("Stopping server...")
	s.transport.Stop()
	s.cancel()
	s.log.Info("Server stopped.")
}

func (s *Server) SendMessage(peerID, msg string) error {
	return s.chatService.SendMessage(peerID, msg)
}

func (s *Server) ListPeers() [](*servicePeerInfo) {
	peers := s.discoveryService.ListKnownPeers()
	ret := make([]*servicePeerInfo, 0, len(peers))
	for _, p := range peers {
		ret = append(ret, &servicePeerInfo{ID: p.ID, Addresses: p.Addresses})
	}
	return ret
}

type servicePeerInfo struct {
	ID        string
	Addresses []string
}
