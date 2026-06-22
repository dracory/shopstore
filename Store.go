package shopstore

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/dracory/neat"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
)

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	categoryTableName      string
	discountTableName      string
	mediaTableName         string
	orderTableName         string
	orderLineItemTableName string
	productTableName       string
	db                     *neat.Database
	timeoutSeconds         int64
	automigrateEnabled     bool
	debugEnabled           bool
	sqlLogger              *slog.Logger
}

// logSql logs sql to the sql logger
func (store *Store) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	}
}

// MigrateUp creates or updates database tables to match the current schema
func (store *Store) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if err := store.categoryTableCreate(); err != nil {
		return err
	}
	if err := store.discountTableCreate(); err != nil {
		return err
	}
	if err := store.mediaTableCreate(); err != nil {
		return err
	}
	if err := store.orderTableCreate(); err != nil {
		return err
	}
	if err := store.orderLineItemTableCreate(); err != nil {
		return err
	}
	if err := store.productTableCreate(); err != nil {
		return err
	}

	if err := migration_001_product_table_add_parent_id(store); err != nil {
		return err
	}

	if err := migration_002_product_table_add_variant_dimensions(store); err != nil {
		return err
	}

	return nil
}

// MigrateDown drops the shop store tables
func (store *Store) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	_ = store.db.Schema().DropIfExists(store.categoryTableName)
	_ = store.db.Schema().DropIfExists(store.discountTableName)
	_ = store.db.Schema().DropIfExists(store.mediaTableName)
	_ = store.db.Schema().DropIfExists(store.orderLineItemTableName)
	_ = store.db.Schema().DropIfExists(store.orderTableName)
	_ = store.db.Schema().DropIfExists(store.productTableName)
	return nil
}

func (store *Store) DB() *sql.DB {
	db, _ := store.db.DB()
	return db
}

// EnableDebug - enables the debug option
func (store *Store) EnableDebug(debug bool, sqlLogger ...*slog.Logger) {
	store.debugEnabled = debug
	if store.debugEnabled {
		if len(sqlLogger) < 1 {
			store.sqlLogger = slog.Default()
			return
		}
		store.sqlLogger = sqlLogger[0]
	} else {
		store.sqlLogger = nil
	}
}

func (store *Store) CategoryTableName() string {
	return store.categoryTableName
}

func (store *Store) DiscountTableName() string {
	return store.discountTableName
}

func (store *Store) MediaTableName() string {
	return store.mediaTableName
}

func (store *Store) OrderTableName() string {
	return store.orderTableName
}

func (store *Store) OrderLineItemTableName() string {
	return store.orderLineItemTableName
}

func (store *Store) ProductTableName() string {
	return store.productTableName
}

func (store *Store) categoryTableCreate() error {
	if store.db.Schema().HasTable(store.categoryTableName) {
		return nil
	}
	return store.db.Schema().Create(store.categoryTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_STATUS, 20)
		table.String(COLUMN_PARENT_ID, 40)
		table.String(COLUMN_TITLE, 255)
		table.Text(COLUMN_DESCRIPTION)
		table.Text(COLUMN_METAS)
		table.Text(COLUMN_MEMO)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})
}

func (store *Store) discountTableCreate() error {
	if store.db.Schema().HasTable(store.discountTableName) {
		return nil
	}
	return store.db.Schema().Create(store.discountTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_STATUS, 20)
		table.String(COLUMN_TITLE, 255)
		table.Text(COLUMN_DESCRIPTION)
		table.String(COLUMN_TYPE, 20)
		table.Decimal(COLUMN_AMOUNT)
		table.String(COLUMN_CODE, 100)
		table.DateTime(COLUMN_STARTS_AT)
		table.DateTime(COLUMN_ENDS_AT)
		table.Text(COLUMN_METAS)
		table.Text(COLUMN_MEMO)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})
}

func (store *Store) mediaTableCreate() error {
	if store.db.Schema().HasTable(store.mediaTableName) {
		return nil
	}
	return store.db.Schema().Create(store.mediaTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_STATUS, 20)
		table.String(COLUMN_ENTITY_ID, 40)
		table.Integer(COLUMN_SEQUENCE)
		table.String(COLUMN_MEDIA_TYPE, 50)
		table.String(COLUMN_MEDIA_URL, 510)
		table.String(COLUMN_TITLE, 255)
		table.Text(COLUMN_DESCRIPTION)
		table.Text(COLUMN_MEMO)
		table.Text(COLUMN_METAS)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})
}

func (store *Store) orderTableCreate() error {
	if store.db.Schema().HasTable(store.orderTableName) {
		return nil
	}
	return store.db.Schema().Create(store.orderTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_STATUS, 20)
		table.String(COLUMN_CUSTOMER_ID, 40)
		table.Integer(COLUMN_QUANTITY)
		table.Decimal(COLUMN_PRICE)
		table.Text(COLUMN_METAS)
		table.Text(COLUMN_MEMO)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})
}

func (store *Store) orderLineItemTableCreate() error {
	if store.db.Schema().HasTable(store.orderLineItemTableName) {
		return nil
	}
	return store.db.Schema().Create(store.orderLineItemTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_STATUS, 20)
		table.String(COLUMN_ORDER_ID, 40)
		table.String(COLUMN_PRODUCT_ID, 40)
		table.String(COLUMN_TITLE, 255)
		table.Integer(COLUMN_QUANTITY)
		table.Decimal(COLUMN_PRICE)
		table.Text(COLUMN_METAS)
		table.Text(COLUMN_MEMO)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})
}

func (store *Store) productTableCreate() error {
	if store.db.Schema().HasTable(store.productTableName) {
		return nil
	}
	return store.db.Schema().Create(store.productTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_STATUS, 20)
		table.String(COLUMN_PARENT_ID, 40)
		table.String(COLUMN_TITLE, 255)
		table.Text(COLUMN_DESCRIPTION)
		table.Text(COLUMN_SHORT_DESCRIPTION)
		table.Integer(COLUMN_QUANTITY)
		table.Decimal(COLUMN_PRICE)
		table.Text(COLUMN_VARIANT_MATRIX_SCHEMA)
		table.Text(COLUMN_VARIANT_MATRIX_VALUES)
		table.Text(COLUMN_METAS)
		table.Text(COLUMN_MEMO)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})
}
