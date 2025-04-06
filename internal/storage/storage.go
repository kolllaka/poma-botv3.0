package storage

type Store interface {
	GetDuration(music *StoreDuration) error
	StoreDuration(music *StoreDuration) error
}
