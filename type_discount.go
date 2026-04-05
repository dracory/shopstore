package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/str"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == CONSTANTS ==============================================================

const DISCOUNT_STATUS_DRAFT = "draft"
const DISCOUNT_STATUS_ACTIVE = "active"
const DISCOUNT_STATUS_INACTIVE = "inactive"

const DISCOUNT_TYPE_AMOUNT = "amount"
const DISCOUNT_TYPE_PERCENT = "percent"

const DISCOUNT_DURATION_FOREVER = "forever"
const DISCOUNT_DURATION_MONTHS = "months"
const DISCOUNT_DURATION_ONCE = "once"

// == CLASS ==================================================================

// Discount represents a discount/promotion in the shop store.
// Discounts support temporal validity (start/end dates), amount-based discounts,
// soft deletion, metadata storage, and status management.
type Discount struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// Compile-time interface compliance check
var _ DiscountInterface = (*Discount)(nil)

// == CONSTRUCTORS ===========================================================

// NewDiscount creates a new discount with default values:
// - Status: draft
// - Type: percent
// - Amount: 0.00
// - Code: randomly generated 12-character code
// - Title: empty
// - Description: empty
// - StartsAt: null datetime
// - EndsAt: null datetime
// - Memo: empty
// - CreatedAt: current UTC time
// - UpdatedAt: current UTC time
// - SoftDeletedAt: max datetime (not deleted)
// - Metas: empty map
func NewDiscount() DiscountInterface {
	code := generateDiscountCode()

	d := (&Discount{}).
		SetID(GenerateShortID()).
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetType(DISCOUNT_TYPE_PERCENT).
		SetTitle("").
		SetDescription("").
		SetAmount(0.00).
		SetCode(code).
		SetStartsAt(sb.NULL_DATETIME).
		SetEndsAt(sb.NULL_DATETIME).
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	d.SetMetas(map[string]string{})

	return d
}

// generateDiscountCode generates a random 12-character discount code
// using alphanumeric characters (excluding confusing characters like 0, O, 1, I).
func generateDiscountCode() string {
	code, err := str.RandomFromGamma(12, "BCDFGHJKLMNPQRSTVWXYZ23456789")

	if err != nil {
		code = str.Random(12)
	}

	return code
}

// NewDiscountFromExistingData creates a discount from existing data map.
// Used when hydrating from database or external sources.
func NewDiscountFromExistingData(data map[string]string) DiscountInterface {
	o := &Discount{}
	o.Hydrate(data)
	return o
}

// == METHODS ================================================================

// == SETTERS AND GETTERS ====================================================

// GetAmount returns the discount amount as a float64.
func (d *Discount) GetAmount() float64 {
	amountStr := d.Get(COLUMN_AMOUNT)
	amount := cast.ToFloat64(amountStr)

	return amount
}

// SetAmount sets the discount amount.
func (d *Discount) SetAmount(amount float64) DiscountInterface {
	amountStr := cast.ToString(amount)
	d.Set(COLUMN_AMOUNT, amountStr)
	return d
}

// GetCode returns the unique discount code.
func (d *Discount) GetCode() string {
	return d.Get(COLUMN_CODE)
}

// SetCode sets the unique discount code.
func (d *Discount) SetCode(code string) DiscountInterface {
	d.Set(COLUMN_CODE, code)
	return d
}

// GetCreatedAt returns the creation timestamp as a string.
func (d *Discount) GetCreatedAt() string {
	return d.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
func (d *Discount) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := d.GetCreatedAt()
	return carbon.Parse(createdAt)
}

// SetCreatedAt sets the creation timestamp.
func (d *Discount) SetCreatedAt(createdAt string) DiscountInterface {
	d.Set(COLUMN_CREATED_AT, createdAt)
	return d
}

// GetDescription returns the discount description.
func (d *Discount) GetDescription() string {
	return d.Get(COLUMN_DESCRIPTION)
}

// SetDescription sets the discount description.
func (d *Discount) SetDescription(description string) DiscountInterface {
	d.Set(COLUMN_DESCRIPTION, description)
	return d
}

// GetEndsAt returns the end date/time as a string.
func (d *Discount) GetEndsAt() string {
	return d.Get(COLUMN_ENDS_AT)
}

// GetEndsAtCarbon returns the end date/time as a Carbon instance.
func (d *Discount) GetEndsAtCarbon() *carbon.Carbon {
	endsAt := d.GetEndsAt()
	return carbon.Parse(endsAt)
}

// SetEndsAt sets the end date/time.
func (d *Discount) SetEndsAt(endsAt string) DiscountInterface {
	d.Set(COLUMN_ENDS_AT, endsAt)
	return d
}

// GetID returns the unique identifier.
func (d *Discount) GetID() string {
	return d.Get(COLUMN_ID)
}

// SetID sets the unique identifier.
func (d *Discount) SetID(id string) DiscountInterface {
	d.Set(COLUMN_ID, id)
	return d
}

// GetMemo returns the internal memo.
func (d *Discount) GetMemo() string {
	return d.Get(COLUMN_MEMO)
}

// SetMemo sets the internal memo.
func (d *Discount) SetMemo(memo string) DiscountInterface {
	d.Set(COLUMN_MEMO, memo)
	return d
}

