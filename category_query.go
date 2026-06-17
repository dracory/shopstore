package shopstore

import "errors"

type CategoryQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) CategoryQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) CategoryQueryInterface

	HasID() bool
	ID() string
	SetID(id string) CategoryQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) CategoryQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) CategoryQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) CategoryQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) CategoryQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) CategoryQueryInterface

	HasParentID() bool
	ParentID() string
	SetParentID(parentID string) CategoryQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) CategoryQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) CategoryQueryInterface

	HasTitleLike() bool
	TitleLike() string
	SetTitleLike(titleLike string) CategoryQueryInterface

	hasProperty(name string) bool
}

func NewCategoryQuery() CategoryQueryInterface {
	return &categoryQueryImplementation{
		properties: make(map[string]any),
	}
}

type categoryQueryImplementation struct {
	properties map[string]any
}

func (c *categoryQueryImplementation) Validate() error {
	if c.HasID() && c.ID() == "" {
		return errors.New("category query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("category query. id_in cannot be empty")
	}

	if c.HasParentID() && c.ParentID() == "" {
		return errors.New("category query. parent_id cannot be empty")
	}

	if c.HasStatus() && c.Status() == "" {
		return errors.New("category query. status cannot be empty")
	}

	if c.HasTitleLike() && c.TitleLike() == "" {
		return errors.New("category query. title_like cannot be empty")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("category query. order_by cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("category query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("category query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("category query. offset must be greater than or equal to 0")
	}

	return nil
}

func (c *categoryQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *categoryQueryImplementation) SetColumns(columns []string) CategoryQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *categoryQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *categoryQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *categoryQueryImplementation) SetCountOnly(countOnly bool) CategoryQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *categoryQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *categoryQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *categoryQueryImplementation) SetID(id string) CategoryQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *categoryQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *categoryQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *categoryQueryImplementation) SetIDIn(idIn []string) CategoryQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *categoryQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *categoryQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *categoryQueryImplementation) SetLimit(limit int) CategoryQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *categoryQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *categoryQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *categoryQueryImplementation) SetOffset(offset int) CategoryQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *categoryQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *categoryQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *categoryQueryImplementation) SetOrderBy(orderBy string) CategoryQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *categoryQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *categoryQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *categoryQueryImplementation) SetSortDirection(sortDirection string) CategoryQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *categoryQueryImplementation) HasParentID() bool {
	return c.hasProperty("parent_id")
}

func (c *categoryQueryImplementation) ParentID() string {
	if !c.HasParentID() {
		return ""
	}

	return c.properties["parent_id"].(string)
}

func (c *categoryQueryImplementation) SetParentID(parentID string) CategoryQueryInterface {
	c.properties["parent_id"] = parentID

	return c
}

func (c *categoryQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *categoryQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *categoryQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) CategoryQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *categoryQueryImplementation) HasStatus() bool {
	return c.hasProperty("status")
}

func (c *categoryQueryImplementation) Status() string {
	if !c.HasStatus() {
		return ""
	}

	return c.properties["status"].(string)
}

func (c *categoryQueryImplementation) SetStatus(status string) CategoryQueryInterface {
	c.properties["status"] = status

	return c
}

func (c *categoryQueryImplementation) HasTitleLike() bool {
	return c.hasProperty("title_like")
}

func (c *categoryQueryImplementation) TitleLike() string {
	if !c.HasTitleLike() {
		return ""
	}

	return c.properties["title_like"].(string)
}

func (c *categoryQueryImplementation) SetTitleLike(titleLike string) CategoryQueryInterface {
	c.properties["title_like"] = titleLike

	return c
}

func (c *categoryQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
