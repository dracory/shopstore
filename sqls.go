package shopstore

import (
	"context"
	"fmt"
	"strings"
)

func migration_001_product_table_add_parent_id(store *Store) (err error) {
	db, _ := store.db.DB()
	_, err = db.ExecContext(context.Background(), "ALTER TABLE "+store.productTableName+" ADD COLUMN "+COLUMN_PARENT_ID+" TEXT DEFAULT '0'")
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}
	fmt.Println("Migration 001 completed: added parent_id column")
	return nil
}

func migration_002_product_table_add_variant_dimensions(store *Store) (err error) {
	db, _ := store.db.DB()

	// Add variant_matrix_schema column
	_, err = db.ExecContext(context.Background(), "ALTER TABLE "+store.productTableName+" ADD COLUMN "+COLUMN_VARIANT_MATRIX_SCHEMA+" TEXT DEFAULT '{}'")
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}
	fmt.Println("Migration 002a completed: added variant_matrix_schema column")

	// Add variant_matrix_values column
	_, err = db.ExecContext(context.Background(), "ALTER TABLE "+store.productTableName+" ADD COLUMN "+COLUMN_VARIANT_MATRIX_VALUES+" TEXT DEFAULT '{}'")
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}
	fmt.Println("Migration 002b completed: added variant_matrix_values column")

	return nil
}
