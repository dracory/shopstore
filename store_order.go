package shopstore

import (
	"context"
	"errors"
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

func (store *Store) OrderCount(ctx context.Context, options OrderQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.orderQuery(options)

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

func (store *Store) OrderCreate(ctx context.Context, order OrderInterface) error {
	order.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	order.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	order.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := order.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.orderTableName).
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

	order.MarkAsNotDirty()

	return nil
}

func (store *Store) OrderDelete(ctx context.Context, order OrderInterface) error {
	if order == nil {
		return errors.New("order is nil")
	}

	return store.OrderDeleteByID(ctx, order.ID())
}

func (store *Store) OrderDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.orderTableName).
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

func (store *Store) OrderSoftDelete(ctx context.Context, order OrderInterface) error {
	if order == nil {
		return errors.New("order is empty")
	}

	order.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.OrderUpdate(ctx, order)
}

func (store *Store) OrderSoftDeleteByID(ctx context.Context, id string) error {
	order, err := store.OrderFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.OrderSoftDelete(ctx, order)
}

func (store *Store) OrderFindByID(ctx context.Context, id string) (OrderInterface, error) {
	if id == "" {
		return nil, errors.New("order id is empty")
	}

	q := NewOrderQuery().SetID(id).SetLimit(1)
	list, err := store.OrderList(ctx, q)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) OrderList(ctx context.Context, options OrderQueryInterface) ([]OrderInterface, error) {
	q, columns, err := store.orderQuery(options)

	if err != nil {
		return []OrderInterface{}, err
	}

	sqlStr, _, errSql := q.Select(columns...).ToSQL()

	if errSql != nil {
		return []OrderInterface{}, nil
	}

	store.logSql("select", sqlStr)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)

	if err != nil {
		return []OrderInterface{}, err
	}

	list := []OrderInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewOrderFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) OrderUpdate(ctx context.Context, order OrderInterface) error {
	if order == nil {
		return errors.New("order is nil")
	}

	order.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := order.DataChanged()

	delete(dataChanged, "id")   // ID is not updateable
	delete(dataChanged, "hash") // Hash is not updateable
	delete(dataChanged, "data") // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.orderTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(order.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	order.MarkAsNotDirty()

	return err
}

func (store *Store) orderQuery(options OrderQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("order options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.orderTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasCustomerID() {
		q = q.Where(goqu.C(COLUMN_CUSTOMER_ID).Eq(options.CustomerID()))
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
		return q, columns, nil // soft deleted orders requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}

func (store *Store) OrderLineItemCount(ctx context.Context, options OrderLineItemQueryInterface) (int64, error) {
	q, _, err := store.orderLineItemQuery(options.SetCountOnly(true))

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

func (store *Store) OrderLineItemCreate(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	if orderLineItem == nil {
		return errors.New("orderLineItem is nil")
	}

	orderLineItem.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	orderLineItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	orderLineItem.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := orderLineItem.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.orderLineItemTableName).
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

	orderLineItem.MarkAsNotDirty()

	return nil
}

func (store *Store) OrderLineItemDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order line id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.orderLineItemTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

func (store *Store) OrderLineItemDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	return store.OrderLineItemDeleteByID(ctx, orderLineItem.ID())
}

func (store *Store) OrderLineItemFindByID(ctx context.Context, id string) (OrderLineItemInterface, error) {
	if id == "" {
		return nil, errors.New("order line id is empty")
	}

	list, err := store.OrderLineItemList(ctx, NewOrderLineItemQuery().
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

func (store *Store) OrderLineItemList(ctx context.Context, options OrderLineItemQueryInterface) ([]OrderLineItemInterface, error) {
	q, columns, err := store.orderLineItemQuery(options)

	if err != nil {
		return []OrderLineItemInterface{}, err
	}

	sqlStr, params, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return nil, errSql
	}

	store.logSql("select", sqlStr, params...)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return []OrderLineItemInterface{}, err
	}

	list := []OrderLineItemInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewOrderLineItemFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) OrderLineItemSoftDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	if orderLineItem == nil {
		return errors.New("order line is empty")
	}

	orderLineItem.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.OrderLineItemUpdate(ctx, orderLineItem)
}

func (store *Store) OrderLineItemSoftDeleteByID(ctx context.Context, id string) error {
	item, err := store.OrderLineItemFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.OrderLineItemSoftDelete(ctx, item)
}

func (store *Store) OrderLineItemUpdate(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	if orderLineItem == nil {
		return errors.New("orderLineItem is nil")
	}

	orderLineItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	dataChanged := orderLineItem.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.orderLineItemTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(orderLineItem.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	orderLineItem.MarkAsNotDirty()

	return err
}

func (store *Store) orderLineItemQuery(options OrderLineItemQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.orderLineItemTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasOrderID() {
		q = q.Where(goqu.C(COLUMN_ORDER_ID).Eq(options.OrderID()))
	}

	if options.HasOrderIDIn() {
		q = q.Where(goqu.C(COLUMN_ORDER_ID).In(options.OrderIDIn()))
	}

	if options.HasProductID() {
		q = q.Where(goqu.C(COLUMN_PRODUCT_ID).Eq(options.ProductID()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasStatusIn() {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn()))
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
		return q, columns, nil // soft deleted line items requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
