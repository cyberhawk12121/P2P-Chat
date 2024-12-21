package repository

import (
	"github.com/cyberhawk12121/p2pchat/internal/db"
	"github.com/cyberhawk12121/p2pchat/internal/models"
)

type PeerRepository interface {
	AddPeer(peer *models.PeerInfo)
	GetPeers() []*models.PeerInfo
}

type peerRepository struct {
	db *db.InMemoryDB
}

func NewPeerRepository(db *db.InMemoryDB) PeerRepository {
	return &peerRepository{db: db}
}

func (r *peerRepository) AddPeer(peer *models.PeerInfo) {
	r.db.AddPeer(peer)
}

func (r *peerRepository) GetPeers() []*models.PeerInfo {
	return r.db.GetPeers()
}
