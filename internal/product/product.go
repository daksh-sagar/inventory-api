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
