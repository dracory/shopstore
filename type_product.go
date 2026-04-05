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

// Product represents a product in the shop store.
// Products support parent-child relationships (variants), pricing, stock management,
// variant matrix configuration, soft deletion, metadata storage, and status management.
type Product struct {
	dataobject.DataObject
}

// == INTERFACES ===============================================================

// Compile-time interface compliance check
var _ ProductInterface = (*Product)(nil)

// == CONSTRUCTORS =============================================================

// NewProduct creates a new product with default values:
// - Status: draft
// - Title: empty
// - Description: empty
// - ShortDescription: empty
// - Quantity: 0
// - Price: 0.00 (free)
// - ParentID: empty (not a variant)
// - Memo: empty
// - CreatedAt: current UTC time
// - UpdatedAt: current UTC time
// - SoftDeletedAt: max datetime (not deleted)
// - Metas: empty map
// - VariantMatrixSchema: empty
// - VariantMatrixValues: empty map
func NewProduct() ProductInterface {
	o := (&Product{}).
		SetID(GenerateShortID()).
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("").
		SetDescription("").
		SetShortDescription("").
		SetQuantityInt(0). // By default 0
		SetPriceFloat(0).  // Free. By default
		SetParentID("").   // No parent by default (not a variant)
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	_ = o.SetMetas(map[string]string{})
	_ = o.SetVariantMatrixSchema(VariantMatrixSchema{})
	_ = o.SetVariantMatrixValues(map[string]string{})

	return o
}

