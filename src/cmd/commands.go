package main

import (
	"context"

	"github.com/joaooliveira247/go_olist_challenge/src/db"
	"github.com/urfave/cli/v3"
)

func createTables(_ context.Context, cmd *cli.Command) error {
	gormDB, err := db.GetDBConnection()

	if err != nil {
		return err
	}

	if err := db.CreateTables(gormDB); err != nil {
		return err
	}

	return nil
}