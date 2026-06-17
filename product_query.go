package shopstore

import "errors"

type ProductQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) ProductQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) ProductQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) ProductQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) ProductQueryInterface

	HasID() bool
	ID() string
	SetID(id string) ProductQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) ProductQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) ProductQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) ProductQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) ProductQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) ProductQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) ProductQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) ProductQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) ProductQueryInterface

	HasTitleLike() bool
	TitleLike() string
	SetTitleLike(titleLike string) ProductQueryInterface

	HasParentID() bool
	ParentID() string
	SetParentID(parentID string) ProductQueryInterface

	hasProperty(name string) bool
}

func NewProductQuery() ProductQueryInterface {
	return &productQueryImplementation{
		properties: make(map[string]any),
	}
}

type productQueryImplementation struct {
	properties map[string]any
}

func (c *productQueryImplementation) Validate() error {

	if c.HasCreatedAtGte() && c.CreatedAtGte() == "" {
		return errors.New("product query. created_at_gte cannot be empty")
	}

	if c.HasCreatedAtLte() && c.CreatedAtLte() == "" {
		return errors.New("product query. created_at_lte cannot be empty")
	}

	if c.HasID() && c.ID() == "" {
		return errors.New("product query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("product query. id_in cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("product query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("product query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("product query. offset must be greater than or equal to 0")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("product query. order_by cannot be empty")
	}

	if c.HasStatus() && c.Status() == "" {
		return errors.New("product query. status cannot be empty")
	}

	if c.HasStatusIn() && len(c.StatusIn()) == 0 {
		return errors.New("product query. status_in cannot be empty")
	}

	if c.HasTitleLike() && c.TitleLike() == "" {
		return errors.New("product query. title_like cannot be empty")
	}

	return nil
}

func (c *productQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *productQueryImplementation) SetColumns(columns []string) ProductQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *productQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *productQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *productQueryImplementation) SetCountOnly(countOnly bool) ProductQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *productQueryImplementation) HasCreatedAtGte() bool {
	return c.hasProperty("created_at_gte")
}

func (c *productQueryImplementation) CreatedAtGte() string {
	if !c.HasCreatedAtGte() {
		return ""
	}

	return c.properties["created_at_gte"].(string)
}

func (c *productQueryImplementation) SetCreatedAtGte(createdAtGte string) ProductQueryInterface {
	c.properties["created_at_gte"] = createdAtGte

	return c
}

func (c *productQueryImplementation) HasCreatedAtLte() bool {
	return c.hasProperty("created_at_lte")
}

func (c *productQueryImplementation) CreatedAtLte() string {
	if !c.HasCreatedAtLte() {
		return ""
	}

	return c.properties["created_at_lte"].(string)
}

func (c *productQueryImplementation) SetCreatedAtLte(createdAtLte string) ProductQueryInterface {
	c.properties["created_at_lte"] = createdAtLte

	return c
}

func (c *productQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *productQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *productQueryImplementation) SetID(id string) ProductQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *productQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *productQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *productQueryImplementation) SetIDIn(idIn []string) ProductQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *productQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *productQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *productQueryImplementation) SetLimit(limit int) ProductQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *productQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *productQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *productQueryImplementation) SetOffset(offset int) ProductQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *productQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *productQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *productQueryImplementation) SetOrderBy(orderBy string) ProductQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *productQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *productQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *productQueryImplementation) SetSortDirection(sortDirection string) ProductQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *productQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *productQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *productQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) ProductQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *productQueryImplementation) HasStatus() bool {
	return c.hasProperty("status")
}

func (c *productQueryImplementation) Status() string {
	if !c.HasStatus() {
		return ""
	}

	return c.properties["status"].(string)
}

func (c *productQueryImplementation) SetStatus(status string) ProductQueryInterface {
	c.properties["status"] = status

	return c
}

func (c *productQueryImplementation) HasStatusIn() bool {
	return c.hasProperty("status_in")
}

func (c *productQueryImplementation) StatusIn() []string {
	if !c.HasStatusIn() {
		return []string{}
	}

	return c.properties["status_in"].([]string)
}

func (c *productQueryImplementation) SetStatusIn(statusIn []string) ProductQueryInterface {
	c.properties["status_in"] = statusIn

	return c
}

func (c *productQueryImplementation) HasTitleLike() bool {
	return c.hasProperty("title_like")
}

func (c *productQueryImplementation) TitleLike() string {
	if !c.HasTitleLike() {
		return ""
	}

	return c.properties["title_like"].(string)
}

func (c *productQueryImplementation) SetTitleLike(titleLike string) ProductQueryInterface {
	c.properties["title_like"] = titleLike

	return c
}

func (c *productQueryImplementation) HasParentID() bool {
	return c.hasProperty("parent_id")
}

func (c *productQueryImplementation) ParentID() string {
	if !c.HasParentID() {
		return ""
	}

	return c.properties["parent_id"].(string)
}

func (c *productQueryImplementation) SetParentID(parentID string) ProductQueryInterface {
	c.properties["parent_id"] = parentID

	return c
}

func (c *productQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
