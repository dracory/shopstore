package shopstore

import "errors"

// ProductCategoryQueryInterface defines the query contract for
// filtering product-category links.
type ProductCategoryQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) ProductCategoryQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) ProductCategoryQueryInterface

	HasCategoryID() bool
	CategoryID() string
	SetCategoryID(categoryID string) ProductCategoryQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) ProductCategoryQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) ProductCategoryQueryInterface

	HasID() bool
	ID() string
	SetID(id string) ProductCategoryQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) ProductCategoryQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) ProductCategoryQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) ProductCategoryQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) ProductCategoryQueryInterface

	HasProductID() bool
	ProductID() string
	SetProductID(productID string) ProductCategoryQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) ProductCategoryQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) ProductCategoryQueryInterface

	hasProperty(name string) bool
}

// NewProductCategoryQuery creates a new product-category query.
func NewProductCategoryQuery() ProductCategoryQueryInterface {
	return &productCategoryQueryImplementation{
		properties: make(map[string]any),
	}
}

type productCategoryQueryImplementation struct {
	properties map[string]any
}

func (c *productCategoryQueryImplementation) Validate() error {
	if c.HasCategoryID() && c.CategoryID() == "" {
		return errors.New("product category query. category_id cannot be empty")
	}

	if c.HasCreatedAtGte() && c.CreatedAtGte() == "" {
		return errors.New("product category query. created_at_gte cannot be empty")
	}

	if c.HasCreatedAtLte() && c.CreatedAtLte() == "" {
		return errors.New("product category query. created_at_lte cannot be empty")
	}

	if c.HasID() && c.ID() == "" {
		return errors.New("product category query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("product category query. id_in cannot be empty")
	}

	if c.HasProductID() && c.ProductID() == "" {
		return errors.New("product category query. product_id cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("product category query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("product category query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("product category query. offset must be greater than or equal to 0")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("product category query. order_by cannot be empty")
	}

	return nil
}

func (c *productCategoryQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *productCategoryQueryImplementation) SetColumns(columns []string) ProductCategoryQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *productCategoryQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *productCategoryQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *productCategoryQueryImplementation) SetCountOnly(countOnly bool) ProductCategoryQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *productCategoryQueryImplementation) HasCategoryID() bool {
	return c.hasProperty("category_id")
}

func (c *productCategoryQueryImplementation) CategoryID() string {
	if !c.HasCategoryID() {
		return ""
	}

	return c.properties["category_id"].(string)
}

func (c *productCategoryQueryImplementation) SetCategoryID(categoryID string) ProductCategoryQueryInterface {
	c.properties["category_id"] = categoryID

	return c
}

func (c *productCategoryQueryImplementation) HasCreatedAtGte() bool {
	return c.hasProperty("created_at_gte")
}

func (c *productCategoryQueryImplementation) CreatedAtGte() string {
	if !c.HasCreatedAtGte() {
		return ""
	}

	return c.properties["created_at_gte"].(string)
}

func (c *productCategoryQueryImplementation) SetCreatedAtGte(createdAtGte string) ProductCategoryQueryInterface {
	c.properties["created_at_gte"] = createdAtGte

	return c
}

func (c *productCategoryQueryImplementation) HasCreatedAtLte() bool {
	return c.hasProperty("created_at_lte")
}

func (c *productCategoryQueryImplementation) CreatedAtLte() string {
	if !c.HasCreatedAtLte() {
		return ""
	}

	return c.properties["created_at_lte"].(string)
}

func (c *productCategoryQueryImplementation) SetCreatedAtLte(createdAtLte string) ProductCategoryQueryInterface {
	c.properties["created_at_lte"] = createdAtLte

	return c
}

func (c *productCategoryQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *productCategoryQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *productCategoryQueryImplementation) SetID(id string) ProductCategoryQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *productCategoryQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *productCategoryQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *productCategoryQueryImplementation) SetIDIn(idIn []string) ProductCategoryQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *productCategoryQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *productCategoryQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *productCategoryQueryImplementation) SetLimit(limit int) ProductCategoryQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *productCategoryQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *productCategoryQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *productCategoryQueryImplementation) SetOffset(offset int) ProductCategoryQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *productCategoryQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *productCategoryQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *productCategoryQueryImplementation) SetOrderBy(orderBy string) ProductCategoryQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *productCategoryQueryImplementation) HasProductID() bool {
	return c.hasProperty("product_id")
}

func (c *productCategoryQueryImplementation) ProductID() string {
	if !c.HasProductID() {
		return ""
	}

	return c.properties["product_id"].(string)
}

func (c *productCategoryQueryImplementation) SetProductID(productID string) ProductCategoryQueryInterface {
	c.properties["product_id"] = productID

	return c
}

func (c *productCategoryQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *productCategoryQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *productCategoryQueryImplementation) SetSortDirection(sortDirection string) ProductCategoryQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *productCategoryQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *productCategoryQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *productCategoryQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) ProductCategoryQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *productCategoryQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
