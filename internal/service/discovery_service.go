package service

import (
	"github.com/cyberhawk12121/p2pchat/internal/models"
	"github.com/cyberhawk12121/p2pchat/internal/repository"
	"github.com/cyberhawk12121/p2pchat/internal/transport"
)

type DiscoveryService interface {
	ListKnownPeers() []*models.PeerInfo
	// Additional methods: FindPeer, etc.
}

type discoveryService struct {
	peerRepo     repository.PeerRepository
	p2pTransport transport.P2PTransport
}

func NewDiscoveryService(peerRepo repository.PeerRepository, p2p transport.P2PTransport) DiscoveryService {
	return &discoveryService{peerRepo: peerRepo, p2pTransport: p2p}
}

func (d *discoveryService) ListKnownPeers() []*models.PeerInfo {
	return d.peerRepo.GetPeers()
}
