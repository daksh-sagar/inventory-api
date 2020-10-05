package handlers

import (
	"github.com/daksh-sagar/garagesale/internal/platform/web"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/daksh-sagar/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
)

// Product has handler methods
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List sends a list of Products to the client as json response
func (p *Product) List(w http.ResponseWriter, _ *http.Request) error {
	products, err := product.List(p.DB)

	if err != nil {
		return err
	}

	return web.Respond(w, products, http.StatusOK)
}

// Retrieve sends a single Product to the client as json response
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(p.DB, id)

	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for product %d", id)
		}
	}

	return web.Respond(w, prod, http.StatusOK)
}

// Create creates a product and sends the created product to the client
func (p *Product) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(p.DB, &np, time.Now())

	if err != nil {
		return err
	}

	return web.Respond(w, prod, http.StatusCreated)
}
