package shopstore

import "errors"

type SettingQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) SettingQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) SettingQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) SettingQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) SettingQueryInterface

	HasID() bool
	ID() string
	SetID(id string) SettingQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) SettingQueryInterface

	HasKey() bool
	Key() string
	SetKey(key string) SettingQueryInterface

	HasKeyLike() bool
	KeyLike() string
	SetKeyLike(keyLike string) SettingQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) SettingQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) SettingQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) SettingQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) SettingQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) SettingQueryInterface

	hasProperty(name string) bool
}

func NewSettingQuery() SettingQueryInterface {
	return &settingQueryImplementation{
		properties: make(map[string]any),
	}
}

type settingQueryImplementation struct {
	properties map[string]any
}

func (c *settingQueryImplementation) Validate() error {
	if c.HasCreatedAtGte() && c.CreatedAtGte() == "" {
		return errors.New("setting query. created_at_gte cannot be empty")
	}

	if c.HasCreatedAtLte() && c.CreatedAtLte() == "" {
		return errors.New("setting query. created_at_lte cannot be empty")
	}

	if c.HasID() && c.ID() == "" {
		return errors.New("setting query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("setting query. id_in cannot be empty")
	}

	if c.HasKey() && c.Key() == "" {
		return errors.New("setting query. key cannot be empty")
	}

	if c.HasKeyLike() && c.KeyLike() == "" {
		return errors.New("setting query. key_like cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("setting query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("setting query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("setting query. offset must be greater than or equal to 0")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("setting query. order_by cannot be empty")
	}

	return nil
}

func (c *settingQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *settingQueryImplementation) SetColumns(columns []string) SettingQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *settingQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *settingQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *settingQueryImplementation) SetCountOnly(countOnly bool) SettingQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *settingQueryImplementation) HasCreatedAtGte() bool {
	return c.hasProperty("created_at_gte")
}

func (c *settingQueryImplementation) CreatedAtGte() string {
	if !c.HasCreatedAtGte() {
		return ""
	}

	return c.properties["created_at_gte"].(string)
}

func (c *settingQueryImplementation) SetCreatedAtGte(createdAtGte string) SettingQueryInterface {
	c.properties["created_at_gte"] = createdAtGte

	return c
}

func (c *settingQueryImplementation) HasCreatedAtLte() bool {
	return c.hasProperty("created_at_lte")
}

func (c *settingQueryImplementation) CreatedAtLte() string {
	if !c.HasCreatedAtLte() {
		return ""
	}

	return c.properties["created_at_lte"].(string)
}

func (c *settingQueryImplementation) SetCreatedAtLte(createdAtLte string) SettingQueryInterface {
	c.properties["created_at_lte"] = createdAtLte

	return c
}

func (c *settingQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *settingQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *settingQueryImplementation) SetID(id string) SettingQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *settingQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *settingQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *settingQueryImplementation) SetIDIn(idIn []string) SettingQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *settingQueryImplementation) HasKey() bool {
	return c.hasProperty("key")
}

func (c *settingQueryImplementation) Key() string {
	if !c.HasKey() {
		return ""
	}

	return c.properties["key"].(string)
}

func (c *settingQueryImplementation) SetKey(key string) SettingQueryInterface {
	c.properties["key"] = key

	return c
}

func (c *settingQueryImplementation) HasKeyLike() bool {
	return c.hasProperty("key_like")
}

func (c *settingQueryImplementation) KeyLike() string {
	if !c.HasKeyLike() {
		return ""
	}

	return c.properties["key_like"].(string)
}

func (c *settingQueryImplementation) SetKeyLike(keyLike string) SettingQueryInterface {
	c.properties["key_like"] = keyLike

	return c
}

func (c *settingQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *settingQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *settingQueryImplementation) SetLimit(limit int) SettingQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *settingQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *settingQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *settingQueryImplementation) SetOffset(offset int) SettingQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *settingQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *settingQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *settingQueryImplementation) SetOrderBy(orderBy string) SettingQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *settingQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *settingQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *settingQueryImplementation) SetSortDirection(sortDirection string) SettingQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *settingQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *settingQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *settingQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) SettingQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *settingQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
