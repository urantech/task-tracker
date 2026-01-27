package migrations

import (
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var EmbedMigrations embed.FS

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Run() {
	goose.SetBaseFS(EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	log.Println("Running migrations...")

	if err := goose.Up(m.db, "."); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations completed successfully")
}
