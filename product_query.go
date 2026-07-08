package shopstore

import "errors"

const (
	propertyColumns             = "columns"
	propertyCountOnly           = "count_only"
	propertyCreatedAtGte        = "created_at_gte"
	propertyCreatedAtLte        = "created_at_lte"
	propertyID                  = "id"
	propertyIDIn                = "id_in"
	propertyIDNotIn             = "id_not_in"
	propertyLimit               = "limit"
	propertyOffset              = "offset"
	propertyOrderBy             = "order_by"
	propertySortDirection       = "sort_direction"
	propertySoftDeletedIncluded = "soft_deleted_included"
	propertyStatus              = "status"
	propertyStatusIn            = "status_in"
	propertyTitleLike           = "title_like"
	propertyParentID            = "parent_id"
	propertyMetasIn             = "metas_in"
	propertyMetasNotIn          = "metas_not_in"
)

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

	HasIDNotIn() bool
	IDNotIn() []string
	SetIDNotIn(idNotIn []string) ProductQueryInterface

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

	HasMetasIn() bool
	MetasIn() map[string]string
	SetMetasIn(metasIn map[string]string) ProductQueryInterface

	HasMetasNotIn() bool
	MetasNotIn() map[string]string
	SetMetasNotIn(metasNotIn map[string]string) ProductQueryInterface

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

	if c.HasIDNotIn() && len(c.IDNotIn()) == 0 {
		return errors.New("product query. id_not_in cannot be empty")
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

	if c.HasMetasIn() {
		if len(c.MetasIn()) == 0 {
			return errors.New("product query. metas_in cannot be empty")
		}
		for k, v := range c.MetasIn() {
			if k == "" || v == "" {
				return errors.New("product query. metas_in keys and values cannot be empty")
			}
		}
	}

	if c.HasMetasNotIn() {
		if len(c.MetasNotIn()) == 0 {
			return errors.New("product query. metas_not_in cannot be empty")
		}
		for k, v := range c.MetasNotIn() {
			if k == "" || v == "" {
				return errors.New("product query. metas_not_in keys and values cannot be empty")
			}
		}
	}

	return nil
}

func (c *productQueryImplementation) Columns() []string {
	if !c.hasProperty(propertyColumns) {
		return []string{}
	}

	return c.properties[propertyColumns].([]string)
}

func (c *productQueryImplementation) SetColumns(columns []string) ProductQueryInterface {
	c.properties[propertyColumns] = columns

	return c
}

func (c *productQueryImplementation) HasCountOnly() bool {
	return c.hasProperty(propertyCountOnly)
}

func (c *productQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties[propertyCountOnly].(bool)
}

func (c *productQueryImplementation) SetCountOnly(countOnly bool) ProductQueryInterface {
	c.properties[propertyCountOnly] = countOnly

	return c
}

func (c *productQueryImplementation) HasCreatedAtGte() bool {
	return c.hasProperty(propertyCreatedAtGte)
}

func (c *productQueryImplementation) CreatedAtGte() string {
	if !c.HasCreatedAtGte() {
		return ""
	}

	return c.properties[propertyCreatedAtGte].(string)
}

func (c *productQueryImplementation) SetCreatedAtGte(createdAtGte string) ProductQueryInterface {
	c.properties[propertyCreatedAtGte] = createdAtGte

	return c
}

func (c *productQueryImplementation) HasCreatedAtLte() bool {
	return c.hasProperty(propertyCreatedAtLte)
}

func (c *productQueryImplementation) CreatedAtLte() string {
	if !c.HasCreatedAtLte() {
		return ""
	}

	return c.properties[propertyCreatedAtLte].(string)
}

func (c *productQueryImplementation) SetCreatedAtLte(createdAtLte string) ProductQueryInterface {
	c.properties[propertyCreatedAtLte] = createdAtLte

	return c
}

func (c *productQueryImplementation) HasID() bool {
	return c.hasProperty(propertyID)
}

func (c *productQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties[propertyID].(string)
}

func (c *productQueryImplementation) SetID(id string) ProductQueryInterface {
	c.properties[propertyID] = id

	return c
}

func (c *productQueryImplementation) HasIDIn() bool {
	return c.hasProperty(propertyIDIn)
}

func (c *productQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties[propertyIDIn].([]string)
}

func (c *productQueryImplementation) SetIDIn(idIn []string) ProductQueryInterface {
	c.properties[propertyIDIn] = idIn

	return c
}

func (c *productQueryImplementation) HasIDNotIn() bool {
	return c.hasProperty(propertyIDNotIn)
}

