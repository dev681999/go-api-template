package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE "users" (
				"id" bigserial,
				"email" text NOT NULL UNIQUE,
				"name" text NOT NULL,
				"password" text NOT NULL,
				"active" boolean NOT NULL,
				"created_at" timestamptz NOT NULL,
				"updated_at" timestamptz NOT NULL,
				"deleted_at" timestamptz,
				PRIMARY KEY ("id"),
				UNIQUE ("email")
			)
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`
			DROP TABLE "users";
		`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20201217061014_initial", up, down, opts)
}
