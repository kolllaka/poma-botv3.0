package services

import (
	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/music"
	"github.com/kolllaka/poma-botv3.0/internal/storage"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
)

type Service interface {
	GetMyPlaylist(isReward bool) (Responce, error)

	GetYoutubePlaylistBy(playlist model.Playlist, isReward bool) (Responce, error)
	GetYoutubeMusicBy(music model.Music, isReward bool) (Responce, error)

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
func (s *service) GetMyPlaylist(isReward bool) (Responce, error) {
	playlist := model.Playlist{}
	playlist, err := s.fmusic.GetPlaylistBy("")
	if err != nil {
		s.logger.Error("error from GetPlaylistBy", logging.ErrAttr(err))

		return Responce{}, err
	}

	var datas []*model.Music

	for _, song := range playlist.Musics {
		if err := s.GetDuration(song); err != nil {
			s.logger.Warn("error GetDuration", logging.AnyAttr("song", &song), logging.ErrAttr(err))
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

	return data, nil
}

// GetPlaylist implements Service.
func (s *service) GetYoutubePlaylistBy(playlist model.Playlist, isReward bool) (Responce, error) {
	filledPlaylist, err := s.ymusic.GetPlaylistBy(playlist.Link)
	if err != nil {
		s.logger.Error("error from GetPlaylistBy", logging.ErrAttr(err))

		return Responce{}, err
	}
	filledPlaylist.Author = playlist.Author

	data := Responce{
		Source:   YOUTUBE,
		IsReward: isReward,
		Data:     filledPlaylist,
	}

	return data, nil
}

// GetMusic implements Service.
func (s *service) GetYoutubeMusicBy(music model.Music, isReward bool) (Responce, error) {
	filledMusic, err := s.ymusic.GetMusicBy(music.Link)
	if err != nil {
		return Responce{}, err
	}
	filledMusic.Author = music.Author

	data := Responce{
		Source:   YOUTUBE,
		IsReward: isReward,
		Data:     filledMusic,
	}

	return data, nil
}

// GetDuration implements Service.
func (s *service) GetDuration(music *model.Music) error {
	storeDuration := storage.StoreDuration{
		Link: music.Name,
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
