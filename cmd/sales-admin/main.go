package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/daksh-sagar/garagesale/internal/platform/database"
	"github.com/daksh-sagar/garagesale/internal/schema"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := database.Open()

	if err != nil {
		return errors.Wrap(err, "opening database")
	}

	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "applying migrations")
		}
		fmt.Println("Migrations complete")
	case "seed":
		if err := schema.Seed(db); err != nil {
			errors.Wrap(err, "applying seed data")
		}
		fmt.Println("Seed data inserted")
	}
	return nil
}
