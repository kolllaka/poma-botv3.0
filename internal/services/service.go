package services

import (
	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/music"
	"github.com/kolllaka/poma-botv3.0/internal/storage"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
)

type Service interface {
	GetMyPlaylist(isReward bool) Responce

	GetYoutubePlaylistBy(textWithLink string, isReward bool) Responce
	GetYoutubeMusicBy(textWithLink string, isReward bool) Responce

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
func (s *service) GetMyPlaylist(isReward bool) Responce {
	playlist := model.Playlist{}
	playlist, err := s.fmusic.GetPlaylistBy("")
	if err != nil {
		s.logger.Error("error from GetPlaylistBy", logging.ErrAttr(err))

		return Responce{}
	}

	var datas []*model.Music

	for _, song := range playlist.Musics {
		if err := s.GetDuration(song); err != nil {
			s.logger.Warn("error GetDuration", logging.AnyAttr("song", &song), logging.ErrAttr(err))

			continue
		}

		datas = append(datas, song)
	}

	s.logger.Info("загруженно треков", logging.IntAttr("count", len(datas)))

	data := Responce{
		Source:   FILE,
		IsReward: isReward,
		Data:     datas,
	}

	s.logger.Debug("playlist", logging.AnyAttr("data", data))

	return data
}

// GetPlaylist implements Service.
func (s *service) GetYoutubePlaylistBy(textWithLink string, isReward bool) Responce {
	playlist, err := s.ymusic.GetPlaylistBy(textWithLink)
	if err != nil {
		s.logger.Error("error from GetPlaylistBy", logging.ErrAttr(err))

		return Responce{}
	}

	data := Responce{
		Source:   YOUTUBE,
		IsReward: isReward,
		Data:     playlist,
	}

	return data
}

// GetMusic implements Service.
func (s *service) GetYoutubeMusicBy(textWithLink string, isReward bool) Responce {
	music, err := s.ymusic.GetMusicBy(textWithLink)
	if err != nil {
		s.logger.Error("error from GetMusicBy", logging.ErrAttr(err))

		return Responce{}
	}

	data := Responce{
		Source:   YOUTUBE,
		IsReward: isReward,
		Data:     music,
	}

	return data
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
