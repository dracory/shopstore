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

// migration_003_discount_table_add_max_uses adds the max_uses, max_uses_count,
// and max_uses_per_customer columns to the discount table if they don't exist.
// This handles existing tables created before usage-limit columns were added.
func migration_003_discount_table_add_max_uses(store *Store) error {
	if !store.db.Schema().HasColumn(store.discountTableName, COLUMN_MAX_USES) {
		err := store.db.Schema().Table(store.discountTableName, func(table contractsschema.Blueprint) {
			table.Integer(COLUMN_MAX_USES).Default(DEFAULT_MAX_USES)
		})
		if err != nil {
			return err
		}
	}

	if !store.db.Schema().HasColumn(store.discountTableName, COLUMN_MAX_USES_COUNT) {
		err := store.db.Schema().Table(store.discountTableName, func(table contractsschema.Blueprint) {
			table.Integer(COLUMN_MAX_USES_COUNT).Default(0)
		})
		if err != nil {
			return err
		}
	}

	if !store.db.Schema().HasColumn(store.discountTableName, COLUMN_MAX_USES_PER_CUSTOMER) {
		err := store.db.Schema().Table(store.discountTableName, func(table contractsschema.Blueprint) {
			table.Integer(COLUMN_MAX_USES_PER_CUSTOMER).Default(DEFAULT_MAX_USES_PER_CUSTOMER)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// migration_005_discount_table_metas_to_longtext changes the metas column
// from TEXT to LONGTEXT to accommodate growing JSON data such as per-customer
// redemption counts.
func migration_005_discount_table_metas_to_longtext(store *Store) error {
	if store.db.Schema().HasColumn(store.discountTableName, COLUMN_METAS) {
		err := store.db.Schema().Table(store.discountTableName, func(table contractsschema.Blueprint) {
			table.LongText(COLUMN_METAS).Default("{}").Change()
		})
		if err != nil {
			return err
		}
	}

	return nil
}
