package product

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("id provided not a valid UUID")
)

// List return all products
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT * FROM products`
	if err := db.SelectContext(ctx, &products, q); err != nil {
		return nil, err
	}

	return products, nil
}

// Retrieve returns a single product
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	var p Product
	const q = `SELECT * FROM products WHERE product_id = $1`

	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

// Create creates a new Product
func Create(ctx context.Context, db *sqlx.DB, np *NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}

	const q = `INSERT INTO products
		(product_id, name, cost, quantity, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", *np)
	}

	return &p, nil
}
