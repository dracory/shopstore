package shopstore

import "errors"

type DiscountQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) DiscountQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) DiscountQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) DiscountQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) DiscountQueryInterface

	HasCode() bool
	Code() string
	SetCode(code string) DiscountQueryInterface

	HasEndsAtGte() bool
	EndsAtGte() string
	SetEndsAtGte(endsAtGte string) DiscountQueryInterface

	HasEndsAtLte() bool
	EndsAtLte() string
	SetEndsAtLte(endsAtLte string) DiscountQueryInterface

	HasID() bool
	ID() string
	SetID(id string) DiscountQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) DiscountQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) DiscountQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) DiscountQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(discountBy string) DiscountQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) DiscountQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) DiscountQueryInterface

	HasStartsAtGte() bool
	StartsAtGte() string
	SetStartsAtGte(startsAtGte string) DiscountQueryInterface

	HasStartsAtLte() bool
	StartsAtLte() string
	SetStartsAtLte(startsAtLte string) DiscountQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) DiscountQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) DiscountQueryInterface

	HasType() bool
	Type() string
	SetType(discountType string) DiscountQueryInterface

	hasProperty(name string) bool
}

func NewDiscountQuery() DiscountQueryInterface {
	return &discountQueryImplementation{
		properties: make(map[string]any),
	}
}

type discountQueryImplementation struct {
	properties map[string]any
}

func (c *discountQueryImplementation) Validate() error {

	if c.HasCreatedAtGte() && c.CreatedAtGte() == "" {
		return errors.New("discount query. created_at_gte cannot be empty")
	}

	if c.HasCreatedAtLte() && c.CreatedAtLte() == "" {
		return errors.New("discount query. created_at_lte cannot be empty")
	}

	if c.HasCode() && c.Code() == "" {
		return errors.New("discount query. code cannot be empty")
	}

	if c.HasID() && c.ID() == "" {
		return errors.New("discount query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("discount query. id_in cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("discount query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("discount query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("discount query. offset must be greater than or equal to 0")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("discount query. order_by cannot be empty")
	}

	if c.HasStatus() && c.Status() == "" {
		return errors.New("discount query. status cannot be empty")
	}

	if c.HasStatusIn() && len(c.StatusIn()) == 0 {
		return errors.New("discount query. status_in cannot be empty")
	}

	if c.HasEndsAtGte() && c.EndsAtGte() == "" {
		return errors.New("discount query. ends_at_gte cannot be empty")
	}

	if c.HasEndsAtLte() && c.EndsAtLte() == "" {
		return errors.New("discount query. ends_at_lte cannot be empty")
	}

	if c.HasStartsAtGte() && c.StartsAtGte() == "" {
		return errors.New("discount query. starts_at_gte cannot be empty")
	}

	if c.HasStartsAtLte() && c.StartsAtLte() == "" {
		return errors.New("discount query. starts_at_lte cannot be empty")
	}

	if c.HasType() && c.Type() == "" {
		return errors.New("discount query. type cannot be empty")
	}

	return nil
}

func (c *discountQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *discountQueryImplementation) SetColumns(columns []string) DiscountQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *discountQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *discountQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *discountQueryImplementation) SetCountOnly(countOnly bool) DiscountQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *discountQueryImplementation) HasCode() bool {
	return c.hasProperty("code")
}

func (c *discountQueryImplementation) Code() string {
	if !c.HasCode() {
		return ""
	}

	return c.properties["code"].(string)
}

func (c *discountQueryImplementation) SetCode(code string) DiscountQueryInterface {
	c.properties["code"] = code

	return c
}

func (c *discountQueryImplementation) HasCreatedAtGte() bool {
	return c.hasProperty("created_at_gte")
}

func (c *discountQueryImplementation) CreatedAtGte() string {
	if !c.HasCreatedAtGte() {
		return ""
	}

	return c.properties["created_at_gte"].(string)
}

func (c *discountQueryImplementation) SetCreatedAtGte(createdAtGte string) DiscountQueryInterface {
	c.properties["created_at_gte"] = createdAtGte

	return c
}

func (c *discountQueryImplementation) HasCreatedAtLte() bool {
	return c.hasProperty("created_at_lte")
}

