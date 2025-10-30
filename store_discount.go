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

func (store *Store) DiscountCount(ctx context.Context, options DiscountQueryInterface) (int64, error) {
	q, _, err := store.discountQuery(options.SetCountOnly(true))

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

	if store.debugEnabled {
		log.Println(sqlStr)
	}

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

func (store *Store) DiscountCreate(ctx context.Context, discount DiscountInterface) error {
	discount.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	discount.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	discount.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := discount.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.discountTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	discount.MarkAsNotDirty()

	return nil
}

func (store *Store) DiscountDelete(ctx context.Context, discount DiscountInterface) error {
	if discount == nil {
		return errors.New("discount is nil")
	}

	return store.DiscountDeleteByID(ctx, discount.ID())
}

func (store *Store) DiscountDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("discount id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.discountTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

func (store *Store) DiscountFindByID(ctx context.Context, id string) (DiscountInterface, error) {
	if id == "" {
		return nil, errors.New("discount id is empty")
	}

	list, err := store.DiscountList(ctx, NewDiscountQuery().
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

func (store *Store) DiscountFindByCode(ctx context.Context, code string) (DiscountInterface, error) {
	if code == "" {
		return nil, errors.New("discount code is empty")
	}

	list, err := store.DiscountList(ctx, NewDiscountQuery().
		SetStatus(DISCOUNT_STATUS_ACTIVE).
		SetCode(code).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) DiscountList(ctx context.Context, options DiscountQueryInterface) ([]DiscountInterface, error) {
	q, columns, err := store.discountQuery(options)

	if err != nil {
		return []DiscountInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []DiscountInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, sqlParams...)

	if err != nil {
		return []DiscountInterface{}, err
	}

	list := []DiscountInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewDiscountFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) DiscountSoftDelete(ctx context.Context, discount DiscountInterface) error {
	if discount == nil {
		return errors.New("discount is nil")
	}

	discount.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.DiscountUpdate(ctx, discount)
}

func (store *Store) DiscountSoftDeleteByID(ctx context.Context, id string) error {
	discount, err := store.DiscountFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.DiscountSoftDelete(ctx, discount)
}

func (store *Store) DiscountUpdate(ctx context.Context, discount DiscountInterface) error {
	if discount == nil {
		return errors.New("discount is nil")
	}

	discount.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := discount.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.discountTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(discount.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	discount.MarkAsNotDirty()

	return err
}

func (store *Store) discountQuery(options DiscountQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		options = NewDiscountQuery()
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.discountTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasStatusIn() {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn()))
	}

	if options.HasCode() {
		q = q.Where(goqu.C(COLUMN_CODE).Eq(options.Code()))
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
		return q, columns, nil // soft deleted discounts requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
