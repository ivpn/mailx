package model

type MessageType string

const (
	Forward MessageType = "Forward"
	Block   MessageType = "Block"
	Reply   MessageType = "Reply"
	Send    MessageType = "Send"
)

type Message struct {
	BaseModel
	AliasID string      `json:"alias_id"`
	Type    MessageType `json:"type"`
	Size    int         `json:"size"`
}
