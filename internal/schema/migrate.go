package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in production
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add products",
		Script: `
CREATE TABLE products (
	product_id   UUID,
	name         TEXT,
	cost         INT,
	quantity     INT,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,
	PRIMARY KEY (product_id)
);`,
	},
	{
		Version:     2,
		Description: "Add sales",
		Script: `
CREATE TABLE sales (
	sale_id      UUID,
	product_id   UUID,
	quantity     INT,
	paid         INT,
	date_created TIMESTAMP,
	PRIMARY KEY (sale_id),
	FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {

	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
