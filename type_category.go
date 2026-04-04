package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

// == CLASS ====================================================================

type Category struct {
	dataobject.DataObject
}

var _ CategoryInterface = (*Category)(nil)

// == CONSTRUCTORS =============================================================

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

func NewCategoryFromExistingData(data map[string]string) CategoryInterface {
	o := &Category{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

func (category *Category) IsActive() bool {
	return category.GetStatus() == CATEGORY_STATUS_ACTIVE
}

func (category *Category) IsDraft() bool {
	return category.GetStatus() == CATEGORY_STATUS_DRAFT
}

func (category *Category) IsInactive() bool {
	return category.GetStatus() == CATEGORY_STATUS_INACTIVE
}

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

func (category *Category) GetCreatedAt() string {
	return category.Get(COLUMN_CREATED_AT)
}

func (category *Category) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.GetCreatedAt(), carbon.UTC)
}

func (category *Category) SetCreatedAt(createdAt string) CategoryInterface {
	category.Set(COLUMN_CREATED_AT, createdAt)
	return category
}

func (category *Category) GetDescription() string {
	return category.Get(COLUMN_DESCRIPTION)
}

func (category *Category) SetDescription(description string) CategoryInterface {
	category.Set(COLUMN_DESCRIPTION, description)
	return category
}

func (category *Category) GetID() string {
	return category.Get(COLUMN_ID)
}

func (category *Category) SetID(id string) CategoryInterface {
	category.Set(COLUMN_ID, id)
	return category
}

func (category *Category) GetMemo() string {
	return category.Get(COLUMN_MEMO)
}

func (category *Category) SetMemo(memo string) CategoryInterface {
	category.Set(COLUMN_MEMO, memo)
	return category
}

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

func (category *Category) SetMeta(name string, value string) error {
	return category.MetasUpsert(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (category *Category) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	category.Set(COLUMN_METAS, string(mapString))
	return nil
}

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

func (category *Category) MetaRemove(name string) error {
	metas, err := category.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return category.SetMetas(metas)
}

func (category *Category) MetasRemove(names []string) error {
	for _, name := range names {
		if err := category.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

func (category *Category) GetParentID() string {
	return category.Get(COLUMN_PARENT_ID)
}

func (category *Category) SetParentID(parentID string) CategoryInterface {
	category.Set(COLUMN_PARENT_ID, parentID)
	return category
}

func (category *Category) GetSoftDeletedAt() string {
	return category.Get(COLUMN_SOFT_DELETED_AT)
}

func (category *Category) SetSoftDeletedAt(softDeletedAt string) CategoryInterface {
	category.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return category
}

func (category *Category) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.GetSoftDeletedAt(), carbon.UTC)
}

func (category *Category) GetStatus() string {
	return category.Get(COLUMN_STATUS)
}

func (category *Category) SetStatus(status string) CategoryInterface {
	category.Set(COLUMN_STATUS, status)
	return category
}

func (category *Category) GetTitle() string {
	return category.Get(COLUMN_TITLE)
}

func (category *Category) SetTitle(title string) CategoryInterface {
	category.Set(COLUMN_TITLE, title)
	return category
}

func (category *Category) GetUpdatedAt() string {
	return category.Get(COLUMN_UPDATED_AT)
}

func (category *Category) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.GetUpdatedAt(), carbon.UTC)
}

func (category *Category) SetUpdatedAt(updatedAt string) CategoryInterface {
	category.Set(COLUMN_UPDATED_AT, updatedAt)
	return category
}
