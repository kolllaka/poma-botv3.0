package model

type Playlist struct {
	Name   string   `json:"title"`
	Link   string   `json:"link"`
	Author string   `json:"author,omitempty"`
	Musics []*Music `jsoc:"musics"`
}

type Music struct {
	Name     string `json:"title"`
	Link     string `json:"link"`
	Author   string `json:"author,omitempty"`
	Duration int    `json:"duration"`
}
