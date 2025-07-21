package shopstore

import (
	"context"
	"database/sql"
	"log"
	"log/slog"

	"github.com/gouniverse/base/database"
)

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	categoryTableName      string
	discountTableName      string
	mediaTableName         string
	orderTableName         string
	orderLineItemTableName string
	productTableName       string
	db                     *sql.DB
	dbDriverName           string
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

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	sqls := []string{
		store.sqlCategoryTableCreate(),
		store.sqlDiscountTableCreate(),
		store.sqlMediaTableCreate(),
		store.sqlOrderTableCreate(),
		store.sqlOrderLineItemTableCreate(),
		store.sqlProductTableCreate(),
	}

	for _, sql := range sqls {

		store.logSql("create table", sql)

		_, err := store.db.Exec(sql)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (store *Store) DB() *sql.DB {
	return store.db
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

func (store *Store) toQuerableContext(context context.Context) database.QueryableContext {
	if database.IsQueryableContext(context) {
		return context.(database.QueryableContext)
	}

	return database.Context(context, store.db)
}
