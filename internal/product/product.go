package product

import (
	"github.com/jmoiron/sqlx"
)

// List return all products
func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT * FROM products`
	if err := db.Select(&products, q); err != nil {
		return nil, err
	}

	return products, nil
}

// Retrieve returns a single product
func Retrieve(db *sqlx.DB, id string) (*Product, error) {
	var p Product
	const q = `SELECT * FROM products WHERE product_id = $1`

	if err := db.Get(&p, q, id); err != nil {
		return nil, err
	}

	return &p, nil
}
