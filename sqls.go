package shopstore

import (
	contractsschema "github.com/dracory/neat/contracts/database/schema"
)

// migration_001_product_table_add_parent_id adds the parent_id column to the product table if it doesn't exist.
// This handles existing tables that were created before parent_id was added to the schema.
func migration_001_product_table_add_parent_id(store *Store) error {
	if store.db.Schema().HasColumn(store.productTableName, COLUMN_PARENT_ID) {
		return nil
	}

	return store.db.Schema().Table(store.productTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_PARENT_ID, 40).Default("0")
	})
}

// migration_002_product_table_add_variant_dimensions adds variant matrix columns to the product table if they don't exist.
// This handles existing tables that were created before variant columns were added to the schema.
func migration_002_product_table_add_variant_dimensions(store *Store) error {
	if !store.db.Schema().HasColumn(store.productTableName, COLUMN_VARIANT_MATRIX_SCHEMA) {
		err := store.db.Schema().Table(store.productTableName, func(table contractsschema.Blueprint) {
			table.Text(COLUMN_VARIANT_MATRIX_SCHEMA).Default("{}")
		})
		if err != nil {
			return err
		}
	}

	if !store.db.Schema().HasColumn(store.productTableName, COLUMN_VARIANT_MATRIX_VALUES) {
		err := store.db.Schema().Table(store.productTableName, func(table contractsschema.Blueprint) {
			table.Text(COLUMN_VARIANT_MATRIX_VALUES).Default("{}")
		})
		if err != nil {
			return err
		}
	}

	return nil
}
