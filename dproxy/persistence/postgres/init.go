package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

func Connect(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := migrateDB(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "dproxy_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS events (
						init_time TIMESTAMP NOT NULL,
						initiator VARCHAR(254) NOT NULL
					)`,
					`CREATE INDEX IF NOT EXISTS events_idx ON events (
						initiator
					)`,
				},
				Down: []string{"DROP TABLE events"},
			},
		},
	}
	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	return err
}
