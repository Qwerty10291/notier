package storage

import "github.com/jmoiron/sqlx"

func InitDb(connStr string) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("sqlite3", connStr)
	if err != nil {
		return nil, err
	}
	err = initNotesTable(conn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func initNotesTable(conn *sqlx.DB) error {
	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS notes(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT,
		created_at INTEGER
		)`)
	return err
}
