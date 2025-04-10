package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
)

var (
	trueFiles = map[string]bool{
		".mp4":  true,
		".mp3":  true,
		".webm": true,
	}
)

type files struct {
	logger *logging.Logger
	path   string
}

func New(
	logger *logging.Logger,
	path string,
) *files {
	return &files{
		logger: logger,
		path:   path,
	}
}

// GetMusicBy implements MusicService.
func (f *files) GetMusicBy(name string) (model.Music, error) {
	panic("unimplemented")
}

// GetPlaylistBy implements MusicService.
func (f *files) GetPlaylistBy(name string) (model.Playlist, error) {
	var playlist model.Playlist
	var musics []*model.Music

	files, err := os.ReadDir(f.path)
	if err != nil {
		f.logger.Error("failed to open files", logging.StringAttr("path", f.path), logging.ErrAttr(err))

		return model.Playlist{}, err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if _, ok := trueFiles[ext]; !ok {
			continue
		}

		musics = append(musics, &model.Music{
			Name: file.Name(),
			Link: fmt.Sprintf("./audio/%s", file.Name()),
		})
	}

	playlist.Musics = musics

	f.logger.Debug("load playlist", logging.AnyAttr("musics", musics))

	return playlist, nil
}
