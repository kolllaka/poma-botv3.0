package router

import (
	"encoding/json"
)

type MsgFromSocket struct {
	Reason string          `json:"reason,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

type MusicFromSocket struct {
	Data struct {
		Title    string `json:"title,omitempty"`
		Link     string `json:"link,omitempty"`
		Duration int    `json:"duration,omitempty"`
	} `json:"data"`
}
