package eventsub

import (
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	TYPE_CHANNEL_RAID           = "channel.raid"
	TYPE_CHANNEL_FOLLOW         = "channel.follow"
	TYPE_CHANNEL_SUBSCRIBE      = "channel.subscribe"
	TYPE_CHANNEL_SUBSCRIBE_GIFT = "channel.subscription.gift"
	TYPE_CHANNEL_RE_SUBSCRIBE   = "channel.subscription.message"
	TYPE_CHANNEL_CHEER          = "channel.cheer"

	SUBSCRIPTIONS_PATH = "subscriptions"
)

type Client interface {
	Connect() error
	Close() error

	// keepaliveTimeoutSeconds in range 10-300.
	SetKeepaliveTimeoutSeconds(keepaliveTimeoutSeconds int)
	GetEventChan() chan EventRawMessage

	// Need set ToBroadcasterUserID or FromBroadcasterUserID. ONLY one!
	// No authorization required.
	SubscribeChannelRaid(condition Condition)
	// need set BroadcasterUserId and ModeratorUserId
	// Must have moderator:read:followers scope.
	SubscribeChannelFollow(condition Condition)
	// need set BroadcasterUserId
	// Must have channel:read:subscriptions scope.
	SubscribeChannelSubscribe(condition Condition)
	// need set BroadcasterUserId
	// Must have channel:read:subscriptions scope.
	SubscribeChannelSubscribeGift(condition Condition)
	// need set BroadcasterUserId
	// Must have channel:read:subscriptions scope.
	SubscribeChannelReSubscribe(condition Condition)
	// need set BroadcasterUserId
	// Must have bits:read scope.
	SubscribeChannelCheer(condition Condition)
}

// keepaliveTimeoutSeconds 10 - 600
func (c *client) Connect() error {
	socket, err := c.connect("")
	if err != nil {
		return fmt.Errorf("error to connect to ws: %w", err)
	}

	c.socket = socket

	go c.reader()

	return nil
}

func (c *client) Close() error {
	c.socket.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)

	return c.socket.Close()
}

func (c *client) SetKeepaliveTimeoutSeconds(keepaliveTimeoutSeconds int) {
	c.keepaliveTimeoutSeconds = keepaliveTimeoutSeconds
}

func (c *client) GetEventChan() chan EventRawMessage {
	return c.EventChan
}

// Need set ToBroadcasterUserID or FromBroadcasterUserID. ONLY one!
// No authorization required.
func (c *client) SubscribeChannelRaid(condition Condition) {
	sData := Subscription{
		Type:      TYPE_CHANNEL_RAID,
		Version:   "1",
		Condition: condition,
		Transport: c.transport,
	}

	c.subscribes = append(c.subscribes, sData)
}

// need set BroadcasterUserId and ModeratorUserId
// Must have moderator:read:followers scope.
func (c *client) SubscribeChannelFollow(condition Condition) {
	sData := Subscription{
		Type:      TYPE_CHANNEL_FOLLOW,
		Version:   "2",
		Condition: condition,
		Transport: c.transport,
	}

	c.subscribes = append(c.subscribes, sData)
}

// need set BroadcasterUserId
// Must have channel:read:subscriptions scope.
func (c *client) SubscribeChannelSubscribe(condition Condition) {
	sData := Subscription{
		Type:      TYPE_CHANNEL_SUBSCRIBE,
		Version:   "1",
		Condition: condition,
		Transport: c.transport,
	}

	c.subscribes = append(c.subscribes, sData)
}

// need set BroadcasterUserId
// Must have channel:read:subscriptions scope.
func (c *client) SubscribeChannelSubscribeGift(condition Condition) {
	sData := Subscription{
		Type:      TYPE_CHANNEL_SUBSCRIBE_GIFT,
		Version:   "1",
		Condition: condition,
		Transport: c.transport,
	}

	c.subscribes = append(c.subscribes, sData)
}

// need set BroadcasterUserId
// Must have channel:read:subscriptions scope.
func (c *client) SubscribeChannelReSubscribe(condition Condition) {
	sData := Subscription{
		Type:      TYPE_CHANNEL_RE_SUBSCRIBE,
		Version:   "1",
		Condition: condition,
		Transport: c.transport,
	}

	c.subscribes = append(c.subscribes, sData)
}

// need set BroadcasterUserId
// Must have bits:read scope.
func (c *client) SubscribeChannelCheer(condition Condition) {
	sData := Subscription{
		Type:      TYPE_CHANNEL_CHEER,
		Version:   "1",
		Condition: condition,
		Transport: c.transport,
	}

	c.subscribes = append(c.subscribes, sData)
}
