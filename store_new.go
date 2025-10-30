package shopstore

import (
	"database/sql"
	"errors"

	"github.com/dracory/sb"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	CategoryTableName      string
	DiscountTableName      string
	MediaTableName         string
	OrderTableName         string
	OrderLineItemTableName string
	ProductTableName       string
	DB                     *sql.DB
	DbDriverName           string
	AutomigrateEnabled     bool
	DebugEnabled           bool
}

// NewStore creates a new block store
func NewStore(opts NewStoreOptions) (*Store, error) {
	if opts.CategoryTableName == "" {
		return nil, errors.New("shop store: CategoryTableName is required")
	}

	if opts.DiscountTableName == "" {
		return nil, errors.New("shop store: DiscountTableName is required")
	}

	if opts.MediaTableName == "" {
		return nil, errors.New("shop store: MediaTableName is required")
	}

	if opts.OrderTableName == "" {
		return nil, errors.New("shop store: OrderTableName is required")
	}

	if opts.OrderLineItemTableName == "" {
		return nil, errors.New("shop store: OrderLineItemTableName is required")
	}

	if opts.ProductTableName == "" {
		return nil, errors.New("shop store: ProductTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("shop store: DB is required")
	}

	if opts.DbDriverName == "" {
		opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
	}

	store := &Store{
		categoryTableName:      opts.CategoryTableName,
		discountTableName:      opts.DiscountTableName,
		mediaTableName:         opts.MediaTableName,
		orderTableName:         opts.OrderTableName,
		orderLineItemTableName: opts.OrderLineItemTableName,
		productTableName:       opts.ProductTableName,
		automigrateEnabled:     opts.AutomigrateEnabled,
		db:                     opts.DB,
		dbDriverName:           opts.DbDriverName,
		debugEnabled:           opts.DebugEnabled,
	}

	store.timeoutSeconds = 2 * 60 * 60 // 2 hours

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