// GetMeta returns a specific metadata value by name. Returns empty string if not found.
func (d *Discount) GetMeta(name string) string {
	metas, err := d.GetMetas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// MetaRemove removes a single metadata entry.
func (d *Discount) MetaRemove(name string) error {
	metas, err := d.GetMetas()

	if err != nil {
		return err
	}

	delete(metas, name)

	return d.SetMetas(metas)
}

// SetMeta sets a single metadata value.
func (d *Discount) SetMeta(name string, value string) error {
	return d.MetasUpsert(map[string]string{name: value})
}

// GetMetas returns all metadata as a map. Returns empty map if no metas stored.
func (d *Discount) GetMetas() (map[string]string, error) {
	metasStr := d.Get(COLUMN_METAS)

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

// MetasRemove removes multiple metadata entries.
func (d *Discount) MetasRemove(names []string) error {
	for _, name := range names {
		err := d.MetaRemove(name)

		if err != nil {
			return err
		}
	}

	return nil
}

// MetasUpsert merges the provided metadata with existing values.
func (d *Discount) MetasUpsert(metas map[string]string) error {
	currentMetas, err := d.GetMetas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return d.SetMetas(currentMetas)
}

// SetMetas replaces all metadata with the provided map.
// Warning: this overwrites any existing metadata.
func (d *Discount) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)

	if err != nil {
		return err
	}

	d.Set(COLUMN_METAS, string(mapString))

	return nil
}

// GetSoftDeletedAt returns the soft deletion timestamp.
func (d *Discount) GetSoftDeletedAt() string {
	return d.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
func (d *Discount) GetSoftDeletedAtCarbon() *carbon.Carbon {
	deletedAt := d.GetSoftDeletedAt()
	return carbon.Parse(deletedAt)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (d *Discount) SetSoftDeletedAt(deletedAt string) DiscountInterface {
	d.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return d
}

// GetStartsAt returns the start date/time as a string.
func (d *Discount) GetStartsAt() string {
	return d.Get(COLUMN_STARTS_AT)
}

// GetStartsAtCarbon returns the start date/time as a Carbon instance.
func (d *Discount) GetStartsAtCarbon() *carbon.Carbon {
	startsAt := d.GetStartsAt()
	return carbon.Parse(startsAt)
}

// SetStartsAt sets the start date/time.
func (d *Discount) SetStartsAt(startsAt string) DiscountInterface {
	d.Set(COLUMN_STARTS_AT, startsAt)
	return d
}

// GetStatus returns the current status.
func (d *Discount) GetStatus() string {
	return d.Get(COLUMN_STATUS)
}

// SetStatus sets the current status.
func (d *Discount) SetStatus(status string) DiscountInterface {
	d.Set(COLUMN_STATUS, status)
	return d
}

// GetTitle returns the discount title.
func (d *Discount) GetTitle() string {
	return d.Get(COLUMN_TITLE)
}

// SetTitle sets the discount title.
func (d *Discount) SetTitle(title string) DiscountInterface {
	d.Set(COLUMN_TITLE, title)
	return d
}

// GetType returns the discount type (amount or percent).
func (d *Discount) GetType() string {
	return d.Get(COLUMN_TYPE)
}

// SetType sets the discount type.
func (d *Discount) SetType(type_ string) DiscountInterface {
	d.Set(COLUMN_TYPE, type_)
	return d
}

// GetUpdatedAt returns the last update timestamp.
func (d *Discount) GetUpdatedAt() string {
	return d.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
func (d *Discount) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := d.GetUpdatedAt()
	return carbon.Parse(updatedAt)
}

// SetUpdatedAt sets the last update timestamp.
func (d *Discount) SetUpdatedAt(updatedAt string) DiscountInterface {
	d.Set(COLUMN_UPDATED_AT, updatedAt)
	return d
}

// IsActive returns true if the discount status is active.
func (d *Discount) IsActive() bool {
	return d.GetStatus() == DISCOUNT_STATUS_ACTIVE
}

// IsDraft returns true if the discount status is draft.
func (d *Discount) IsDraft() bool {
	return d.GetStatus() == DISCOUNT_STATUS_DRAFT
}

// IsInactive returns true if the discount status is inactive.
func (d *Discount) IsInactive() bool {
	return d.GetStatus() == DISCOUNT_STATUS_INACTIVE
}

// IsStarted returns true if the discount period has started (starts_at <= now).
func (d *Discount) IsStarted() bool {
	startsAt := d.GetStartsAt()
	if startsAt == sb.NULL_DATETIME || startsAt == "" {
		return false
	}
	startsAtCarbon := d.GetStartsAtCarbon()
	if startsAtCarbon == nil {
		return false
	}
	return !startsAtCarbon.IsFuture()
}

// IsEnded returns true if the discount period has ended (ends_at <= now).
func (d *Discount) IsEnded() bool {
	endsAt := d.GetEndsAt()
	if endsAt == sb.NULL_DATETIME || endsAt == "" {
		return false
	}
	endsAtCarbon := d.GetEndsAtCarbon()
	if endsAtCarbon == nil {
		return false
	}
	return endsAtCarbon.IsPast()
}

// IsExpired is an alias for IsEnded.
func (d *Discount) IsExpired() bool {
	return d.IsEnded()
}

// IsValidNow returns true if the discount is active, started, and not ended.
func (d *Discount) IsValidNow() bool {
	return d.IsActive() && d.IsStarted() && !d.IsEnded()
}
