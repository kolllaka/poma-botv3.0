package music

import "github.com/kolllaka/poma-botv3.0/internal/model"

type MusicService interface {
	GetMusicBy(link string, music *model.Music) error
	GetPlaylistBy(link string, playlist *model.Playlist) error
}
