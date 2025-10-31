package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
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
		SetID(uid.HumanUid()).
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("").       // By default empty, root category
		SetDescription(""). // By default empty
		SetMemo("").        // By default empty
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	_ = o.SetMetas(map[string]string{})

	return o
}

func NewMediaFromExistingData(data map[string]string) *Media {
	o := &Media{}
	o.Hydrate(data)
	return o
}

// == SETTESR AND GETTERS =====================================================

func (o *Media) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *Media) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *Media) SetCreatedAt(createdAt string) MediaInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *Media) Description() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *Media) SetDescription(description string) MediaInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

func (o *Media) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *Media) SetEntityID(entityID string) MediaInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *Media) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *Media) SetID(id string) MediaInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *Media) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *Media) SetMemo(memo string) MediaInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (order *Media) Metas() (map[string]string, error) {
	metasStr := order.Get(COLUMN_METAS)

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

func (order *Media) Meta(name string) string {
	metas, err := order.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (order *Media) SetMeta(name string, value string) error {
	return order.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (order *Media) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	order.Set(COLUMN_METAS, string(mapString))
	return nil
}

func (order *Media) UpsertMetas(metas map[string]string) error {
	currentMetas, err := order.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return order.SetMetas(currentMetas)
}

func (o *Media) Sequence() int {
	return cast.ToInt(o.Get(COLUMN_SEQUENCE))
}

func (o *Media) SetSequence(sequence int) MediaInterface {
	o.Set(COLUMN_SEQUENCE, cast.ToString(sequence))
	return o
}

func (o *Media) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *Media) SetStatus(status string) MediaInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *Media) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *Media) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *Media) SetSoftDeletedAt(softDeletedAt string) MediaInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *Media) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *Media) SetTitle(title string) MediaInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *Media) Type() string {
	return o.Get(COLUMN_MEDIA_TYPE)
}

func (o *Media) SetType(type_ string) MediaInterface {
	o.Set(COLUMN_MEDIA_TYPE, type_)
	return o
}

func (o *Media) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *Media) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

func (o *Media) SetUpdatedAt(updatedAt string) MediaInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *Media) URL() string {
	return o.Get(COLUMN_MEDIA_URL)
}

func (o *Media) SetURL(url string) MediaInterface {
	o.Set(COLUMN_MEDIA_URL, url)
	return o
}
