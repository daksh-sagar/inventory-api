package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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

	data, err := json.Marshal(products)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error marshalling data", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing", err)
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

	data, err := json.Marshal(prod)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error marshalling data", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing", err)
	}
}