// NewProductFromExistingData creates a product from existing data map.
// Used when hydrating from database or external sources.
func NewProductFromExistingData(data map[string]string) ProductInterface {
	o := &Product{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

// IsActive returns true if the product status is active.
func (product *Product) IsActive() bool {
	return product.GetStatus() == PRODUCT_STATUS_ACTIVE
}

// IsDisabled returns true if the product status is disabled.
func (product *Product) IsDisabled() bool {
	return product.GetStatus() == PRODUCT_STATUS_DISABLED
}

// IsDraft returns true if the product status is draft.
func (product *Product) IsDraft() bool {
	return product.GetStatus() == PRODUCT_STATUS_DRAFT
}

// IsSoftDeleted returns true if the product is soft deleted.
func (product *Product) IsSoftDeleted() bool {
	return product.GetSoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// IsVariant returns true if this product is a variant of another product.
func (product *Product) IsVariant() bool {
	return product.GetParentID() != ""
}

// IsParent returns true if this product has variants (has variant matrix schema).
func (product *Product) IsParent() bool {
	return product.HasVariantMatrixSchema()
}

// GetParentID returns the parent product ID (empty if not a variant).
func (product *Product) GetParentID() string {
	return product.Get(COLUMN_PARENT_ID)
}

// SetParentID sets the parent product ID.
func (product *Product) SetParentID(parentID string) ProductInterface {
	product.Set(COLUMN_PARENT_ID, parentID)
	return product
}

// SetVariantMatrixSchema sets the variant matrix schema configuration.
func (product *Product) SetVariantMatrixSchema(schema VariantMatrixSchema) error {
	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	product.Set(COLUMN_VARIANT_MATRIX_SCHEMA, string(jsonBytes))
	return nil
}

// SetVariantMatrixValues sets the variant attribute values for this product.
func (product *Product) SetVariantMatrixValues(values map[string]string) error {
	jsonBytes, err := json.Marshal(values)
	if err != nil {
		return err
	}
	product.Set(COLUMN_VARIANT_MATRIX_VALUES, string(jsonBytes))
	return nil
}

// GetVariantMatrixValues returns the variant attribute values for this product.
func (product *Product) GetVariantMatrixValues() (map[string]string, error) {
	valuesJSON := product.Get(COLUMN_VARIANT_MATRIX_VALUES)
	if valuesJSON == "" || valuesJSON == "null" {
		return map[string]string{}, nil
	}
	var values map[string]string
	err := json.Unmarshal([]byte(valuesJSON), &values)
	return values, err
}

// GetVariantMatrixSchema returns the variant matrix schema configuration.
func (product *Product) GetVariantMatrixSchema() (VariantMatrixSchema, error) {
	schemaJSON := product.Get(COLUMN_VARIANT_MATRIX_SCHEMA)
	if schemaJSON == "" || schemaJSON == "null" {
		return VariantMatrixSchema{}, nil
	}
	var schema VariantMatrixSchema
	err := json.Unmarshal([]byte(schemaJSON), &schema)
	return schema, err
}

// HasVariantMatrixSchema returns true if a variant matrix schema is defined.
func (product *Product) HasVariantMatrixSchema() bool {
	dimJSON := product.Get(COLUMN_VARIANT_MATRIX_SCHEMA)
	return dimJSON != "" && dimJSON != "null"
}

// HasStock returns true if the product quantity is greater than 0.
func (product *Product) HasStock() bool {
	return product.GetQuantityInt() > 0
}

// IsOutOfStock returns true if the product quantity is less than or equal to 0.
func (product *Product) IsOutOfStock() bool {
	return product.GetQuantityInt() <= 0
}

// IsPaid returns true if the product price is greater than 0.
func (product *Product) IsPaid() bool {
	return product.GetPriceFloat() > 0
}

// IsFree returns true if the product price is less than or equal to 0.
func (product *Product) IsFree() bool {
	return product.GetPriceFloat() <= 0
}

// Slug returns the URL-friendly slug generated from the product title.
func (product *Product) Slug() string {
	title := product.GetTitle()
	return str.Slugify(title, '-')
}

// == GETTERS & SETTERS ========================================================

// GetCreatedAt returns the creation timestamp as a string.
func (product *Product) GetCreatedAt() string {
	return product.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
func (product *Product) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(product.GetCreatedAt(), carbon.UTC)
}

// SetCreatedAt sets the creation timestamp.
func (product *Product) SetCreatedAt(createdAt string) ProductInterface {
	product.Set(COLUMN_CREATED_AT, createdAt)
	return product
}

// GetDescription returns the full product description.
func (product *Product) GetDescription() string {
	return product.Get(COLUMN_DESCRIPTION)
}

// SetDescription sets the full product description.
func (product *Product) SetDescription(description string) ProductInterface {
	product.Set(COLUMN_DESCRIPTION, description)
	return product
}

// GetID returns the unique identifier.
func (product *Product) GetID() string {
	return product.Get(COLUMN_ID)
}

// SetID sets the unique identifier.
func (product *Product) SetID(id string) ProductInterface {
	product.Set(COLUMN_ID, id)
	return product
}

// GetMemo returns the internal memo.
func (product *Product) GetMemo() string {
	return product.Get(COLUMN_MEMO)
}

// SetMemo sets the internal memo.
func (product *Product) SetMemo(memo string) ProductInterface {
	product.Set(COLUMN_MEMO, memo)
	return product
}

// GetMetas returns all metadata as a map. Returns empty map if no metas stored.
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

// GetMeta returns a specific metadata value by name. Returns empty string if not found.
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

// SetMeta sets a single metadata value.
func (product *Product) SetMeta(name string, value string) error {
	return product.MetasUpsert(map[string]string{name: value})
}

// SetMetas replaces all metadata with the provided map.
// Warning: this overwrites any existing metadata.
func (product *Product) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	product.Set(COLUMN_METAS, string(mapString))
	return nil
}

// MetasUpsert merges the provided metadata with existing values.
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

// MetaRemove removes a single metadata entry.
func (product *Product) MetaRemove(name string) error {
	metas, err := product.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return product.SetMetas(metas)
}

// MetasRemove removes multiple metadata entries.
func (product *Product) MetasRemove(names []string) error {
	for _, name := range names {
		if err := product.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

// GetPrice returns the price as a string.
func (product *Product) GetPrice() string {
	return product.Get(COLUMN_PRICE)
}

// SetPrice sets the price from a string.
func (product *Product) SetPrice(price string) ProductInterface {
	product.Set(COLUMN_PRICE, price)
	return product
}

// GetPriceFloat returns the price as a float64.
func (product *Product) GetPriceFloat() float64 {
	price := cast.ToFloat64(product.Get(COLUMN_PRICE))
	return price
}

// SetPriceFloat sets the price from a float64.
func (product *Product) SetPriceFloat(price float64) ProductInterface {
	product.SetPrice(cast.ToString(price))
	return product
}

// GetQuantity returns the stock quantity as a string.
func (product *Product) GetQuantity() string {
	return product.Get(COLUMN_QUANTITY)
}

// SetQuantity sets the stock quantity from a string.
func (product *Product) SetQuantity(quantity string) ProductInterface {
	product.Set(COLUMN_QUANTITY, quantity)
	return product
}

// GetQuantityInt returns the stock quantity as an int64.
func (product *Product) GetQuantityInt() int64 {
	quantity := cast.ToInt64(product.GetQuantity())
	return quantity
}

// SetQuantityInt sets the stock quantity from an int64.
func (product *Product) SetQuantityInt(quantity int64) ProductInterface {
	product.SetQuantity(cast.ToString(quantity))
	return product
}

// GetShortDescription returns the short/abbreviated description.
func (product *Product) GetShortDescription() string {
	return product.Get(COLUMN_SHORT_DESCRIPTION)
}

// SetShortDescription sets the short/abbreviated description.
func (product *Product) SetShortDescription(shortDescription string) ProductInterface {
	product.Set(COLUMN_SHORT_DESCRIPTION, shortDescription)
	return product
}

// GetSoftDeletedAt returns the soft deletion timestamp.
func (product *Product) GetSoftDeletedAt() string {
	return product.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
func (product *Product) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(product.GetSoftDeletedAt(), carbon.UTC)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (product *Product) SetSoftDeletedAt(deletedAt string) ProductInterface {
	product.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return product
}

// GetStatus returns the current status.
func (product *Product) GetStatus() string {
	return product.Get(COLUMN_STATUS)
}

// SetStatus sets the current status.
func (product *Product) SetStatus(status string) ProductInterface {
	product.Set(COLUMN_STATUS, status)
	return product
}

// GetTitle returns the product title.
func (product *Product) GetTitle() string {
	return product.Get(COLUMN_TITLE)
}

// SetTitle sets the product title.
func (product *Product) SetTitle(title string) ProductInterface {
	product.Set(COLUMN_TITLE, title)
	return product
}

// GetUpdatedAt returns the last update timestamp.
func (product *Product) GetUpdatedAt() string {
	return product.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
func (product *Product) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(product.GetUpdatedAt(), carbon.UTC)
}

// SetUpdatedAt sets the last update timestamp.
func (product *Product) SetUpdatedAt(updatedAt string) ProductInterface {
	product.Set(COLUMN_UPDATED_AT, updatedAt)
	return product
}
