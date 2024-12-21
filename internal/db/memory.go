package db

import (
	"sync"

	"github.com/cyberhawk12121/p2pchat/internal/models"
)

type InMemoryDB struct {
	mu    sync.Mutex
	peers map[string]*models.PeerInfo
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		peers: make(map[string]*models.PeerInfo),
	}
}

// Just add the peers
func (db *InMemoryDB) AddPeer(peer *models.PeerInfo) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.peers[peer.ID] = peer
}

func (db *InMemoryDB) GetPeers() []*models.PeerInfo {
	db.mu.Lock()
	defer db.mu.Unlock()
	peers := make([]*models.PeerInfo, 0, len(db.peers))
	for _, p := range db.peers {
		peers = append(peers, p)
	}

	return peers
}
