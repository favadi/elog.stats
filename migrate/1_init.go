package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	// TODO: https://godoc.org/github.com/go-pg/migrations#Register
	// migrations.Register returns an error, handle it
	migrations.Register(func(db migrations.DB) error {
		queries := []string{
			`CREATE TABLE events (
				ip_client TEXT,
				ip_server TEXT,
				message   TEXT,
				tags 			JSONB
			)`,
			`CREATE INDEX events_ip_client_idx on events(ip_client)`,
			`CREATE INDEX events_ip_server_idx on events(ip_server)`,
			`CREATE INDEX events_tags_idx ON events USING gin (tags)`,
		}
		for _, query := range queries {
			_, err := db.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		_, err := db.Exec(`DROP TABLE events`)
		return err
	})
}
