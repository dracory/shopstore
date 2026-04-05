package shopstore

import (
	"encoding/json"
	"strings"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == TYPE ====================================================================

// Media represents a media file (image, video, etc) in the shop store.
// Media files are associated with entities via EntityID, support sequencing for ordering,
// soft deletion, metadata storage, and type classification.
type Media struct {
	dataobject.DataObject
}

// == INTERFACES ===============================================================

// Compile-time interface compliance check
var _ MediaInterface = (*Media)(nil)

// == CONSTRUCTORS =============================================================

// NewMedia creates a new media entity with default values:
// - Status: draft
// - Title: empty
// - Description: empty
// - Memo: empty
// - CreatedAt: current UTC time
// - UpdatedAt: current UTC time
// - SoftDeletedAt: max datetime (not deleted)
// - Metas: empty map
func NewMedia() MediaInterface {
	o := (&Media{}).
		SetID(GenerateShortID()).
		SetStatus(MEDIA_STATUS_DRAFT).
		SetTitle("").       // By default empty, root category
		SetDescription(""). // By default empty
		SetMemo("").        // By default empty
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	_ = o.SetMetas(map[string]string{})

	return o
}

// NewMediaFromExistingData creates a media entity from existing data map.
// Used when hydrating from database or external sources.
func NewMediaFromExistingData(data map[string]string) MediaInterface {
	o := &Media{}
	o.Hydrate(data)
	return o
}

// == SETTESR AND GETTERS =====================================================

// GetCreatedAt returns the creation timestamp as a string.
func (m *Media) GetCreatedAt() string {
	return m.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
func (m *Media) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(m.GetCreatedAt(), carbon.UTC)
}

// SetCreatedAt sets the creation timestamp.
func (m *Media) SetCreatedAt(createdAt string) MediaInterface {
	m.Set(COLUMN_CREATED_AT, createdAt)
	return m
}

// GetDescription returns the media description.
func (m *Media) GetDescription() string {
	return m.Get(COLUMN_DESCRIPTION)
}

// SetDescription sets the media description.
func (m *Media) SetDescription(description string) MediaInterface {
	m.Set(COLUMN_DESCRIPTION, description)
	return m
}

// GetEntityID returns the associated entity ID.
func (m *Media) GetEntityID() string {
	return m.Get(COLUMN_ENTITY_ID)
}

// SetEntityID sets the associated entity ID.
func (m *Media) SetEntityID(entityID string) MediaInterface {
	m.Set(COLUMN_ENTITY_ID, entityID)
	return m
}

// GetID returns the unique identifier.
func (m *Media) GetID() string {
	return m.Get(COLUMN_ID)
}

// SetID sets the unique identifier.
func (m *Media) SetID(id string) MediaInterface {
	m.Set(COLUMN_ID, id)
	return m
}

// GetMemo returns the internal memo.
func (m *Media) GetMemo() string {
	return m.Get(COLUMN_MEMO)
}

// SetMemo sets the internal memo.
func (m *Media) SetMemo(memo string) MediaInterface {
	m.Set(COLUMN_MEMO, memo)
	return m
}

// GetMetas returns all metadata as a map. Returns empty map if no metas stored.
func (m *Media) GetMetas() (map[string]string, error) {
	metasStr := m.Get(COLUMN_METAS)

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
func (m *Media) GetMeta(name string) string {
	metas, err := m.GetMetas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// SetMeta sets a single metadata value.
func (m *Media) SetMeta(name string, value string) error {
	return m.MetasUpsert(map[string]string{name: value})
}

// SetMetas replaces all metadata with the provided map.
// Warning: this overwrites any existing metadata.
func (m *Media) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	m.Set(COLUMN_METAS, string(mapString))
	return nil
}

// MetasUpsert merges the provided metadata with existing values.
func (m *Media) MetasUpsert(metas map[string]string) error {
	currentMetas, err := m.GetMetas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return m.SetMetas(currentMetas)
}

// MetaRemove removes a single metadata entry.
func (m *Media) MetaRemove(name string) error {
	metas, err := m.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return m.SetMetas(metas)
}

// MetasRemove removes multiple metadata entries.
func (m *Media) MetasRemove(names []string) error {
	for _, name := range names {
		if err := m.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

// GetSequence returns the display sequence/order.
func (m *Media) GetSequence() int {
	return cast.ToInt(m.Get(COLUMN_SEQUENCE))
}

// SetSequence sets the display sequence/order.
func (m *Media) SetSequence(sequence int) MediaInterface {
	m.Set(COLUMN_SEQUENCE, cast.ToString(sequence))
	return m
}

// GetStatus returns the current status.
func (m *Media) GetStatus() string {
	return m.Get(COLUMN_STATUS)
}

// SetStatus sets the current status.
func (m *Media) SetStatus(status string) MediaInterface {
	m.Set(COLUMN_STATUS, status)
	return m
}

// GetSoftDeletedAt returns the soft deletion timestamp.
func (m *Media) GetSoftDeletedAt() string {
	return m.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
func (m *Media) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(m.GetSoftDeletedAt(), carbon.UTC)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (m *Media) SetSoftDeletedAt(softDeletedAt string) MediaInterface {
	m.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return m
}

// GetTitle returns the media title.
func (m *Media) GetTitle() string {
	return m.Get(COLUMN_TITLE)
}

// SetTitle sets the media title.
func (m *Media) SetTitle(title string) MediaInterface {
	m.Set(COLUMN_TITLE, title)
	return m
}

// GetType returns the media type (e.g., image/jpeg, video/mp4).
func (m *Media) GetType() string {
	return m.Get(COLUMN_MEDIA_TYPE)
}

// SetType sets the media type.
func (m *Media) SetType(type_ string) MediaInterface {
	m.Set(COLUMN_MEDIA_TYPE, type_)
	return m
}

// GetUpdatedAt returns the last update timestamp.
func (m *Media) GetUpdatedAt() string {
	return m.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
func (m *Media) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(m.GetUpdatedAt(), carbon.UTC)
}

// SetUpdatedAt sets the last update timestamp.
func (m *Media) SetUpdatedAt(updatedAt string) MediaInterface {
	m.Set(COLUMN_UPDATED_AT, updatedAt)
	return m
}

// GetURL returns the media file URL.
func (m *Media) GetURL() string {
	return m.Get(COLUMN_MEDIA_URL)
}

// SetURL sets the media file URL.
func (m *Media) SetURL(url string) MediaInterface {
	m.Set(COLUMN_MEDIA_URL, url)
	return m
}

// IsActive returns true if the media status is active.
func (m *Media) IsActive() bool {
	return m.GetStatus() == MEDIA_STATUS_ACTIVE
}

// IsDraft returns true if the media status is draft.
func (m *Media) IsDraft() bool {
	return m.GetStatus() == MEDIA_STATUS_DRAFT
}

// IsInactive returns true if the media status is inactive.
func (m *Media) IsInactive() bool {
	return m.GetStatus() == MEDIA_STATUS_INACTIVE
}

// IsSoftDeleted returns true if the media is soft deleted.
func (m *Media) IsSoftDeleted() bool {
	return m.GetSoftDeletedAt() != sb.MAX_DATETIME
}

// IsImage returns true if the media type starts with "image/".
func (m *Media) IsImage() bool {
	return strings.HasPrefix(m.GetType(), "image/")
}

// IsVideo returns true if the media type starts with "video/".
func (m *Media) IsVideo() bool {
	return strings.HasPrefix(m.GetType(), "video/")
}
