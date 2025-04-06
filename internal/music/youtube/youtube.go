package youtube

import (
	"fmt"

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
func (y *youtube) GetMusicBy(text string, music *model.Music) error {
	songId, err := parseSongIdByText(text)
	if err != nil {
		y.logger.Warn("error parseSongIdByText", logging.StringAttr("text", text), logging.ErrAttr(err))

		return err
	}

	song, err := y.youtubeapi.ReqSongInfoById(songId)
	if err != nil {
		y.logger.Warn("error ReqSongInfoById", logging.StringAttr("songId", songId), logging.ErrAttr(err))

		return err
	}

	if len(song.Items) < 1 {
		y.logger.Warn("song not found", logging.StringAttr("songId", songId))

		return fmt.Errorf("song not found")
	}

	music.Name = song.Items[0].Snippet.Title
	music.Link = songId
	music.Duration = formatTimeFromYoutube(song.Items[0].ContentDetails.Duration)

	return nil
}

// GetPlaylistBy implements MusicService.
func (y *youtube) GetPlaylistBy(link string, playlist *model.Playlist) error {
	panic("unimplemented")
}
