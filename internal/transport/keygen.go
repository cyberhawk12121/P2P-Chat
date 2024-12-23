package transport

import (
	"crypto/rand"
	"os"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
)

func loadOrCreatePrivateKey() (crypto.PrivKey, error) {
	// (1) Check if a key file exists
	if _, err := os.Stat("mykey"); err == nil {
		// If it exists, load it
		keyBytes, err := os.ReadFile("mykey")
		if err != nil {
			return nil, err
		}
		// Unmarshal the private key
		privKey, err := crypto.UnmarshalPrivateKey(keyBytes)
		if err != nil {
			return nil, err
		}
		return privKey, nil
	}

	// If no key file then generate one and save it in "mykey"
	privKey, err := generatePrivKey()
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func generatePrivKey() (crypto.PrivKey, error) {
	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Marshal and save it
	keyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile("mykey", keyBytes, 0600); err != nil {
		return nil, err
	}

	return privKey, nil
}

func createHostWithPersistentKey() (host.Host, error) {
	privKey, err := loadOrCreatePrivateKey()
	if err != nil {
		return nil, err
	}

	return libp2p.New(
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.Identity(privKey),
	)
}
