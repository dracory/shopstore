package shopstore

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (store *Store) ProductCount(ctx context.Context, options ProductQueryInterface) (int64, error) {
	q, _, err := store.productQuery(options.SetCountOnly(true))

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	store.logSql("count", sqlStr, params...)

	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (store *Store) ProductCreate(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is nil")
	}

	product.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	product.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	product.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := product.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.productTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("insert", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	product.MarkAsNotDirty()

	return nil
}

func (store *Store) ProductDelete(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is nil")
	}

	return store.ProductDeleteByID(ctx, product.ID())
}

func (store *Store) ProductDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.productTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

func (store *Store) ProductSoftDelete(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is empty")
	}

	product.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.ProductUpdate(ctx, product)
}

func (store *Store) ProductSoftDeleteByID(ctx context.Context, id string) error {
	product, err := store.ProductFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.ProductSoftDelete(ctx, product)
}

func (store *Store) ProductFindByID(ctx context.Context, id string) (ProductInterface, error) {
	if id == "" {
		return nil, errors.New("product id is empty")
	}

	list, err := store.ProductList(ctx, NewProductQuery().
		SetID(id).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) ProductList(ctx context.Context, options ProductQueryInterface) ([]ProductInterface, error) {
	q, columns, err := store.productQuery(options)

	if err != nil {
		return []ProductInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []ProductInterface{}, nil
	}

	store.logSql("select", sqlStr)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, sqlParams...)
	if err != nil {
		return []ProductInterface{}, err
	}

	list := []ProductInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewProductFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) ProductUpdate(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is nil")
	}

	product.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := product.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.productTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(product.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	product.MarkAsNotDirty()

	return err
}

func (store *Store) productQuery(options ProductQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("product options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.productTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasTitleLike() {
		// Sanitize the input by escaping special characters
		// This is a basic example, a more robust solution should be used in production
		searchTerm := strings.ReplaceAll(options.TitleLike(), "'", "''")
		searchTerm = strings.ReplaceAll(searchTerm, "%", "%")
		searchTerm = strings.ReplaceAll(searchTerm, "_", "_")
		q = q.Where(goqu.C(COLUMN_TITLE).ILike(`%` + searchTerm + `%`))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasStatusIn() {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn()))
	}

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Between(exp.NewRangeVal(options.CreatedAtGte(), options.CreatedAtLte())))
	} else if options.HasCreatedAtGte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()))
	} else if options.HasCreatedAtLte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()))
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(cast.ToUint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(cast.ToUint(options.Offset()))
		}
	}

	sortOrder := lo.Ternary(options.HasSortDirection(), options.SortDirection(), sb.DESC)

	if options.HasOrderBy() {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted products requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
