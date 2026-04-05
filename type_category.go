package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

// == CLASS ====================================================================

// Category represents a product category in the shop store.
// Categories support hierarchical structures with parent-child relationships,
// soft deletion, metadata storage, and status management.
type Category struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// Compile-time interface compliance check
var _ CategoryInterface = (*Category)(nil)

// == CONSTRUCTORS =============================================================

// NewCategory creates a new category with default values:
// - Status: draft
// - ParentID: empty (root category)
// - Description: empty
// - Memo: empty
// - CreatedAt: current UTC time
// - UpdatedAt: current UTC time
// - SoftDeletedAt: max datetime (not deleted)
// - Metas: empty map
func NewCategory() CategoryInterface {
	o := (&Category{}).
		SetID(GenerateShortID()).
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetParentID("").    // By default empty, root category
		SetDescription(""). // By default empty
		SetMemo("").        // By default empty
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	_ = o.SetMetas(map[string]string{})

	return o
}

// NewCategoryFromExistingData creates a category from existing data map.
// Used when hydrating from database or external sources.
func NewCategoryFromExistingData(data map[string]string) CategoryInterface {
	o := &Category{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

// IsActive returns true if the category status is active.
func (category *Category) IsActive() bool {
	return category.GetStatus() == CATEGORY_STATUS_ACTIVE
}

// IsDraft returns true if the category status is draft.
func (category *Category) IsDraft() bool {
	return category.GetStatus() == CATEGORY_STATUS_DRAFT
}

// IsInactive returns true if the category status is inactive.
func (category *Category) IsInactive() bool {
	return category.GetStatus() == CATEGORY_STATUS_INACTIVE
}

// IsSoftDeleted returns true if the category is soft deleted.
func (category *Category) IsSoftDeleted() bool {
	return category.GetSoftDeletedAt() != sb.MAX_DATETIME
}

// IsRoot returns true if the category has no parent (parent_id is empty)
func (category *Category) IsRoot() bool {
	return category.GetParentID() == ""
}

// IsChild returns true if the category has a parent (parent_id is not empty)
func (category *Category) IsChild() bool {
	return category.GetParentID() != ""
}

// == GETTERS & SETTERS ========================================================

// GetCreatedAt returns the creation timestamp as a string.
func (category *Category) GetCreatedAt() string {
	return category.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
func (category *Category) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.GetCreatedAt(), carbon.UTC)
}

// SetCreatedAt sets the creation timestamp.
func (category *Category) SetCreatedAt(createdAt string) CategoryInterface {
	category.Set(COLUMN_CREATED_AT, createdAt)
	return category
}

// GetDescription returns the category description.
func (category *Category) GetDescription() string {
	return category.Get(COLUMN_DESCRIPTION)
}

// SetDescription sets the category description.
func (category *Category) SetDescription(description string) CategoryInterface {
	category.Set(COLUMN_DESCRIPTION, description)
	return category
}

// GetID returns the unique identifier.
func (category *Category) GetID() string {
	return category.Get(COLUMN_ID)
}

// SetID sets the unique identifier.
func (category *Category) SetID(id string) CategoryInterface {
	category.Set(COLUMN_ID, id)
	return category
}

// GetMemo returns the internal memo.
func (category *Category) GetMemo() string {
	return category.Get(COLUMN_MEMO)
}

// SetMemo sets the internal memo.
func (category *Category) SetMemo(memo string) CategoryInterface {
	category.Set(COLUMN_MEMO, memo)
	return category
}

// GetMetas returns all metadata as a map. Returns empty map if no metas stored.
func (category *Category) GetMetas() (map[string]string, error) {
	metasStr := category.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson := map[string]string{}
	errJson := json.Unmarshal([]byte(metasStr), &metasJson)
	if errJson != nil {
		return map[string]string{}, errJson
	}

	if metasJson == nil {
		metasJson = map[string]string{}
	}

	return metasJson, nil
}

// GetMeta returns a specific metadata value by name. Returns empty string if not found.
func (category *Category) GetMeta(name string) string {
	metas, err := category.GetMetas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// SetMeta sets a single metadata value.
func (category *Category) SetMeta(name string, value string) error {
	return category.MetasUpsert(map[string]string{name: value})
}

// SetMetas replaces all metadata with the provided map.
// Warning: this overwrites any existing metadata.
func (category *Category) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	category.Set(COLUMN_METAS, string(mapString))
	return nil
}

// MetasUpsert merges the provided metadata with existing values.
func (category *Category) MetasUpsert(metas map[string]string) error {
	currentMetas, err := category.GetMetas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return category.SetMetas(currentMetas)
}

// MetaRemove removes a single metadata entry.
func (category *Category) MetaRemove(name string) error {
	metas, err := category.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return category.SetMetas(metas)
}

// MetasRemove removes multiple metadata entries.
func (category *Category) MetasRemove(names []string) error {
	for _, name := range names {
		if err := category.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

// GetParentID returns the parent category ID (empty if root category).
func (category *Category) GetParentID() string {
	return category.Get(COLUMN_PARENT_ID)
}

// SetParentID sets the parent category ID. Empty string makes it a root category.
func (category *Category) SetParentID(parentID string) CategoryInterface {
	category.Set(COLUMN_PARENT_ID, parentID)
	return category
}

// GetSoftDeletedAt returns the soft deletion timestamp.
func (category *Category) GetSoftDeletedAt() string {
	return category.Get(COLUMN_SOFT_DELETED_AT)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (category *Category) SetSoftDeletedAt(softDeletedAt string) CategoryInterface {
	category.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return category
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
func (category *Category) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.GetSoftDeletedAt(), carbon.UTC)
}

// GetStatus returns the current status.
func (category *Category) GetStatus() string {
	return category.Get(COLUMN_STATUS)
}

// SetStatus sets the current status.
func (category *Category) SetStatus(status string) CategoryInterface {
	category.Set(COLUMN_STATUS, status)
	return category
}

// GetTitle returns the category title.
func (category *Category) GetTitle() string {
	return category.Get(COLUMN_TITLE)
}

// SetTitle sets the category title.
func (category *Category) SetTitle(title string) CategoryInterface {
	category.Set(COLUMN_TITLE, title)
	return category
}

// GetUpdatedAt returns the last update timestamp.
func (category *Category) GetUpdatedAt() string {
	return category.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
func (category *Category) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.GetUpdatedAt(), carbon.UTC)
}

// SetUpdatedAt sets the last update timestamp.
func (category *Category) SetUpdatedAt(updatedAt string) CategoryInterface {
	category.Set(COLUMN_UPDATED_AT, updatedAt)
	return category
}
