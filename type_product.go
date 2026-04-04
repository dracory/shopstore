package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/str"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == CLASS ====================================================================

type Product struct {
	dataobject.DataObject
}

// == INTERFACES ===============================================================

var _ ProductInterface = (*Product)(nil)

// == CONSTRUCTORS =============================================================

func NewProduct() ProductInterface {
	o := (&Product{}).
		SetID(GenerateShortID()).
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("").
		SetDescription("").
		SetShortDescription("").
		SetQuantityInt(0). // By default 0
		SetPriceFloat(0).  // Free. By default
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	_ = o.SetMetas(map[string]string{})

	return o
}

func NewProductFromExistingData(data map[string]string) ProductInterface {
	o := &Product{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

func (product *Product) IsActive() bool {
	return product.GetStatus() == PRODUCT_STATUS_ACTIVE
}

func (product *Product) IsDisabled() bool {
	return product.GetStatus() == PRODUCT_STATUS_DISABLED
}

func (product *Product) IsDraft() bool {
	return product.GetStatus() == PRODUCT_STATUS_DRAFT
}

func (product *Product) IsSoftDeleted() bool {
	return product.GetSoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

func (product *Product) IsFree() bool {
	return product.GetPriceFloat() <= 0
}

// HasStock returns true if the product quantity is greater than 0
func (product *Product) HasStock() bool {
	return product.GetQuantityInt() > 0
}

// IsOutOfStock returns true if the product quantity is less than or equal to 0
func (product *Product) IsOutOfStock() bool {
	return product.GetQuantityInt() <= 0
}

// IsPaid returns true if the product price is greater than 0
func (product *Product) IsPaid() bool {
	return product.GetPriceFloat() > 0
}

func (product *Product) Slug() string {
	title := product.GetTitle()
	return str.Slugify(title, '-')
}

// == GETTERS & SETTERS ========================================================

func (product *Product) GetCreatedAt() string {
	return product.Get(COLUMN_CREATED_AT)
}

func (product *Product) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(product.GetCreatedAt(), carbon.UTC)
}

func (product *Product) SetCreatedAt(createdAt string) ProductInterface {
	product.Set(COLUMN_CREATED_AT, createdAt)
	return product
}

func (product *Product) GetDescription() string {
	return product.Get(COLUMN_DESCRIPTION)
}

func (product *Product) SetDescription(description string) ProductInterface {
	product.Set(COLUMN_DESCRIPTION, description)
	return product
}

func (product *Product) GetID() string {
	return product.Get(COLUMN_ID)
}

func (product *Product) SetID(id string) ProductInterface {
	product.Set(COLUMN_ID, id)
	return product
}

func (product *Product) GetMemo() string {
	return product.Get(COLUMN_MEMO)
}

func (product *Product) SetMemo(memo string) ProductInterface {
	product.Set(COLUMN_MEMO, memo)
	return product
}

func (product *Product) GetMetas() (map[string]string, error) {
	metasStr := product.Get(COLUMN_METAS)

	if metasStr == "" || metasStr == "null" {
		metasStr = "{}"
	}

	var metasJson map[string]string
	errJson := json.Unmarshal([]byte(metasStr), &metasJson)
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return metasJson, nil
}

func (product *Product) GetMeta(name string) string {
	metas, err := product.GetMetas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (product *Product) SetMeta(name string, value string) error {
	return product.MetasUpsert(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (product *Product) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	product.Set(COLUMN_METAS, string(mapString))
	return nil
}

func (product *Product) MetasUpsert(metas map[string]string) error {
	currentMetas, err := product.GetMetas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return product.SetMetas(currentMetas)
}

func (product *Product) MetaRemove(name string) error {
	metas, err := product.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return product.SetMetas(metas)
}

func (product *Product) MetasRemove(names []string) error {
	for _, name := range names {
		if err := product.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

func (product *Product) GetPrice() string {
	return product.Get(COLUMN_PRICE)
}

func (product *Product) SetPrice(price string) ProductInterface {
	product.Set(COLUMN_PRICE, price)
	return product
}

func (product *Product) GetPriceFloat() float64 {
	price := cast.ToFloat64(product.Get(COLUMN_PRICE))
	return price
}

func (product *Product) SetPriceFloat(price float64) ProductInterface {
	product.SetPrice(cast.ToString(price))
	return product
}

func (product *Product) GetQuantity() string {
	return product.Get(COLUMN_QUANTITY)
}

func (product *Product) SetQuantity(quantity string) ProductInterface {
	product.Set(COLUMN_QUANTITY, quantity)
	return product
}

func (product *Product) GetQuantityInt() int64 {
	quantity := cast.ToInt64(product.GetQuantity())
	return quantity
}

func (product *Product) SetQuantityInt(quantity int64) ProductInterface {
	product.SetQuantity(cast.ToString(quantity))
	return product
}

func (product *Product) GetShortDescription() string {
	return product.Get(COLUMN_SHORT_DESCRIPTION)
}

func (product *Product) SetShortDescription(shortDescription string) ProductInterface {
	product.Set(COLUMN_SHORT_DESCRIPTION, shortDescription)
	return product
}

func (product *Product) GetSoftDeletedAt() string {
	return product.Get(COLUMN_SOFT_DELETED_AT)
}

func (product *Product) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(product.GetSoftDeletedAt(), carbon.UTC)
}

func (product *Product) SetSoftDeletedAt(deletedAt string) ProductInterface {
	product.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return product
}

func (product *Product) GetStatus() string {
	return product.Get(COLUMN_STATUS)
}

func (product *Product) SetStatus(status string) ProductInterface {
	product.Set(COLUMN_STATUS, status)
	return product
}

func (product *Product) GetTitle() string {
	return product.Get(COLUMN_TITLE)
}

func (product *Product) SetTitle(title string) ProductInterface {
	product.Set(COLUMN_TITLE, title)
	return product
}

func (product *Product) GetUpdatedAt() string {
	return product.Get(COLUMN_UPDATED_AT)
}

func (product *Product) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(product.GetUpdatedAt(), carbon.UTC)
}

func (product *Product) SetUpdatedAt(updatedAt string) ProductInterface {
	product.Set(COLUMN_UPDATED_AT, updatedAt)
	return product
}
