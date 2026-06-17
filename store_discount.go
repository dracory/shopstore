package shopstore

import (
	"context"
	"errors"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (store *Store) DiscountCount(ctx context.Context, options DiscountQueryInterface) (int64, error) {
	q, err := store.discountQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (store *Store) DiscountCreate(ctx context.Context, discount DiscountInterface) error {
	discount.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	discount.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	discount.SetSoftDeletedAt(MAX_DATETIME)

	data := discount.Data()
	row := map[string]any{}
	for k, v := range data {
		row[k] = v
	}

	err := store.db.Query().Table(store.discountTableName).Create(row)
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

	return store.DiscountDeleteByID(ctx, discount.GetID())
}

func (store *Store) DiscountDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("discount id is empty")
	}

	_, err := store.db.Query().Table(store.discountTableName).Where(COLUMN_ID+" = ?", id).Delete()
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
	q, err := store.discountQuery(options)
	if err != nil {
		return []DiscountInterface{}, err
	}

	var results []map[string]any
	if err := q.Get(&results); err != nil {
		return []DiscountInterface{}, err
	}

	list := []DiscountInterface{}

	lo.ForEach(results, func(result map[string]any, index int) {
		model := NewDiscountFromExistingData(mapAnyToString(result))
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

	row := map[string]any{}
	for k, v := range dataChanged {
		row[k] = v
	}

	_, err := store.db.Query().Table(store.discountTableName).Where(COLUMN_ID+" = ?", discount.GetID()).Update(row)

	discount.MarkAsNotDirty()

	return err
}

func (store *Store) discountQuery(options DiscountQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		options = NewDiscountQuery()
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.discountTableName)

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

	if options.HasCode() {
		q = q.Where(COLUMN_CODE+" = ?", options.Code())
	}

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" BETWEEN ? AND ?", options.CreatedAtGte(), options.CreatedAtLte())
	} else if options.HasCreatedAtGte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	} else if options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
	}

	if options.HasStartsAtGte() && options.HasStartsAtLte() {
		q = q.Where(COLUMN_STARTS_AT+" BETWEEN ? AND ?", options.StartsAtGte(), options.StartsAtLte())
	} else if options.HasStartsAtGte() {
		q = q.Where(COLUMN_STARTS_AT+" >= ?", options.StartsAtGte())
	} else if options.HasStartsAtLte() {
		q = q.Where(COLUMN_STARTS_AT+" <= ?", options.StartsAtLte())
	}

	if options.HasEndsAtGte() && options.HasEndsAtLte() {
		q = q.Where(COLUMN_ENDS_AT+" BETWEEN ? AND ?", options.EndsAtGte(), options.EndsAtLte())
	} else if options.HasEndsAtGte() {
		q = q.Where(COLUMN_ENDS_AT+" >= ?", options.EndsAtGte())
	} else if options.HasEndsAtLte() {
		q = q.Where(COLUMN_ENDS_AT+" <= ?", options.EndsAtLte())
	}

	if options.HasType() {
		q = q.Where(COLUMN_TYPE+" = ?", options.Type())
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
