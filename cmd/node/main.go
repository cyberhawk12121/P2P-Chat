package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cyberhawk12121/p2pchat/internal/server"
	"github.com/cyberhawk12121/p2pchat/pkg/logger"
)

// Step 1: Start a node
// Step 2: Maintain a list of nodes
// Step 3: Send message to particular nodes and rest should not be able to see it
// Step 4:
func main() {
	// Start a logger
	log := logger.NewLogger("debug")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := server.NewServer(ctx, log)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	// Start the server: sets up host, DHT, handle streams, etc
	go srv.Start()

	// Simple CLI to send messages
	reader := bufio.NewReader(os.Stdin)
	log.Info("Type `/peers` to list known peers, `/send <peerID> <msg>` to send a message. `/exit` to quit.")

	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "/peers" {
			peers := srv.ListPeers()
			for _, p := range peers {
				fmt.Println("Known peers: ", p.ID)
			}
			continue
		}

		if strings.HasPrefix(line, "/send") {
			parts := strings.SplitN(line, " ", 3)
			if len(parts) < 3 {
				fmt.Println("Usage: /send <peerID> <message>")
				continue
			}
			peerID := parts[1]
			msg := parts[2]
			err := srv.SendMessage(peerID, msg)

			if err != nil {
				fmt.Println("Error sending message: ", err)
			} else {
				fmt.Println("Message sent.")
			}
			continue
		}
	}

}
