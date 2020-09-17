package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daksh-sagar/garagesale/schema"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	db, err := openDB()

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

	ps := ProductService{db: db}

	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ps.List),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "root"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

// Product is what we sell
type Product struct {
	ID          string    `db:"product_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"dateCreated"`
	DateUpdated time.Time `db:"date_updated" json:"dateUpdated"`
}

// ProductService has handler methods
type ProductService struct {
	db *sqlx.DB
}

// List sends a list of Products to the client as json response
func (p *ProductService) List(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	const q = `SELECT * FROM products`
	if err := p.db.Select(&products, q); err != nil {
		log.Printf("error: selecting products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(products)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error marshalling data", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	if _, err := w.Write(data); err != nil {
		log.Println("error writing", err)
	}
}
