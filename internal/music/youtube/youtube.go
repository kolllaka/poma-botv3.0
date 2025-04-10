package youtube

import (
	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
	"github.com/kolllaka/poma-botv3.0/pkg/youtubeapi"
)

type youtube struct {
	logger     *logging.Logger
	youtubeapi youtubeapi.Api
}

func New(
	logger *logging.Logger,
	youtubeapi youtubeapi.Api,
) *youtube {
	return &youtube{
		logger:     logger,
		youtubeapi: youtubeapi,
	}
}

// GetMusicBy implements MusicService.
func (y *youtube) GetMusicBy(text string) (model.Music, error) {
	songId, err := parseSongIdByText(text)
	if err != nil {
		y.logger.Warn("bad user input", logging.StringAttr("text", text), logging.ErrAttr(err))

		return model.Music{}, model.ErrBadLink
	}

	song, err := y.youtubeapi.ReqSongInfoById(songId)
	if err != nil {
		y.logger.Warn("error ReqSongInfoById", logging.StringAttr("songId", songId), logging.ErrAttr(err))

		return model.Music{}, err
	}

	if len(song.Items) < 1 {
		y.logger.Warn("song not found", logging.StringAttr("songId", songId))

		return model.Music{}, nil
	}

	var music model.Music

	music.Name = song.Items[0].Snippet.Title
	music.Link = songId
	music.Duration = formatTimeFromYoutube(song.Items[0].ContentDetails.Duration)

	return music, nil
}

// GetPlaylistBy implements MusicService.
func (y *youtube) GetPlaylistBy(link string) (model.Playlist, error) {
	panic("unimplemented")
}