func (c *productQueryImplementation) IDNotIn() []string {
	if !c.HasIDNotIn() {
		return []string{}
	}

	return c.properties[propertyIDNotIn].([]string)
}

func (c *productQueryImplementation) SetIDNotIn(idNotIn []string) ProductQueryInterface {
	c.properties[propertyIDNotIn] = idNotIn

	return c
}

func (c *productQueryImplementation) HasLimit() bool {
	return c.hasProperty(propertyLimit)
}

func (c *productQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties[propertyLimit].(int)
}

func (c *productQueryImplementation) SetLimit(limit int) ProductQueryInterface {
	c.properties[propertyLimit] = limit

	return c
}

func (c *productQueryImplementation) HasOffset() bool {
	return c.hasProperty(propertyOffset)
}

func (c *productQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties[propertyOffset].(int)
}

func (c *productQueryImplementation) SetOffset(offset int) ProductQueryInterface {
	c.properties[propertyOffset] = offset

	return c
}

func (c *productQueryImplementation) HasOrderBy() bool {
	return c.hasProperty(propertyOrderBy)
}

func (c *productQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties[propertyOrderBy].(string)
}

func (c *productQueryImplementation) SetOrderBy(orderBy string) ProductQueryInterface {
	c.properties[propertyOrderBy] = orderBy

	return c
}

func (c *productQueryImplementation) HasSortDirection() bool {
	return c.hasProperty(propertySortDirection)
}

func (c *productQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties[propertySortDirection].(string)
}

func (c *productQueryImplementation) SetSortDirection(sortDirection string) ProductQueryInterface {
	c.properties[propertySortDirection] = sortDirection

	return c
}

func (c *productQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty(propertySoftDeletedIncluded)
}

func (c *productQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties[propertySoftDeletedIncluded].(bool)
}

func (c *productQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) ProductQueryInterface {
	c.properties[propertySoftDeletedIncluded] = softDeletedIncluded

	return c
}

func (c *productQueryImplementation) HasStatus() bool {
	return c.hasProperty(propertyStatus)
}

func (c *productQueryImplementation) Status() string {
	if !c.HasStatus() {
		return ""
	}

	return c.properties[propertyStatus].(string)
}

func (c *productQueryImplementation) SetStatus(status string) ProductQueryInterface {
	c.properties[propertyStatus] = status

	return c
}

func (c *productQueryImplementation) HasStatusIn() bool {
	return c.hasProperty(propertyStatusIn)
}

func (c *productQueryImplementation) StatusIn() []string {
	if !c.HasStatusIn() {
		return []string{}
	}

	return c.properties[propertyStatusIn].([]string)
}

func (c *productQueryImplementation) SetStatusIn(statusIn []string) ProductQueryInterface {
	c.properties[propertyStatusIn] = statusIn

	return c
}

func (c *productQueryImplementation) HasTitleLike() bool {
	return c.hasProperty(propertyTitleLike)
}

func (c *productQueryImplementation) TitleLike() string {
	if !c.HasTitleLike() {
		return ""
	}

	return c.properties[propertyTitleLike].(string)
}

func (c *productQueryImplementation) SetTitleLike(titleLike string) ProductQueryInterface {
	c.properties[propertyTitleLike] = titleLike

	return c
}

func (c *productQueryImplementation) HasParentID() bool {
	return c.hasProperty(propertyParentID)
}

func (c *productQueryImplementation) ParentID() string {
	if !c.HasParentID() {
		return ""
	}

	return c.properties[propertyParentID].(string)
}

func (c *productQueryImplementation) SetParentID(parentID string) ProductQueryInterface {
	c.properties[propertyParentID] = parentID

	return c
}

func (c *productQueryImplementation) HasMetasIn() bool {
	return c.hasProperty(propertyMetasIn)
}

func (c *productQueryImplementation) MetasIn() map[string]string {
	if !c.HasMetasIn() {
		return map[string]string{}
	}

	return c.properties[propertyMetasIn].(map[string]string)
}

func (c *productQueryImplementation) SetMetasIn(metasIn map[string]string) ProductQueryInterface {
	c.properties[propertyMetasIn] = metasIn

	return c
}

func (c *productQueryImplementation) HasMetasNotIn() bool {
	return c.hasProperty(propertyMetasNotIn)
}

func (c *productQueryImplementation) MetasNotIn() map[string]string {
	if !c.HasMetasNotIn() {
		return map[string]string{}
	}

	return c.properties[propertyMetasNotIn].(map[string]string)
}

func (c *productQueryImplementation) SetMetasNotIn(metasNotIn map[string]string) ProductQueryInterface {
	c.properties[propertyMetasNotIn] = metasNotIn

	return c
}

func (c *productQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
