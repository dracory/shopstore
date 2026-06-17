package shopstore

import "errors"

type OrderQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) OrderQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) OrderQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) OrderQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) OrderQueryInterface

	HasCustomerID() bool
	CustomerID() string
	SetCustomerID(customerID string) OrderQueryInterface

	HasID() bool
	ID() string
	SetID(id string) OrderQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) OrderQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) OrderQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) OrderQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) OrderQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) OrderQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) OrderQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) OrderQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) OrderQueryInterface

	hasProperty(name string) bool
}

func NewOrderQuery() OrderQueryInterface {
	return &orderQueryImplementation{
		properties: make(map[string]any),
	}
}

type orderQueryImplementation struct {
	properties map[string]any
}

func (c *orderQueryImplementation) Validate() error {

	if c.HasCreatedAtGte() && c.CreatedAtGte() == "" {
		return errors.New("order query. created_at_gte cannot be empty")
	}

	if c.HasCreatedAtLte() && c.CreatedAtLte() == "" {
		return errors.New("order query. created_at_lte cannot be empty")
	}

	if c.HasCustomerID() && c.CustomerID() == "" {
		return errors.New("order query. customer_id cannot be empty")
	}

	if c.HasID() && c.ID() == "" {
		return errors.New("order query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("order query. id_in cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("order query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("order query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("order query. offset must be greater than or equal to 0")
	}

	if c.HasStatus() && c.Status() == "" {
		return errors.New("order query. status cannot be empty")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("order query. order_by cannot be empty")
	}

	return nil
}

func (c *orderQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *orderQueryImplementation) SetColumns(columns []string) OrderQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *orderQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *orderQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *orderQueryImplementation) SetCountOnly(countOnly bool) OrderQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *orderQueryImplementation) HasCreatedAtGte() bool {
	return c.hasProperty("created_at_gte")
}

func (c *orderQueryImplementation) CreatedAtGte() string {
	if !c.HasCreatedAtGte() {
		return ""
	}

	return c.properties["created_at_gte"].(string)
}

func (c *orderQueryImplementation) SetCreatedAtGte(createdAtGte string) OrderQueryInterface {
	c.properties["created_at_gte"] = createdAtGte

	return c
}

func (c *orderQueryImplementation) HasCreatedAtLte() bool {
	return c.hasProperty("created_at_lte")
}

func (c *orderQueryImplementation) CreatedAtLte() string {
	if !c.HasCreatedAtLte() {
		return ""
	}

	return c.properties["created_at_lte"].(string)
}

func (c *orderQueryImplementation) SetCreatedAtLte(createdAtLte string) OrderQueryInterface {
	c.properties["created_at_lte"] = createdAtLte

	return c
}

func (c *orderQueryImplementation) HasCustomerID() bool {
	return c.hasProperty("customer_id")
}

func (c *orderQueryImplementation) CustomerID() string {
	if !c.HasCustomerID() {
		return ""
	}

	return c.properties["customer_id"].(string)
}

func (c *orderQueryImplementation) SetCustomerID(customerID string) OrderQueryInterface {
	c.properties["customer_id"] = customerID

	return c
}

func (c *orderQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *orderQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *orderQueryImplementation) SetID(id string) OrderQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *orderQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *orderQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *orderQueryImplementation) SetIDIn(idIn []string) OrderQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *orderQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *orderQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *orderQueryImplementation) SetLimit(limit int) OrderQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *orderQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *orderQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *orderQueryImplementation) SetOffset(offset int) OrderQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *orderQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *orderQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *orderQueryImplementation) SetOrderBy(orderBy string) OrderQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *orderQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *orderQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *orderQueryImplementation) SetSortDirection(sortDirection string) OrderQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *orderQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *orderQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *orderQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) OrderQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *orderQueryImplementation) HasStatus() bool {
	return c.hasProperty("status")
}

func (c *orderQueryImplementation) Status() string {
	if !c.HasStatus() {
		return ""
	}

	return c.properties["status"].(string)
}

func (c *orderQueryImplementation) SetStatus(status string) OrderQueryInterface {
	c.properties["status"] = status

	return c
}

func (c *orderQueryImplementation) HasStatusIn() bool {
	return c.hasProperty("status_in")
}

func (c *orderQueryImplementation) StatusIn() []string {
	if !c.HasStatusIn() {
		return []string{}
	}

	return c.properties["status_in"].([]string)
}

func (c *orderQueryImplementation) SetStatusIn(statusIn []string) OrderQueryInterface {
	c.properties["status_in"] = statusIn

	return c
}

func (c *orderQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
