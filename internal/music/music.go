package music

import "github.com/kolllaka/poma-botv3.0/internal/model"

type MusicService interface {
	GetMusicBy(link string) (model.Music, error)
	GetPlaylistBy(link string) (model.Playlist, error)
}
