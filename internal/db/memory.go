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