func (c *discountQueryImplementation) CreatedAtLte() string {
	if !c.HasCreatedAtLte() {
		return ""
	}

	return c.properties["created_at_lte"].(string)
}

func (c *discountQueryImplementation) SetCreatedAtLte(createdAtLte string) DiscountQueryInterface {
	c.properties["created_at_lte"] = createdAtLte

	return c
}

func (c *discountQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *discountQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *discountQueryImplementation) SetID(id string) DiscountQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *discountQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *discountQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *discountQueryImplementation) SetIDIn(idIn []string) DiscountQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *discountQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *discountQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *discountQueryImplementation) SetLimit(limit int) DiscountQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *discountQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *discountQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *discountQueryImplementation) SetOffset(offset int) DiscountQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *discountQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *discountQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *discountQueryImplementation) SetOrderBy(orderBy string) DiscountQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *discountQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *discountQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *discountQueryImplementation) SetSortDirection(sortDirection string) DiscountQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *discountQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *discountQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *discountQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) DiscountQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *discountQueryImplementation) HasStatus() bool {
	return c.hasProperty("status")
}

func (c *discountQueryImplementation) Status() string {
	if !c.HasStatus() {
		return ""
	}

	return c.properties["status"].(string)
}

func (c *discountQueryImplementation) SetStatus(status string) DiscountQueryInterface {
	c.properties["status"] = status

	return c
}

func (c *discountQueryImplementation) HasStatusIn() bool {
	return c.hasProperty("status_in")
}

func (c *discountQueryImplementation) StatusIn() []string {
	if !c.HasStatusIn() {
		return []string{}
	}

	return c.properties["status_in"].([]string)
}

func (c *discountQueryImplementation) SetStatusIn(statusIn []string) DiscountQueryInterface {
	c.properties["status_in"] = statusIn

	return c
}

func (c *discountQueryImplementation) HasEndsAtGte() bool {
	return c.hasProperty("ends_at_gte")
}

func (c *discountQueryImplementation) EndsAtGte() string {
	if !c.HasEndsAtGte() {
		return ""
	}

	return c.properties["ends_at_gte"].(string)
}

func (c *discountQueryImplementation) SetEndsAtGte(endsAtGte string) DiscountQueryInterface {
	c.properties["ends_at_gte"] = endsAtGte

	return c
}

func (c *discountQueryImplementation) HasEndsAtLte() bool {
	return c.hasProperty("ends_at_lte")
}

func (c *discountQueryImplementation) EndsAtLte() string {
	if !c.HasEndsAtLte() {
		return ""
	}

	return c.properties["ends_at_lte"].(string)
}

func (c *discountQueryImplementation) SetEndsAtLte(endsAtLte string) DiscountQueryInterface {
	c.properties["ends_at_lte"] = endsAtLte

	return c
}

func (c *discountQueryImplementation) HasStartsAtGte() bool {
	return c.hasProperty("starts_at_gte")
}

func (c *discountQueryImplementation) StartsAtGte() string {
	if !c.HasStartsAtGte() {
		return ""
	}

	return c.properties["starts_at_gte"].(string)
}

func (c *discountQueryImplementation) SetStartsAtGte(startsAtGte string) DiscountQueryInterface {
	c.properties["starts_at_gte"] = startsAtGte

	return c
}

func (c *discountQueryImplementation) HasStartsAtLte() bool {
	return c.hasProperty("starts_at_lte")
}

func (c *discountQueryImplementation) StartsAtLte() string {
	if !c.HasStartsAtLte() {
		return ""
	}

	return c.properties["starts_at_lte"].(string)
}

func (c *discountQueryImplementation) SetStartsAtLte(startsAtLte string) DiscountQueryInterface {
	c.properties["starts_at_lte"] = startsAtLte

	return c
}

func (c *discountQueryImplementation) HasType() bool {
	return c.hasProperty("type")
}

func (c *discountQueryImplementation) Type() string {
	if !c.HasType() {
		return ""
	}

	return c.properties["type"].(string)
}

func (c *discountQueryImplementation) SetType(discountType string) DiscountQueryInterface {
	c.properties["type"] = discountType

	return c
}

func (c *discountQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
