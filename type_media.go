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

type Media struct {
	dataobject.DataObject
}

// == INTERFACE ===============================================================

var _ MediaInterface = (*Media)(nil)

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

func NewMediaFromExistingData(data map[string]string) MediaInterface {
	o := &Media{}
	o.Hydrate(data)
	return o
}

// == SETTESR AND GETTERS =====================================================

func (m *Media) GetCreatedAt() string {
	return m.Get(COLUMN_CREATED_AT)
}

func (m *Media) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(m.GetCreatedAt(), carbon.UTC)
}

func (m *Media) SetCreatedAt(createdAt string) MediaInterface {
	m.Set(COLUMN_CREATED_AT, createdAt)
	return m
}

func (m *Media) GetDescription() string {
	return m.Get(COLUMN_DESCRIPTION)
}

func (m *Media) SetDescription(description string) MediaInterface {
	m.Set(COLUMN_DESCRIPTION, description)
	return m
}

func (m *Media) GetEntityID() string {
	return m.Get(COLUMN_ENTITY_ID)
}

func (m *Media) SetEntityID(entityID string) MediaInterface {
	m.Set(COLUMN_ENTITY_ID, entityID)
	return m
}

func (m *Media) GetID() string {
	return m.Get(COLUMN_ID)
}

func (m *Media) SetID(id string) MediaInterface {
	m.Set(COLUMN_ID, id)
	return m
}

func (m *Media) GetMemo() string {
	return m.Get(COLUMN_MEMO)
}

func (m *Media) SetMemo(memo string) MediaInterface {
	m.Set(COLUMN_MEMO, memo)
	return m
}

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

func (m *Media) SetMeta(name string, value string) error {
	return m.MetasUpsert(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (m *Media) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	m.Set(COLUMN_METAS, string(mapString))
	return nil
}

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

func (m *Media) MetaRemove(name string) error {
	metas, err := m.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return m.SetMetas(metas)
}

func (m *Media) MetasRemove(names []string) error {
	for _, name := range names {
		if err := m.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

func (m *Media) GetSequence() int {
	return cast.ToInt(m.Get(COLUMN_SEQUENCE))
}

func (m *Media) SetSequence(sequence int) MediaInterface {
	m.Set(COLUMN_SEQUENCE, cast.ToString(sequence))
	return m
}

func (m *Media) GetStatus() string {
	return m.Get(COLUMN_STATUS)
}

func (m *Media) SetStatus(status string) MediaInterface {
	m.Set(COLUMN_STATUS, status)
	return m
}

func (m *Media) GetSoftDeletedAt() string {
	return m.Get(COLUMN_SOFT_DELETED_AT)
}

func (m *Media) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(m.GetSoftDeletedAt(), carbon.UTC)
}

func (m *Media) SetSoftDeletedAt(softDeletedAt string) MediaInterface {
	m.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return m
}

func (m *Media) GetTitle() string {
	return m.Get(COLUMN_TITLE)
}

func (m *Media) SetTitle(title string) MediaInterface {
	m.Set(COLUMN_TITLE, title)
	return m
}

func (m *Media) GetType() string {
	return m.Get(COLUMN_MEDIA_TYPE)
}

func (m *Media) SetType(type_ string) MediaInterface {
	m.Set(COLUMN_MEDIA_TYPE, type_)
	return m
}

func (m *Media) GetUpdatedAt() string {
	return m.Get(COLUMN_UPDATED_AT)
}

func (m *Media) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(m.GetUpdatedAt(), carbon.UTC)
}

func (m *Media) SetUpdatedAt(updatedAt string) MediaInterface {
	m.Set(COLUMN_UPDATED_AT, updatedAt)
	return m
}

func (m *Media) GetURL() string {
	return m.Get(COLUMN_MEDIA_URL)
}

func (m *Media) SetURL(url string) MediaInterface {
	m.Set(COLUMN_MEDIA_URL, url)
	return m
}

// IsActive returns true if the media status is active
func (m *Media) IsActive() bool {
	return m.GetStatus() == MEDIA_STATUS_ACTIVE
}

// IsDraft returns true if the media status is draft
func (m *Media) IsDraft() bool {
	return m.GetStatus() == MEDIA_STATUS_DRAFT
}

// IsInactive returns true if the media status is inactive
func (m *Media) IsInactive() bool {
	return m.GetStatus() == MEDIA_STATUS_INACTIVE
}

// IsSoftDeleted returns true if the media is soft deleted
func (m *Media) IsSoftDeleted() bool {
	return m.GetSoftDeletedAt() != sb.MAX_DATETIME
}

// IsImage returns true if the media type starts with "image/"
func (m *Media) IsImage() bool {
	return strings.HasPrefix(m.GetType(), "image/")
}

// IsVideo returns true if the media type starts with "video/"
func (m *Media) IsVideo() bool {
	return strings.HasPrefix(m.GetType(), "video/")
}
