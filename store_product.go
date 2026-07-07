package shopstore

import (
	"context"
	"errors"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (store *Store) ProductCount(ctx context.Context, options ProductQueryInterface) (int64, error) {
	q, err := store.productQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (store *Store) ProductCreate(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is nil")
	}

	product.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	product.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	product.SetSoftDeletedAt(MAX_DATETIME)

	data := product.Data()
	row := map[string]any{}
	for k, v := range data {
		row[k] = v
	}

	err := store.db.Query().Table(store.productTableName).Create(row)
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

	return store.ProductDeleteByID(ctx, product.GetID())
}

func (store *Store) ProductDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product id is empty")
	}

	_, err := store.db.Query().Table(store.productTableName).Where(COLUMN_ID+" = ?", id).Delete()
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
	q, err := store.productQuery(options)
	if err != nil {
		return []ProductInterface{}, err
	}

	var results []map[string]any
	if err := q.Get(&results); err != nil {
		return []ProductInterface{}, err
	}

	list := []ProductInterface{}

	lo.ForEach(results, func(result map[string]any, index int) {
		model := NewProductFromExistingData(mapAnyToString(result))
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

	row := map[string]any{}
	for k, v := range dataChanged {
		row[k] = v
	}

	_, err := store.db.Query().Table(store.productTableName).Where(COLUMN_ID+" = ?", product.GetID()).Update(row)

	product.MarkAsNotDirty()

	return err
}

func (store *Store) productQuery(options ProductQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		return nil, errors.New("product options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.productTableName)

	if options.HasID() {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasNotID() {
		q = q.Where(COLUMN_ID+" != ?", options.NotID())
	}

	if options.HasIDIn() {
		ids := make([]any, len(options.IDIn()))
		for i, id := range options.IDIn() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.HasNotIDIn() {
		ids := make([]any, len(options.NotIDIn()))
		for i, id := range options.NotIDIn() {
			ids[i] = id
		}
		q = q.WhereNotIn(COLUMN_ID, ids)
	}

	if options.HasTitleLike() {
		searchTerm := strings.ReplaceAll(options.TitleLike(), "'", "''")
		searchTerm = strings.ReplaceAll(searchTerm, "%", "\\%")
		searchTerm = strings.ReplaceAll(searchTerm, "_", "\\_")
		q = q.Where(COLUMN_TITLE+" LIKE ?", "%"+searchTerm+"%")
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

	if options.HasParentID() {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID())
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

// ProductVariantList returns all variants for a given parent product ID
func (store *Store) ProductVariantList(ctx context.Context, parentID string) ([]ProductInterface, error) {
	if parentID == "" {
		return []ProductInterface{}, errors.New("parent ID is empty")
	}

	return store.ProductList(ctx, NewProductQuery().SetParentID(parentID))
}

// ProductIsParent returns true if the product has variants (has variant dimensions defined)
func (store *Store) ProductIsParent(ctx context.Context, productID string) (bool, error) {
	if productID == "" {
		return false, errors.New("product ID is empty")
	}

	product, err := store.ProductFindByID(ctx, productID)
	if err != nil {
		return false, err
	}

	if product == nil {
		return false, nil
	}

	return product.HasVariantMatrixSchema(), nil
}

// ProductGetParent returns the parent product of a variant
func (store *Store) ProductGetParent(ctx context.Context, productID string) (ProductInterface, error) {
	if productID == "" {
		return nil, errors.New("product ID is empty")
	}

	product, err := store.ProductFindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, nil
	}

	if !product.IsVariant() {
		return nil, errors.New("product is not a variant")
	}

	return store.ProductFindByID(ctx, product.GetParentID())
}

// mapAnyToString converts a map[string]any to map[string]string
func mapAnyToString(m map[string]any) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = cast.ToString(v)
	}
	return result
}
