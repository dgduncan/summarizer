package summarizer

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
)

//go:embed queries/tables.sql
var schemaFS embed.FS

func NewSQLiteDB() (*sql.DB, error) {
	_ = os.Remove("./foo.db")
	// Create database connection with WAL mode and a few non-default PRAGMAs.
	// Notes (mattn/go-sqlite3): these `_...` params are applied per-connection,
	// so keep the pool small unless you have a specific reason not to.
	dsn := "file:./foo.db" +
		"?_journal_mode=WAL" +
		"&_foreign_keys=on" +
		"&_busy_timeout=5000" +
		"&_synchronous=NORMAL" +
		"&_cache_size=-64000" +
		"&_temp_store=MEMORY" +
		"&_wal_autocheckpoint=1000" +
		"&_journal_size_limit=67108864" +
		"&_mmap_size=268435456"

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQLite performs best with a small pool; multiple concurrent writers will
	// contend on database-level locks.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Test the connection
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Load embedded SQL schema
	schemaSQL, err := schemaFS.ReadFile("queries/tables.sql")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to read embedded schema: %w", err)
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}
