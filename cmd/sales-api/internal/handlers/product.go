package handlers

import (
	"github.com/daksh-sagar/garagesale/internal/platform/web"
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
func (p *Product) List(w http.ResponseWriter, _ *http.Request) {
	products, err := product.List(p.DB)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error fetching products", err)
		return
	}

	if err := web.Respond(w, products, http.StatusOK); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}

// Retrieve sends a single Product to the client as json response
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(p.DB, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error fetching product", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusOK); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}

// Create creates a product and sends the created product to the client
func (p *Product) Create(w http.ResponseWriter, r *http.Request) {
	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println(err)
		return
	}

	prod, err := product.Create(p.DB, &np, time.Now())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error creating product", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusCreated); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}
