package main

import (
	"context"
	"log"
	"os"

	"github.com/joaooliveira247/go_olist_challenge/src/config"
)

func init() {
	config.LoadEnv()
}

func main() {
	app := Gen()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
