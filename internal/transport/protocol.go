package transport

import (
	"encoding/json"

	"github.com/cyberhawk12121/p2pchat/internal/models"
)

const ChatProtocolID = "/p2pchat/1.0.0"

// Encode message to bytes
func EncodeMessage(msg models.Message) ([]byte, error) {
	return json.Marshal(msg)
}

// Decode message from bytes
func DecodeMessage(data []byte) (models.Message, error) {
	var msg models.Message
	err := json.Unmarshal(data, &msg)
	return msg, err
}
