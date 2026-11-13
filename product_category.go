package shopstore

import (
	"uuid"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == CLASS ==================================================================

// ProductCategory represents a product-category assignment in the shop store.
// It is a join entity linking a product to a category (many-to-many).
type ProductCategory struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// Compile-time interface compliance check
var _ ProductCategoryInterface = (*ProductCategory)(nil)

// == CONSTRUCTORS ===========================================================

// NewProductCategory creates a new product-category link with default values:
// - ProductID: empty
// - CategoryID: empty
// - CreatedAt: current UTC time
// - UpdatedAt: current UTC time
// - SoftDeletedAt: max datetime (not deleted)
func NewProductCategory() ProductCategoryInterface {
	s := (&ProductCategory{}).
		SetID(uuid.NewV7().String()).
		SetProductID("").
		SetCategoryID("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(MAX_DATETIME)

	return s
}

// NewProductCategoryFromExistingData creates a product-category link from
// an existing data map. Used when hydrating from database or external sources.
func NewProductCategoryFromExistingData(data map[string]string) ProductCategoryInterface {
	o := &ProductCategory{}
	o.Hydrate(data)
	return o
}

// == SETTERS AND GETTERS ====================================================

// GetCategoryID returns the linked category identifier.
func (s *ProductCategory) GetCategoryID() string {
	return s.Get(COLUMN_CATEGORY_ID)
}

// SetCategoryID sets the linked category identifier.
func (s *ProductCategory) SetCategoryID(categoryID string) ProductCategoryInterface {
	s.Set(COLUMN_CATEGORY_ID, categoryID)
	return s
}

// GetCreatedAt returns the creation timestamp as a string.
func (s *ProductCategory) GetCreatedAt() string {
	return s.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
func (s *ProductCategory) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := s.GetCreatedAt()
	return carbon.Parse(createdAt)
}

// SetCreatedAt sets the creation timestamp.
func (s *ProductCategory) SetCreatedAt(createdAt string) ProductCategoryInterface {
	s.Set(COLUMN_CREATED_AT, createdAt)
	return s
}

// GetID returns the unique identifier.
func (s *ProductCategory) GetID() string {
	return s.Get(COLUMN_ID)
}

// SetID sets the unique identifier.
func (s *ProductCategory) SetID(id string) ProductCategoryInterface {
	s.Set(COLUMN_ID, id)
	return s
}

// GetProductID returns the linked product identifier.
func (s *ProductCategory) GetProductID() string {
	return s.Get(COLUMN_PRODUCT_ID)
}

// SetProductID sets the linked product identifier.
func (s *ProductCategory) SetProductID(productID string) ProductCategoryInterface {
	s.Set(COLUMN_PRODUCT_ID, productID)
	return s
}

// GetSoftDeletedAt returns the soft deletion timestamp.
func (s *ProductCategory) GetSoftDeletedAt() string {
	return s.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
func (s *ProductCategory) GetSoftDeletedAtCarbon() *carbon.Carbon {
	deletedAt := s.GetSoftDeletedAt()
	return carbon.Parse(deletedAt)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (s *ProductCategory) SetSoftDeletedAt(deletedAt string) ProductCategoryInterface {
	s.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return s
}

// GetUpdatedAt returns the last update timestamp.
func (s *ProductCategory) GetUpdatedAt() string {
	return s.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
func (s *ProductCategory) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := s.GetUpdatedAt()
	return carbon.Parse(updatedAt)
}

// SetUpdatedAt sets the last update timestamp.
func (s *ProductCategory) SetUpdatedAt(updatedAt string) ProductCategoryInterface {
	s.Set(COLUMN_UPDATED_AT, updatedAt)
	return s
}

// IsSoftDeleted returns true if the link is soft deleted.
func (s *ProductCategory) IsSoftDeleted() bool {
	return s.GetSoftDeletedAt() != MAX_DATETIME
}

// MarkAsNotDirty resets the dirty state, clearing all change tracking.
func (s *ProductCategory) MarkAsNotDirty() {
	s.DataObject.MarkAsNotDirty()
}
