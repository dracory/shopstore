package shopstore

import (
	"uuid"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == CLASS ==================================================================

// Setting represents a key/value setting in the shop store.
// Settings hold shop-wide configuration such as the display currency.
type Setting struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// Compile-time interface compliance check
var _ SettingInterface = (*Setting)(nil)

// == CONSTRUCTORS ===========================================================

// NewSetting creates a new setting with default values:
// - Key: empty
// - Value: empty
// - CreatedAt: current UTC time
// - UpdatedAt: current UTC time
// - SoftDeletedAt: max datetime (not deleted)
func NewSetting() SettingInterface {
	s := (&Setting{}).
		SetID(uuid.NewV7().String()).
		SetKey("").
		SetValue("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(MAX_DATETIME)

	return s
}

// NewSettingFromExistingData creates a setting from existing data map.
// Used when hydrating from database or external sources.
func NewSettingFromExistingData(data map[string]string) SettingInterface {
	o := &Setting{}
	o.Hydrate(data)
	return o
}

// == SETTERS AND GETTERS ====================================================

// GetCreatedAt returns the creation timestamp as a string.
func (s *Setting) GetCreatedAt() string {
	return s.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
func (s *Setting) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := s.GetCreatedAt()
	return carbon.Parse(createdAt)
}

// SetCreatedAt sets the creation timestamp.
func (s *Setting) SetCreatedAt(createdAt string) SettingInterface {
	s.Set(COLUMN_CREATED_AT, createdAt)
	return s
}

// GetID returns the unique identifier.
func (s *Setting) GetID() string {
	return s.Get(COLUMN_ID)
}

// SetID sets the unique identifier.
func (s *Setting) SetID(id string) SettingInterface {
	s.Set(COLUMN_ID, id)
	return s
}

// GetKey returns the setting key.
func (s *Setting) GetKey() string {
	return s.Get(COLUMN_SETTING_KEY)
}

// SetKey sets the setting key.
func (s *Setting) SetKey(key string) SettingInterface {
	s.Set(COLUMN_SETTING_KEY, key)
	return s
}

// GetSoftDeletedAt returns the soft deletion timestamp.
func (s *Setting) GetSoftDeletedAt() string {
	return s.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
func (s *Setting) GetSoftDeletedAtCarbon() *carbon.Carbon {
	deletedAt := s.GetSoftDeletedAt()
	return carbon.Parse(deletedAt)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (s *Setting) SetSoftDeletedAt(deletedAt string) SettingInterface {
	s.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return s
}

// GetUpdatedAt returns the last update timestamp.
func (s *Setting) GetUpdatedAt() string {
	return s.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
func (s *Setting) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := s.GetUpdatedAt()
	return carbon.Parse(updatedAt)
}

// SetUpdatedAt sets the last update timestamp.
func (s *Setting) SetUpdatedAt(updatedAt string) SettingInterface {
	s.Set(COLUMN_UPDATED_AT, updatedAt)
	return s
}

// GetValue returns the setting value.
func (s *Setting) GetValue() string {
	return s.Get(COLUMN_SETTING_VALUE)
}

// SetValue sets the setting value.
func (s *Setting) SetValue(value string) SettingInterface {
	s.Set(COLUMN_SETTING_VALUE, value)
	return s
}

// IsSoftDeleted returns true if the setting is soft deleted.
func (s *Setting) IsSoftDeleted() bool {
	return s.GetSoftDeletedAt() != MAX_DATETIME
}

// MarkAsNotDirty resets the dirty state, clearing all change tracking.
func (s *Setting) MarkAsNotDirty() {
	s.DataObject.MarkAsNotDirty()
}
