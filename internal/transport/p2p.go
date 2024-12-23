package transport

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/cyberhawk12121/p2pchat/internal/models"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

type OnMessageFunc func(models.Message)

type P2PTransport interface {
	Start(context.Context) error
	Stop() error
	SelfID() string
	Send(msg models.Message) error
	ConnectToPeer(ctx context.Context, addr string) error
	SetOnMessageCallback(f OnMessageFunc)
}

type p2pTransport struct {
	host host.Host
	dht  *dht.IpfsDHT
	// OnMessage function is basically how we want to process the received message
	onMessage OnMessageFunc
	// bootstrapPeers are nodes that we need to connect to get the initial data of the network,
	//  like peers to connect etc.
	bootstrapPeers []multiaddr.Multiaddr
}

func NewP2PTransport(bootstrapPeers []string) (P2PTransport, error) {
	// Convert string addrs to multiaddr
	mAddrs := make([]multiaddr.Multiaddr, 0, len(bootstrapPeers))
	for _, addr := range bootstrapPeers {
		ma, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}
		mAddrs = append(mAddrs, ma)
	}

	return &p2pTransport{
		bootstrapPeers: mAddrs,
	}, nil
}

func (p *p2pTransport) Start(ctx context.Context) error {
	// Start a host
	h, err := createHostWithPersistentKey()
	if err != nil {
		return err
	}

	// pubKey := h.Peerstore().PubKey(h.ID())
	// privKey := h.Peerstore().PrivKey(h.ID())

	p.host = h

	// Setup DHT
	d, err := dht.New(ctx, h)
	if err != nil {
		return err
	}

	p.dht = d

	// Connect to bootstrap peers
	if err := p.dht.Bootstrap(ctx); err != nil {
		return err
	}

	for _, pa := range p.bootstrapPeers {
		pi, err := peer.AddrInfoFromP2pAddr(pa)
		if err == nil {
			h.Connect(ctx, *pi)
		}
	}

	// Set stream handler
	p.host.SetStreamHandler(protocol.ID(ChatProtocolID), p.handleStream)

	// Wait a bit for DHT to warm up
	time.Sleep(2 * time.Second)

	return nil
}

func (p *p2pTransport) Stop() error {
	if p.dht != nil {
		p.dht.Close()
	}
	if p.host != nil {
		return p.host.Close()
	}

	return nil
}

func (p *p2pTransport) SelfID() string {
	return p.host.ID().String()
}

func (p *p2pTransport) SetOnMessageCallback(f OnMessageFunc) {
	p.onMessage = f
}

func (p *p2pTransport) Send(msg models.Message) error {
	/**
	1. Start a timer with 5 seconds - take maximum 5 seconds to connect and send
	2. Get the To address
	3. If there is no connection, then find the peer (using DHT) and make connection [MISSING]
	4. Create a stream to the peer
	5. Encode the message
	6. Write the message to the stream
	7. Return error if any
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	peerID, err := peer.Decode(msg.To)
	if err != nil {
		return err
	}

	// Ensure we have a connection
	if p.host.Network().Connectedness(peerID) != network.Connected {
		p.dht.FindPeer(ctx, peerID)
	}

	s, err := p.host.NewStream(ctx, peerID, protocol.ID(ChatProtocolID))
	if err != nil {
		return err
	}

	defer s.Close()
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}

	_, err = s.Write(data)
	return err
}

func (p *p2pTransport) handleStream(s network.Stream) {
	// Step 1: From the stream of data, pull data, ReadAll()
	// Step 2: Decode msg from the data - Unmarshal
	// Step 3: Call the onMessage callback
	defer s.Close()
	data, err := io.ReadAll(s)
	if err != nil {
		fmt.Println("Error reading stream: ", err)
		return
	}

	msg, err := DecodeMessage(data)
	if err != nil {
		fmt.Println("Error decoding data: ", err)
		return
	}

	if p.onMessage != nil {
		p.onMessage(msg)
	}
}

func (p *p2pTransport) ConnectToPeer(ctx context.Context, addr string) error {
	ma, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}
	pi, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		return err
	}
	return p.host.Connect(ctx, *pi)
}
