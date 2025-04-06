package sqlite

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/kolllaka/poma-botv3.0/internal/storage"
)

type durationStore struct {
	db *sql.DB
}

func New(db *sql.DB) *durationStore {
	return &durationStore{
		db: db,
	}
}

// GetDuration implements Store.
func (d *durationStore) GetDuration(music *storage.StoreDuration) error {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE name = ?", tableMusic)
	if err := d.db.QueryRow(stmt, filepath.Base(music.Link)).Scan(&music.Link, &music.Duration); err != nil {
		return err
	}

	return nil
}

// StoreDuration implements Store.
func (d *durationStore) StoreDuration(music *storage.StoreDuration) error {
	stmt := fmt.Sprintf("INSERT INTO %s (name, duration) VALUES(?, ?)", tableMusic)
	if _, err := d.db.Exec(
		stmt,
		filepath.Base(music.Link),
		music.Duration,
	); err != nil {
		return err
	}

	return nil
}
