package eventsub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	wsHost = "eventsub.wss.twitch.tv"

	bHost = "api.twitch.tv"
	bPath = "helix/eventsub/subscriptions"
)

type Secrets struct {
	SecretClientId string
	SecretBearer   string
}

type SubscribeFunc func(sessionId string, condition Condition) error

type client struct {
	host  string
	bPath string

	secretClientId string
	secretBearer   string

	keepaliveTimeoutSeconds int
	transport               Transport
	subscribes              []Subscription

	isConnect bool
	socket    *websocket.Conn
	client    http.Client
	EventChan chan EventRawMessage
}

func New(secrets Secrets) Client {
	return &client{
		host:                    bHost,
		bPath:                   bPath,
		secretClientId:          secrets.SecretClientId,
		secretBearer:            secrets.SecretBearer,
		keepaliveTimeoutSeconds: 30,
		transport: Transport{
			Method: "websocket",
		},
		isConnect: true,
		client:    http.Client{},
		EventChan: make(chan EventRawMessage),
	}
}

func (c *client) reader() {
	for {
		msgType, bytes, err := c.socket.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			break
		}

		var msg Message
		if err := json.Unmarshal(bytes, &msg); err != nil {
			continue
		}

		log.Printf("msg from socket %+v\n", msg)

		switch msg.Metadata.MessageType {
		case "session_welcome":
			if c.isConnect {
				sessId := msg.Payload.Session.ID
				for _, sub := range c.subscribes {
					c.subscribe(sessId, sub)
				}

				c.isConnect = true
			}

		case "session_keepalive":
			continue
		case "session_reconnect":
			c.socket, _ = c.connect(msg.Payload.Session.ReconnectURL)
			c.isConnect = false
		case "notification":
			bytes, _ := json.Marshal(msg.Payload.Event)

			log.Printf("event: %+v", string(bytes))
			c.EventChan <- EventRawMessage{
				EventType: msg.Metadata.SubscriptionType,
				Event:     msg.Payload.Event,
			}
		}
	}
}

func (c *client) doRequest(data Subscription, query url.Values) error {
	const op = "eventsub.doRequest"

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   c.bPath,
	}
	u.RawQuery = query.Encode()

	network, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%s: error to marshal: %w", op, err)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(network))
	req.Header.Add("Authorization", c.secretBearer)
	req.Header.Add("Client-Id", c.secretClientId)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	if err != nil {
		return fmt.Errorf("%s: error to NewRequest: %w", op, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: error to Do request: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: error to ReadAll: %w", op, err)
	}

	log.Printf("resp.Body %+v\n", string(body))

	return nil
}

func (c *client) subscribe(sessionId string, sData Subscription) error {
	sData.Transport.SessionID = sessionId
	if err := c.doRequest(sData, nil); err != nil {
		return err
	}

	return nil
}

func (c *client) defaultWSUrl() string {
	u := url.URL{Scheme: "wss", Host: wsHost, Path: "/ws"}
	q := u.Query()
	q.Set("keepalive_timeout_seconds", fmt.Sprint(c.keepaliveTimeoutSeconds))
	u.RawQuery = q.Encode()

	return u.String()
}

func (c *client) connect(url string) (*websocket.Conn, error) {
	if url == "" {
		url = c.defaultWSUrl()
	}

	socket, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return socket, nil
}
