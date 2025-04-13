package eventsub

import (
	"encoding/json"
	"errors"
)

var (
	// ErrShardTooManyTopics returned when a shard has attempted to join too many topics
	ErrShardTooManyTopics = errors.New("too many topics on shard")
	// ErrShardIDOutOfBounds returned when an invalid shard id is provided
	ErrShardIDOutOfBounds = errors.New("shard id out of bounds")
	// ErrNonceTimeout returned when the server doesnt respond to a nonced message in time
	ErrNonceTimeout = errors.New("nonced message timeout")
	// ErrPingTimeout returned when the server takes too long to respond to a ping message
	ErrPingTimeout = errors.New("server took too long to respond to ping")

	// ErrBadMessage returned when the server receives an invalid message
	ErrBadMessage = errors.New("server received an invalid message")
	// ErrBadAuth returned when a topic doesnt have the permissions required
	ErrBadAuth = errors.New("bad authentication for topic")
	// ErrBadTopic returned when an invalid topic was requested
	ErrBadTopic = errors.New("invalid topic")
	// ErrServer returned when something went wrong on the servers end
	ErrServer = errors.New("something went wrong on the servers end")
	// ErrUnknown returned when the server sends back an error that wasnt handled by the reader
	ErrUnknown = errors.New("server sent back an unknown error")

	// ErrInvalidNonceGenerator returned when a provided nonce generator can not be used
	ErrInvalidNonceGenerator = errors.New("nonce generator is invalid")
)

type EventRawMessage struct {
	EventType string
	Event     json.RawMessage
}

type Message struct {
	Metadata Metadata `json:"metadata,omitempty"`
	Payload  Payload  `json:"payload,omitempty"`
}

type Metadata struct {
	MessageID           string `json:"message_id,omitempty"`
	MessageType         string `json:"message_type,omitempty"`
	MessageTimestamp    string `json:"message_timestamp,omitempty"`
	SubscriptionType    string `json:"subscription_type,omitempty"`
	SubscriptionVersion string `json:"subscription_version,omitempty"`
}

type Payload struct {
	Session      Session         `json:"session,omitempty"`
	Subscription Subscription    `json:"subscription,omitempty"`
	Event        json.RawMessage `json:"event,omitempty"`
}

type Session struct {
	ID                      string `json:"id,omitempty"`
	Status                  string `json:"status,omitempty"`
	ConnectedAt             string `json:"connected_at,omitempty"`
	KeepaliveTimeoutSeconds int    `json:"keepalive_timeout_seconds,omitempty"`
	ReconnectURL            string `json:"reconnect_url,omitempty"`
}

type Subscription struct {
	ID        string    `json:"id,omitempty"`
	Status    string    `json:"status,omitempty"`
	Type      string    `json:"type,omitempty"`
	Version   string    `json:"version,omitempty"`
	Condition Condition `json:"condition,omitempty"`
	Transport Transport `json:"transport,omitempty"`
	CreatedAt string    `json:"created_at,omitempty"`
	Cost      int       `json:"cost,omitempty"`
}

type Condition struct {
	FromBroadcasterUserID string `json:"from_broadcaster_user_id,omitempty"`
	ToBroadcasterUserID   string `json:"to_broadcaster_user_id,omitempty"`
	BroadcasterUserId     string `json:"broadcaster_user_id,omitempty"`
	ModeratorUserId       string `json:"moderator_user_id,omitempty"`
}

type Transport struct {
	Method    string `json:"method,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}
