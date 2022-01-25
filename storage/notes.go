package storage

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	Id        int    `json:"id" db:"id"`
	Text      string `json:"text" db:"text"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}

var ErrorNoteNotFound = errors.New("note not found")

func GetAllNotes(conn *sqlx.DB) ([]Note, error) {
	rows := []Note{}
	err := conn.Select(&rows, `SELECT * FROM notes`)
	return rows, err
}

func NewNote(conn *sqlx.DB, text string) (*Note, error) {
	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	timestamp := time.Now().Unix()
	result, err := tx.Exec(`INSERT INTO notes (text, created_at) VALUES (?, ?)`, text, timestamp)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	noteId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Note{
		Id:        int(noteId),
		Text:      text,
		CreatedAt: timestamp,
	}, nil
}

func DeleteNote(conn *sqlx.DB, noteId int) error {
	result, err := conn.Exec("DELETE FROM notes where id = $1", noteId)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrorNoteNotFound
	}
	return nil
}
