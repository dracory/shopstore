package shopstore

import (
	"context"
	"errors"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

func (store *Store) CategoryCount(ctx context.Context, options CategoryQueryInterface) (int64, error) {
	q, err := store.categoryQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (store *Store) CategoryCreate(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	category.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	category.SetSoftDeletedAt(MAX_DATETIME)

	data := category.Data()
	row := map[string]any{}
	for k, v := range data {
		row[k] = v
	}

	err := store.db.Query().Table(store.categoryTableName).Create(row)
	if err != nil {
		return err
	}

	return nil
}

func (store *Store) CategoryDelete(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	return store.CategoryDeleteByID(ctx, category.GetID())
}

func (store *Store) CategoryDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	_, err := store.db.Query().Table(store.categoryTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}

func (store *Store) CategoryFindByID(ctx context.Context, id string) (CategoryInterface, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	q := NewCategoryQuery().SetID(id).SetLimit(1)

	list, err := store.CategoryList(ctx, q)

	if err != nil {
		return nil, err
	}

	if len(list) < 1 {
		return nil, nil
	}

	return list[0], nil
}

func (store *Store) CategoryList(ctx context.Context, options CategoryQueryInterface) ([]CategoryInterface, error) {
	err := options.Validate()

	if err != nil {
		return nil, err
	}

	q, err := store.categoryQuery(options)

	if err != nil {
		return nil, err
	}

	var results []map[string]any
	if err := q.Get(&results); err != nil {
		return nil, err
	}

	list := []CategoryInterface{}

	lo.ForEach(results, func(result map[string]any, index int) {
		model := NewCategoryFromExistingData(mapAnyToString(result))
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) CategorySoftDelete(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.CategoryUpdate(ctx, category)
}

func (store *Store) CategorySoftDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	category, err := store.CategoryFindByID(ctx, id)

	if err != nil {
		return err
	}

	if category == nil {
		return nil
	}

	return store.CategorySoftDelete(ctx, category)
}

func (store *Store) CategoryUpdate(ctx context.Context, category CategoryInterface) (err error) {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := category.DataChanged()

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

	_, err = store.db.Query().Table(store.categoryTableName).Where(COLUMN_ID+" = ?", category.GetID()).Update(row)

	if err != nil {
		return err
	}

	category.MarkAsNotDirty()

	return nil
}

func (store *Store) categoryQuery(options CategoryQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		return nil, errors.New("category options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.categoryTableName)

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

	if options.HasParentID() {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID())
	}

	if options.HasStatus() {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.HasTitleLike() {
		searchTerm := strings.ReplaceAll(options.TitleLike(), "'", "''")
		searchTerm = strings.ReplaceAll(searchTerm, "%", "\\%")
		searchTerm = strings.ReplaceAll(searchTerm, "_", "\\_")
		q = q.Where(COLUMN_TITLE+" LIKE ?", "%"+searchTerm+"%")
	}

	if options.SoftDeletedIncluded() {
		return q, nil // soft deleted blocks requested specifically
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", MAX_DATETIME)

	return q, nil
}
