package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
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
		SetID(uid.HumanUid()).
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

func NewCategoryFromExistingData(data map[string]string) *Category {
	o := &Category{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

func (category *Category) IsActive() bool {
	return category.Status() == CATEGORY_STATUS_ACTIVE
}

func (category *Category) IsDraft() bool {
	return category.Status() == CATEGORY_STATUS_DRAFT
}

func (category *Category) IsInactive() bool {
	return category.Status() == CATEGORY_STATUS_INACTIVE
}

func (category *Category) IsSoftDeleted() bool {
	return category.SoftDeletedAt() != sb.MAX_DATETIME
}

// == GETTERS & SETTERS ========================================================

func (category *Category) CreatedAt() string {
	return category.Get(COLUMN_CREATED_AT)
}

func (category *Category) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.CreatedAt(), carbon.UTC)
}

func (category *Category) SetCreatedAt(createdAt string) CategoryInterface {
	category.Set(COLUMN_CREATED_AT, createdAt)
	return category
}

func (category *Category) CustomerID() string {
	return category.Get(COLUMN_CUSTOMER_ID)
}

func (category *Category) SetCustomerID(id string) CategoryInterface {
	category.Set(COLUMN_CUSTOMER_ID, id)
	return category
}

func (category *Category) Description() string {
	return category.Get(COLUMN_DESCRIPTION)
}

func (category *Category) SetDescription(description string) CategoryInterface {
	category.Set(COLUMN_DESCRIPTION, description)
	return category
}

func (category *Category) ID() string {
	return category.Get(COLUMN_ID)
}

func (category *Category) SetID(id string) CategoryInterface {
	category.Set(COLUMN_ID, id)
	return category
}

func (category *Category) Memo() string {
	return category.Get(COLUMN_MEMO)
}

func (category *Category) SetMemo(memo string) CategoryInterface {
	category.Set(COLUMN_MEMO, memo)
	return category
}

func (category *Category) Metas() (map[string]string, error) {
	metasStr := category.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	var metasJson map[string]string
	errJson := json.Unmarshal([]byte(metasStr), &metasJson)
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return metasJson, nil
}

func (category *Category) Meta(name string) string {
	metas, err := category.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (category *Category) SetMeta(name string, value string) error {
	return category.UpsertMetas(map[string]string{name: value})
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

func (category *Category) UpsertMetas(metas map[string]string) error {
	currentMetas, err := category.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return category.SetMetas(currentMetas)
}

func (category *Category) ParentID() string {
	return category.Get(COLUMN_PARENT_ID)
}

func (category *Category) SetParentID(parentID string) CategoryInterface {
	category.Set(COLUMN_PARENT_ID, parentID)
	return category
}

func (category *Category) SoftDeletedAt() string {
	return category.Get(COLUMN_SOFT_DELETED_AT)
}

func (category *Category) SetSoftDeletedAt(softDeletedAt string) CategoryInterface {
	category.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return category
}

func (category *Category) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.SoftDeletedAt(), carbon.UTC)
}

func (category *Category) Status() string {
	return category.Get(COLUMN_STATUS)
}

func (category *Category) SetStatus(status string) CategoryInterface {
	category.Set(COLUMN_STATUS, status)
	return category
}

func (category *Category) Title() string {
	return category.Get(COLUMN_TITLE)
}

func (category *Category) SetTitle(title string) CategoryInterface {
	category.Set(COLUMN_TITLE, title)
	return category
}

func (category *Category) UpdatedAt() string {
	return category.Get(COLUMN_UPDATED_AT)
}

func (category *Category) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(category.UpdatedAt(), carbon.UTC)
}

func (category *Category) SetUpdatedAt(updatedAt string) CategoryInterface {
	category.Set(COLUMN_UPDATED_AT, updatedAt)
	return category
}
