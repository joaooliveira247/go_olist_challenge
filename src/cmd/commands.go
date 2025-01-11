package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveira247/go_olist_challenge/src/config"
	"github.com/joaooliveira247/go_olist_challenge/src/db"
	"github.com/joaooliveira247/go_olist_challenge/src/routes"
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

func deleteTables(_ context.Context, cmd *cli.Command) error {
	gormDB, err := db.GetDBConnection()

	if err != nil {
		return err
	}

	if err := db.DeleteAllTables(gormDB); err != nil {
		return err
	}

	return nil
}

func runAPI(_ context.Context, cmd *cli.Command) error {
	api := gin.Default()
	routes.RegistryRoutes(api)
	port := config.APIPort
	if cliPort := cmd.Int("port"); cliPort > 0 {
		port = int(cliPort)
	}

	if err := api.Run(fmt.Sprintf(":%d", port)); err != nil {
		return err
	}
	return nil
}

func Gen() *[]cli.Command {
	cmd := &[]cli.Command{
		{
			Name:  "run",
			Usage: "Start API",
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:  "port",
					Value: uint64(8000),
					Usage: "Port that API will run",
				},
			},
			Action: runAPI,
		},
	}
	return cmd
}
