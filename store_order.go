package shopstore

import (
	"context"
	"errors"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (store *Store) OrderCount(ctx context.Context, options OrderQueryInterface) (int64, error) {
	q, err := store.orderQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (store *Store) OrderCreate(ctx context.Context, order OrderInterface) error {
	order.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	order.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	order.SetSoftDeletedAt(MAX_DATETIME)

	data := order.Data()
	row := map[string]any{}
	for k, v := range data {
		row[k] = v
	}

	err := store.db.Query().Table(store.orderTableName).Create(row)
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

	return store.OrderDeleteByID(ctx, order.GetID())
}

func (store *Store) OrderDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order id is empty")
	}

	_, err := store.db.Query().Table(store.orderTableName).Where(COLUMN_ID+" = ?", id).Delete()
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
	q, err := store.orderQuery(options)
	if err != nil {
		return []OrderInterface{}, err
	}

	var results []map[string]any
	if err := q.Get(&results); err != nil {
		return []OrderInterface{}, err
	}

	list := []OrderInterface{}

	lo.ForEach(results, func(result map[string]any, index int) {
		model := NewOrderFromExistingData(mapAnyToString(result))
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

	row := map[string]any{}
	for k, v := range dataChanged {
		row[k] = v
	}

	_, err := store.db.Query().Table(store.orderTableName).Where(COLUMN_ID+" = ?", order.GetID()).Update(row)

	order.MarkAsNotDirty()

	return err
}

func (store *Store) orderQuery(options OrderQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		return nil, errors.New("order options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.orderTableName)

	if options.HasID() {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasIDIn() {
		ids := make([]any, len(options.IDIn()))
		for i, id := range options.IDIn() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.HasCustomerID() {
		q = q.Where(COLUMN_CUSTOMER_ID+" = ?", options.CustomerID())
	}

	if options.HasStatus() {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.HasStatusIn() {
		statuses := make([]any, len(options.StatusIn()))
		for i, status := range options.StatusIn() {
			statuses[i] = status
		}
		q = q.WhereIn(COLUMN_STATUS, statuses)
	}

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" BETWEEN ? AND ?", options.CreatedAtGte(), options.CreatedAtLte())
	} else if options.HasCreatedAtGte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	} else if options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(cast.ToInt(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(cast.ToInt(options.Offset()))
		}
	}

	sortOrder := lo.Ternary(options.HasSortDirection(), options.SortDirection(), "desc")

	if options.HasOrderBy() {
		q = q.OrderBy(options.OrderBy(), sortOrder)
	}

	if !options.SoftDeletedIncluded() {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", MAX_DATETIME)
	}

	return q, nil
}

func (store *Store) OrderLineItemCount(ctx context.Context, options OrderLineItemQueryInterface) (int64, error) {
	q, err := store.orderLineItemQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (store *Store) OrderLineItemCreate(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	if orderLineItem == nil {
		return errors.New("orderLineItem is nil")
	}

	orderLineItem.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	orderLineItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	orderLineItem.SetSoftDeletedAt(MAX_DATETIME)

	data := orderLineItem.Data()
	row := map[string]any{}
	for k, v := range data {
		row[k] = v
	}

	err := store.db.Query().Table(store.orderLineItemTableName).Create(row)
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

	_, err := store.db.Query().Table(store.orderLineItemTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}

func (store *Store) OrderLineItemDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	return store.OrderLineItemDeleteByID(ctx, orderLineItem.GetID())
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
	q, err := store.orderLineItemQuery(options)
	if err != nil {
		return []OrderLineItemInterface{}, err
	}

	var results []map[string]any
	if err := q.Get(&results); err != nil {
		return []OrderLineItemInterface{}, err
	}

	list := []OrderLineItemInterface{}

	lo.ForEach(results, func(result map[string]any, index int) {
		model := NewOrderLineItemFromExistingData(mapAnyToString(result))
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

	row := map[string]any{}
	for k, v := range dataChanged {
		row[k] = v
	}

	_, err := store.db.Query().Table(store.orderLineItemTableName).Where(COLUMN_ID+" = ?", orderLineItem.GetID()).Update(row)

	orderLineItem.MarkAsNotDirty()

	return err
}

func (store *Store) orderLineItemQuery(options OrderLineItemQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		return nil, errors.New("options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.orderLineItemTableName)

	if options.HasID() {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasIDIn() {
		ids := make([]any, len(options.IDIn()))
		for i, id := range options.IDIn() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.HasOrderID() {
		q = q.Where(COLUMN_ORDER_ID+" = ?", options.OrderID())
	}

	if options.HasOrderIDIn() {
		ids := make([]any, len(options.OrderIDIn()))
		for i, id := range options.OrderIDIn() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ORDER_ID, ids)
	}

	if options.HasProductID() {
		q = q.Where(COLUMN_PRODUCT_ID+" = ?", options.ProductID())
	}

	if options.HasStatus() {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.HasStatusIn() {
		statuses := make([]any, len(options.StatusIn()))
		for i, status := range options.StatusIn() {
			statuses[i] = status
		}
		q = q.WhereIn(COLUMN_STATUS, statuses)
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(cast.ToInt(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(cast.ToInt(options.Offset()))
		}
	}

	sortOrder := lo.Ternary(options.HasSortDirection(), options.SortDirection(), "desc")

	if options.HasOrderBy() {
		q = q.OrderBy(options.OrderBy(), sortOrder)
	}

	if !options.SoftDeletedIncluded() {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", MAX_DATETIME)
	}

	return q, nil
}
