package shopstore

import (
	"context"
	"errors"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// ProductCategoryCount returns the number of product-category links
// matching the query options.
func (store *Store) ProductCategoryCount(ctx context.Context, options ProductCategoryQueryInterface) (int64, error) {
	q, err := store.productCategoryQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return -1, err
	}

	return count, nil
}

// ProductCategoryCreate inserts a new product-category link.
func (store *Store) ProductCategoryCreate(ctx context.Context, productCategory ProductCategoryInterface) error {
	if productCategory == nil {
		return errors.New("product category is nil")
	}

	if productCategory.GetProductID() == "" {
		return errors.New("product id is empty")
	}

	if productCategory.GetCategoryID() == "" {
		return errors.New("category id is empty")
	}

	productCategory.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	productCategory.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	productCategory.SetSoftDeletedAt(MAX_DATETIME)

	data := productCategory.Data()
	row := map[string]any{}
	for k, v := range data {
		row[k] = v
	}

	err := store.db.Query().Table(store.productCategoryTableName).Create(row)
	if err != nil {
		return err
	}

	productCategory.MarkAsNotDirty()

	return nil
}

// ProductCategoryDelete permanently deletes a product-category link.
func (store *Store) ProductCategoryDelete(ctx context.Context, productCategory ProductCategoryInterface) error {
	if productCategory == nil {
		return errors.New("product category is nil")
	}

	return store.ProductCategoryDeleteByID(ctx, productCategory.GetID())
}

// ProductCategoryDeleteByID permanently deletes a product-category link by ID.
func (store *Store) ProductCategoryDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product category id is empty")
	}

	_, err := store.db.Query().Table(store.productCategoryTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}

// ProductCategoryDeleteByProductID permanently deletes all category links
// for the given product.
func (store *Store) ProductCategoryDeleteByProductID(ctx context.Context, productID string) error {
	if productID == "" {
		return errors.New("product id is empty")
	}

	_, err := store.db.Query().Table(store.productCategoryTableName).Where(COLUMN_PRODUCT_ID+" = ?", productID).Delete()
	return err
}

// ProductCategoryFindByID returns the product-category link by ID, or nil.
func (store *Store) ProductCategoryFindByID(ctx context.Context, id string) (ProductCategoryInterface, error) {
	if id == "" {
		return nil, errors.New("product category id is empty")
	}

	list, err := store.ProductCategoryList(ctx, NewProductCategoryQuery().
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

// ProductCategoryIDs returns the IDs of all categories linked to a product.
func (store *Store) ProductCategoryIDs(ctx context.Context, productID string) ([]string, error) {
	if productID == "" {
		return []string{}, errors.New("product id is empty")
	}

	list, err := store.ProductCategoryList(ctx, NewProductCategoryQuery().
		SetProductID(productID))

	if err != nil {
		return []string{}, err
	}

	ids := []string{}
	for _, link := range list {
		ids = append(ids, link.GetCategoryID())
	}

	return ids, nil
}

// ProductCategorySet replaces all category links for a product with the
// given category IDs. Existing links not in the list are deleted; new
// links are created. This is the primary way to update a product's
// category assignments.
func (store *Store) ProductCategorySet(ctx context.Context, productID string, categoryIDs []string) error {
	if productID == "" {
		return errors.New("product id is empty")
	}

	existing, err := store.ProductCategoryList(ctx, NewProductCategoryQuery().
		SetProductID(productID))
	if err != nil {
		return err
	}

	wanted := map[string]bool{}
	for _, categoryID := range categoryIDs {
		wanted[categoryID] = true
	}

	for _, link := range existing {
		if !wanted[link.GetCategoryID()] {
			if err := store.ProductCategoryDeleteByID(ctx, link.GetID()); err != nil {
				return err
			}
		}
		delete(wanted, link.GetCategoryID())
	}

	for categoryID := range wanted {
		link := NewProductCategory().
			SetProductID(productID).
			SetCategoryID(categoryID)
		if err := store.ProductCategoryCreate(ctx, link); err != nil {
			return err
		}
	}

	return nil
}

// ProductCategoryList returns the product-category links matching the
// query options.
func (store *Store) ProductCategoryList(ctx context.Context, options ProductCategoryQueryInterface) ([]ProductCategoryInterface, error) {
	q, err := store.productCategoryQuery(options)
	if err != nil {
		return []ProductCategoryInterface{}, err
	}

	var results []map[string]any
	if err := q.Get(&results); err != nil {
		return []ProductCategoryInterface{}, err
	}

	list := []ProductCategoryInterface{}

	lo.ForEach(results, func(result map[string]any, index int) {
		model := NewProductCategoryFromExistingData(mapAnyToString(result))
		list = append(list, model)
	})

	return list, nil
}

// ProductCategorySoftDelete marks a product-category link as soft deleted.
func (store *Store) ProductCategorySoftDelete(ctx context.Context, productCategory ProductCategoryInterface) error {
	if productCategory == nil {
		return errors.New("product category is nil")
	}

	productCategory.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.ProductCategoryUpdate(ctx, productCategory)
}

// ProductCategorySoftDeleteByID marks a product-category link as soft
// deleted by ID.
func (store *Store) ProductCategorySoftDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product category id is empty")
	}

	productCategory, err := store.ProductCategoryFindByID(ctx, id)
	if err != nil {
		return err
	}
	if productCategory == nil {
		return nil
	}

	return store.ProductCategorySoftDelete(ctx, productCategory)
}

// ProductCategoryUpdate updates an existing product-category link.
func (store *Store) ProductCategoryUpdate(ctx context.Context, productCategory ProductCategoryInterface) error {
	if productCategory == nil {
		return errors.New("product category is nil")
	}

	productCategory.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := productCategory.DataChanged()

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

	_, err := store.db.Query().Table(store.productCategoryTableName).Where(COLUMN_ID+" = ?", productCategory.GetID()).Update(row)

	productCategory.MarkAsNotDirty()

	return err
}

func (store *Store) productCategoryQuery(options ProductCategoryQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		options = NewProductCategoryQuery()
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.productCategoryTableName)

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

	if options.HasProductID() {
		q = q.Where(COLUMN_PRODUCT_ID+" = ?", options.ProductID())
	}

	if options.HasCategoryID() {
		q = q.Where(COLUMN_CATEGORY_ID+" = ?", options.CategoryID())
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
