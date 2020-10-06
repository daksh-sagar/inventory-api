package handlers

import (
	"log"
	"net/http"

	"github.com/daksh-sagar/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs a handler that knows about the routes
func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger)

	p := Product{DB: db, Log: logger}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)

	app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)

	return app
}
