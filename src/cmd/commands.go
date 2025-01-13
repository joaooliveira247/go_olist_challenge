package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveira247/go_olist_challenge/src/config"
	"github.com/joaooliveira247/go_olist_challenge/src/db"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
	"github.com/joaooliveira247/go_olist_challenge/src/routes"
	"github.com/joaooliveira247/go_olist_challenge/src/utils"
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

func importAuthorsFromCSV(_ context.Context, cmd *cli.Command) error {
	header := cmd.Bool("header")
	path := cmd.Args().Get(0)

	authors, err := utils.ParseAuthorsFromCSV(path, header)

	if err != nil {
		return err
	}

	gormDB, err := db.GetDBConnection()

	if err != nil {
		return err
	}

	repository := repositories.NewAuthorRepository(gormDB)

	IDs, err := repository.CreateMany(&authors)

	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Inserted with IDs: %s", IDs))

	return nil
}

func Gen() *cli.Command {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Start API",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "port",
						Value:   uint64(8000),
						Aliases: []string{"p"},
						Usage:   "Port that API will run",
					},
				},
				Action: runAPI,
			},
			{
				Name:    "database",
				Aliases: []string{"db"},
				Usage:   "Interact with database",
				Flags:   nil,
				Commands: []*cli.Command{
					{
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "Create all tables",
						Action:  createTables,
					},
					{
						Name:    "delete",
						Aliases: []string{"d"},
						Usage:   "Delete all tables",
						Action:  deleteTables,
					},
				},
			},
			{
				Name:  "import",
				Usage: "Import authors by csv",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "header",
						Value:   true,
						Aliases: []string{"h"},
						Usage:   "Define csv has header",
					},
				},
				Action: importAuthorsFromCSV,
			},
		},
	}
	return cmd
}
