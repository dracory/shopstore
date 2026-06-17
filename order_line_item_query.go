package shopstore

import "errors"

type OrderLineItemQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) OrderLineItemQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) OrderLineItemQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) OrderLineItemQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) OrderLineItemQueryInterface

	HasID() bool
	ID() string
	SetID(id string) OrderLineItemQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) OrderLineItemQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) OrderLineItemQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) OrderLineItemQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) OrderLineItemQueryInterface

	HasOrderID() bool
	OrderID() string
	SetOrderID(orderID string) OrderLineItemQueryInterface

	HasOrderIDIn() bool
	OrderIDIn() []string
	SetOrderIDIn(orderIDIn []string) OrderLineItemQueryInterface

	HasProductID() bool
	ProductID() string
	SetProductID(productID string) OrderLineItemQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) OrderLineItemQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) OrderLineItemQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) OrderLineItemQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) OrderLineItemQueryInterface

	hasProperty(name string) bool
}

func NewOrderLineItemQuery() OrderLineItemQueryInterface {
	return &orderLineItemQueryImplementation{
		properties: make(map[string]any),
	}
}

type orderLineItemQueryImplementation struct {
	properties map[string]any
}

func (c *orderLineItemQueryImplementation) Validate() error {

	if c.HasCreatedAtGte() && c.CreatedAtGte() == "" {
		return errors.New("orderLineItem query. created_at_gte cannot be empty")
	}

	if c.HasCreatedAtLte() && c.CreatedAtLte() == "" {
		return errors.New("orderLineItem query. created_at_lte cannot be empty")
	}

	if c.HasCustomerID() && c.CustomerID() == "" {
		return errors.New("orderLineItem query. customer_id cannot be empty")
	}

	if c.HasID() && c.ID() == "" {
		return errors.New("orderLineItem query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("orderLineItem query. id_in cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("orderLineItem query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("orderLineItem query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("orderLineItem query. offset must be greater than or equal to 0")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("orderLineItem query. order_by cannot be empty")
	}

	if c.HasOrderID() && c.OrderID() == "" {
		return errors.New("orderLineItem query. order_id cannot be empty")
	}

	if c.HasOrderIDIn() {
		if len(c.OrderIDIn()) < 1 {
			return errors.New("orderLineItem query. order_id_in cannot be empty")
		}
	}

	if c.HasProductID() && c.ProductID() == "" {
		return errors.New("orderLineItem query. product_id cannot be empty")
	}

	if c.HasStatus() && c.Status() == "" {
		return errors.New("orderLineItem query. status cannot be empty")
	}

	return nil
}

func (c *orderLineItemQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *orderLineItemQueryImplementation) SetColumns(columns []string) OrderLineItemQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *orderLineItemQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *orderLineItemQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *orderLineItemQueryImplementation) SetCountOnly(countOnly bool) OrderLineItemQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *orderLineItemQueryImplementation) HasCreatedAtGte() bool {
	return c.hasProperty("created_at_gte")
}

func (c *orderLineItemQueryImplementation) CreatedAtGte() string {
	if !c.HasCreatedAtGte() {
		return ""
	}

	return c.properties["created_at_gte"].(string)
}

func (c *orderLineItemQueryImplementation) SetCreatedAtGte(createdAtGte string) OrderLineItemQueryInterface {
	c.properties["created_at_gte"] = createdAtGte

	return c
}

func (c *orderLineItemQueryImplementation) HasCreatedAtLte() bool {
	return c.hasProperty("created_at_lte")
}

func (c *orderLineItemQueryImplementation) CreatedAtLte() string {
	if !c.HasCreatedAtLte() {
		return ""
	}

	return c.properties["created_at_lte"].(string)
}

func (c *orderLineItemQueryImplementation) SetCreatedAtLte(createdAtLte string) OrderLineItemQueryInterface {
	c.properties["created_at_lte"] = createdAtLte

	return c
}

func (c *orderLineItemQueryImplementation) HasCustomerID() bool {
	return c.hasProperty("customer_id")
}

func (c *orderLineItemQueryImplementation) CustomerID() string {
	if !c.HasCustomerID() {
		return ""
	}

	return c.properties["customer_id"].(string)
}

func (c *orderLineItemQueryImplementation) SetCustomerID(customerID string) OrderLineItemQueryInterface {
	c.properties["customer_id"] = customerID

	return c
}

func (c *orderLineItemQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *orderLineItemQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *orderLineItemQueryImplementation) SetID(id string) OrderLineItemQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *orderLineItemQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *orderLineItemQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *orderLineItemQueryImplementation) SetIDIn(idIn []string) OrderLineItemQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *orderLineItemQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *orderLineItemQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *orderLineItemQueryImplementation) SetLimit(limit int) OrderLineItemQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *orderLineItemQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *orderLineItemQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *orderLineItemQueryImplementation) SetOffset(offset int) OrderLineItemQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *orderLineItemQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *orderLineItemQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *orderLineItemQueryImplementation) SetOrderBy(orderBy string) OrderLineItemQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *orderLineItemQueryImplementation) HasOrderID() bool {
	return c.hasProperty("order_id")
}

func (c *orderLineItemQueryImplementation) OrderID() string {
	if !c.HasOrderID() {
		return ""
	}

	return c.properties["order_id"].(string)
}

func (c *orderLineItemQueryImplementation) SetOrderID(orderID string) OrderLineItemQueryInterface {
	c.properties["order_id"] = orderID

	return c
}

func (c *orderLineItemQueryImplementation) HasOrderIDIn() bool {
	return c.hasProperty("order_id_in")
}

func (c *orderLineItemQueryImplementation) OrderIDIn() []string {
	if !c.HasOrderIDIn() {
		return []string{}
	}

	return c.properties["order_id_in"].([]string)
}

func (c *orderLineItemQueryImplementation) SetOrderIDIn(orderIDIn []string) OrderLineItemQueryInterface {
	c.properties["order_id_in"] = orderIDIn

	return c
}

func (c *orderLineItemQueryImplementation) HasProductID() bool {
	return c.hasProperty("product_id")
}

func (c *orderLineItemQueryImplementation) ProductID() string {
	if !c.HasProductID() {
		return ""
	}

	return c.properties["product_id"].(string)
}

func (c *orderLineItemQueryImplementation) SetProductID(productID string) OrderLineItemQueryInterface {
	c.properties["product_id"] = productID

	return c
}

func (c *orderLineItemQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *orderLineItemQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *orderLineItemQueryImplementation) SetSortDirection(sortDirection string) OrderLineItemQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *orderLineItemQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *orderLineItemQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *orderLineItemQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) OrderLineItemQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *orderLineItemQueryImplementation) HasStatus() bool {
	return c.hasProperty("status")
}

func (c *orderLineItemQueryImplementation) Status() string {
	if !c.HasStatus() {
		return ""
	}

	return c.properties["status"].(string)
}

func (c *orderLineItemQueryImplementation) SetStatus(status string) OrderLineItemQueryInterface {
	c.properties["status"] = status

	return c
}

func (c *orderLineItemQueryImplementation) HasStatusIn() bool {
	return c.hasProperty("status_in")
}

func (c *orderLineItemQueryImplementation) StatusIn() []string {
	if !c.HasStatusIn() {
		return []string{}
	}

	return c.properties["status_in"].([]string)
}

func (c *orderLineItemQueryImplementation) SetStatusIn(statusIn []string) OrderLineItemQueryInterface {
	c.properties["status_in"] = statusIn

	return c
}

func (c *orderLineItemQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
