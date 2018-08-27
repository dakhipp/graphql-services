package main

import (
	"github.com/go-pg/pg/orm"
	"github.com/robinjoseph08/go-pg-migrations"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
      CREATE TABLE IF NOT EXISTS users (
        id 					CHAR(27) 			PRIMARY KEY,
        firstName 	VARCHAR(24) 	NOT NULL,
        lastName 		VARCHAR(24) 	NOT NULL
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
