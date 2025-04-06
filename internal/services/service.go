package services

import (
	"bytes"
	"encoding/json"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/music"
	"github.com/kolllaka/poma-botv3.0/internal/storage"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
)

type Service interface {
	GetMyPlaylist() []byte

	GetYoutubePlaylistBy(msg model.RewardMessage) []byte
	GetYoutubeMusicBy(msg model.RewardMessage) []byte

	GetDuration(music *model.Music) error
	StoreDuration(music *model.Music) error
}

type service struct {
	logger *logging.Logger

	fmusic  music.MusicService
	ymusic  music.MusicService
	storage storage.Store
}

func New(
	logger *logging.Logger,
	fmusic music.MusicService,
	ymusic music.MusicService,
	storage storage.Store,
) Service {
	return &service{
		logger:  logger,
		fmusic:  fmusic,
		ymusic:  ymusic,
		storage: storage,
	}
}

// GetMyPlaylist implements Service.
func (s *service) GetMyPlaylist() []byte {
	playlist := model.Playlist{}
	s.fmusic.GetPlaylistBy("", &playlist)

	var datas []Responce

	for _, song := range playlist.Musics {
		if err := s.GetDuration(song); err != nil {
			s.logger.Warn("error GetDuration", logging.AnyAttr("song", &song), logging.ErrAttr(err))

			continue
		}

		data := Responce{
			Source:   FILE,
			IsReward: false,
			Data:     &song,
		}

		datas = append(datas, data)
	}

	s.logger.Info("загруженно треков", logging.IntAttr("count", len(datas)))

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(datas)

	s.logger.Debug("json", logging.StringAttr("datas", network.String()))

	return network.Bytes()
}

// GetPlaylist implements Service.
func (s *service) GetYoutubePlaylistBy(msg model.RewardMessage) []byte {
	playlist := model.Playlist{
		Author: msg.Username,
	}
	s.ymusic.GetPlaylistBy(msg.Text, &playlist)

	data := Responce{
		Source:   YOUTUBE,
		IsReward: msg.IsReward,
		Data:     playlist,
	}

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(data)

	return network.Bytes()
}

// GetMusic implements Service.
func (s *service) GetYoutubeMusicBy(msg model.RewardMessage) []byte {
	music := model.Music{
		Author: msg.Username,
	}
	s.ymusic.GetMusicBy(msg.Text, &music)

	data := Responce{
		Source:   YOUTUBE,
		IsReward: msg.IsReward,
		Data:     music,
	}

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(data)

	return network.Bytes()
}

// GetDuration implements Service.
func (s *service) GetDuration(music *model.Music) error {
	storeDuration := storage.StoreDuration{
		Link: music.Link,
	}

	if err := s.storage.GetDuration(&storeDuration); err != nil {
		return err
	}

	music.Duration = storeDuration.Duration

	return nil
}

// StoreDuration implements Service.
func (s *service) StoreDuration(music *model.Music) error {
	storeDuration := storage.StoreDuration{
		Link:     music.Link,
		Duration: music.Duration,
	}

	return s.storage.StoreDuration(&storeDuration)
}
