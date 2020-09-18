package main

import (
	"flag"
	"log"

	"github.com/daksh-sagar/garagesale/internal/platform/database"
	"github.com/daksh-sagar/garagesale/internal/schema"
)

func main() {
	db, err := database.Open()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("applying migrations", err)
		}
		log.Println("Migrations complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("applying seed data", err)
		}
		log.Println("Seed data inserted")
		return
	}

}
