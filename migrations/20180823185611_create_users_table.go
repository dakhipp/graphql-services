package main

import (
	"github.com/go-pg/pg/orm"
	"github.com/robinjoseph08/go-pg-migrations"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
      CREATE TABLE IF NOT EXISTS users (
        id 					    TEXT 		PRIMARY KEY,
        first_name 	    TEXT 	  NOT NULL,
        last_name 	    TEXT 	  NOT NULL,
        email           TEXT    NOT NULL UNIQUE,
        phone           TEXT    NOT NULL UNIQUE,
        password        TEXT    NOT NULL,
        phone_verified  BOOL    DEFAULT FALSE,
        email_verified  BOOL    DEFAULT FALSE,
        roles           JSONB
      );
    `)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE users")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20180823185611_create_users_table", up, down, opts)
}
