package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/daksh-sagar/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
)

// Product has handler methods
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List sends a list of Products to the client as json response
func (p *Product) List(w http.ResponseWriter, r *http.Request) {
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
