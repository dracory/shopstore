package shopstore

import (
	"context"
	"fmt"

	"github.com/dracory/sb"
)

// sqlCategoryTableCreate returns a SQL string for creating the category table
func (store *Store) sqlCategoryTableCreate() (string, error) {
	sql, err := sb.NewBuilder(store.dbDriverName).
		Table(store.categoryTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:   COLUMN_PARENT_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_DESCRIPTION,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, err
}

// sqlDiscountTableCreate returns a SQL string for creating the discount table
func (store *Store) sqlDiscountTableCreate() (string, error) {
	sql, err := sb.NewBuilder(store.dbDriverName).
		Table(store.discountTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_DESCRIPTION,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:   COLUMN_TYPE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:     COLUMN_AMOUNT,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 2,
		}).
		Column(sb.Column{
			Name:   COLUMN_CODE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name: COLUMN_STARTS_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_ENDS_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, err
}

func (store *Store) sqlMediaTableCreate() (string, error) {
	sql, err := sb.NewBuilder(store.dbDriverName).
		Table(store.mediaTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:   COLUMN_ENTITY_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name: COLUMN_SEQUENCE,
			Type: sb.COLUMN_TYPE_INTEGER,
		}).
		Column(sb.Column{
			Name:   COLUMN_MEDIA_TYPE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:   COLUMN_MEDIA_URL,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 510,
		}).
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_DESCRIPTION,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, err
}

func (store *Store) sqlOrderLineItemTableCreate() (string, error) {
	sql, err := sb.NewBuilder(store.dbDriverName).
		Table(store.orderLineItemTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:   COLUMN_ORDER_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_PRODUCT_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   COLUMN_QUANTITY,
			Type:   sb.COLUMN_TYPE_INTEGER,
			Length: 10,
		}).
		Column(sb.Column{
			Name:     COLUMN_PRICE,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 2,
		}).
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, err
}

func (store *Store) sqlOrderTableCreate() (string, error) {
	sql, err := sb.NewBuilder(store.dbDriverName).
		Table(store.orderTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:   COLUMN_CUSTOMER_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_QUANTITY,
			Type:   sb.COLUMN_TYPE_INTEGER,
			Length: 10,
		}).
		Column(sb.Column{
			Name:     COLUMN_PRICE,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 2,
		}).
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, err
}

func (store *Store) sqlProductTableCreate() (string, error) {
	sql, err := sb.NewBuilder(store.dbDriverName).
		Table(store.productTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:    COLUMN_PARENT_ID,
			Type:    sb.COLUMN_TYPE_STRING,
			Length:  40,
			Default: "0",
		}).
		Column(sb.Column{
			Name:   COLUMN_TITLE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_DESCRIPTION,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_SHORT_DESCRIPTION,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:   COLUMN_QUANTITY,
			Type:   sb.COLUMN_TYPE_INTEGER,
			Length: 10,
		}).
		Column(sb.Column{
			Name:     COLUMN_PRICE,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 2,
		}).
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:    COLUMN_VARIANT_MATRIX_SCHEMA,
			Type:    sb.COLUMN_TYPE_TEXT,
			Default: "{}",
		}).
		Column(sb.Column{
			Name:    COLUMN_VARIANT_MATRIX_VALUES,
			Type:    sb.COLUMN_TYPE_TEXT,
			Default: "{}",
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, err
}

func migration_001_product_table_add_parent_id(store *Store) (err error) {
	builder := sb.NewBuilder(store.dbDriverName)
	sql, params, err := builder.TableColumnExists(store.productTableName, COLUMN_PARENT_ID)
	if err != nil {
		return err
	}

	row := store.db.QueryRowContext(context.Background(), sql, params...)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	// Column already exists, skip migration
	if count > 0 {
		return nil
	}

	sql2, err := builder.TableColumnAdd(store.productTableName, sb.Column{
		Name:    COLUMN_PARENT_ID,
		Type:    sb.COLUMN_TYPE_STRING,
		Length:  40,
		Default: "0",
	})
	if err != nil {
		return err
	}

	_, err = store.db.ExecContext(context.Background(), sql2)
	if err != nil {
		return err
	}

	fmt.Println("Migration 001 completed: added parent_id column")

	return nil
}

func migration_002_product_table_add_variant_dimensions(store *Store) (err error) {
	builder := sb.NewBuilder(store.dbDriverName)

	// Add variant_matrix_schema column
	sql, params, err := builder.TableColumnExists(store.productTableName, COLUMN_VARIANT_MATRIX_SCHEMA)
	if err != nil {
		return err
	}

	row := store.db.QueryRowContext(context.Background(), sql, params...)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		sql2, err := builder.TableColumnAdd(store.productTableName, sb.Column{
			Name:    COLUMN_VARIANT_MATRIX_SCHEMA,
			Type:    sb.COLUMN_TYPE_TEXT,
			Default: "{}",
		})
		if err != nil {
			return err
		}

		_, err = store.db.ExecContext(context.Background(), sql2)
		if err != nil {
			return err
		}
		fmt.Println("Migration 002a completed: added variant_matrix_schema column")
	}

	// Add variant_matrix_values column
	sql, params, err = builder.TableColumnExists(store.productTableName, COLUMN_VARIANT_MATRIX_VALUES)
	if err != nil {
		return err
	}

	row = store.db.QueryRowContext(context.Background(), sql, params...)
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		sql2, err := builder.TableColumnAdd(store.productTableName, sb.Column{
			Name:    COLUMN_VARIANT_MATRIX_VALUES,
			Type:    sb.COLUMN_TYPE_TEXT,
			Default: "{}",
		})
		if err != nil {
			return err
		}

		_, err = store.db.ExecContext(context.Background(), sql2)
		if err != nil {
			return err
		}
		fmt.Println("Migration 002b completed: added variant_matrix_values column")
	}

	return nil
}
