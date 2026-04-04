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

type Discount struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

var _ DiscountInterface = (*Discount)(nil)

// == CONSTRUCTORS ===========================================================

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

func generateDiscountCode() string {
	code, err := str.RandomFromGamma(12, "BCDFGHJKLMNPQRSTVWXYZ23456789")

	if err != nil {
		code = str.Random(12)
	}

	return code
}

func NewDiscountFromExistingData(data map[string]string) DiscountInterface {
	o := &Discount{}
	o.Hydrate(data)
	return o
}

// == METHODS ================================================================

// == SETTERS AND GETTERS ====================================================

func (d *Discount) GetAmount() float64 {
	amountStr := d.Get(COLUMN_AMOUNT)
	amount := cast.ToFloat64(amountStr)

	return amount
}

func (d *Discount) SetAmount(amount float64) DiscountInterface {
	amountStr := cast.ToString(amount)
	d.Set(COLUMN_AMOUNT, amountStr)
	return d
}

func (d *Discount) GetCode() string {
	return d.Get(COLUMN_CODE)
}

func (d *Discount) SetCode(code string) DiscountInterface {
	d.Set(COLUMN_CODE, code)
	return d
}

func (d *Discount) GetCreatedAt() string {
	return d.Get(COLUMN_CREATED_AT)
}

func (d *Discount) GetCreatedAtCarbon() *carbon.Carbon {
	createdAt := d.GetCreatedAt()
	return carbon.Parse(createdAt)
}

func (d *Discount) SetCreatedAt(createdAt string) DiscountInterface {
	d.Set(COLUMN_CREATED_AT, createdAt)
	return d
}

func (d *Discount) GetDescription() string {
	return d.Get(COLUMN_DESCRIPTION)
}

func (d *Discount) SetDescription(description string) DiscountInterface {
	d.Set(COLUMN_DESCRIPTION, description)
	return d
}

func (d *Discount) GetEndsAt() string {
	return d.Get(COLUMN_ENDS_AT)
}

func (d *Discount) GetEndsAtCarbon() *carbon.Carbon {
	endsAt := d.GetEndsAt()
	return carbon.Parse(endsAt)
}

func (d *Discount) SetEndsAt(endsAt string) DiscountInterface {
	d.Set(COLUMN_ENDS_AT, endsAt)
	return d
}

// GetID returns the ID of the discount
func (d *Discount) GetID() string {
	return d.Get(COLUMN_ID)
}

// SetID sets the ID of the discount
func (d *Discount) SetID(id string) DiscountInterface {
	d.Set(COLUMN_ID, id)
	return d
}

func (d *Discount) GetMemo() string {
	return d.Get(COLUMN_MEMO)
}

func (d *Discount) SetMemo(memo string) DiscountInterface {
	d.Set(COLUMN_MEMO, memo)
	return d
}

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

func (d *Discount) MetaRemove(name string) error {
	metas, err := d.GetMetas()

	if err != nil {
		return err
	}

	delete(metas, name)

	return d.SetMetas(metas)
}

func (d *Discount) SetMeta(name string, value string) error {
	return d.MetasUpsert(map[string]string{name: value})
}

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

func (d *Discount) MetasRemove(names []string) error {
	for _, name := range names {
		err := d.MetaRemove(name)

		if err != nil {
			return err
		}
	}

	return nil
}

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

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (d *Discount) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)

	if err != nil {
		return err
	}

	d.Set(COLUMN_METAS, string(mapString))

	return nil
}

func (d *Discount) GetSoftDeletedAt() string {
	return d.Get(COLUMN_SOFT_DELETED_AT)
}

func (d *Discount) GetSoftDeletedAtCarbon() *carbon.Carbon {
	deletedAt := d.GetSoftDeletedAt()
	return carbon.Parse(deletedAt)
}

func (d *Discount) SetSoftDeletedAt(deletedAt string) DiscountInterface {
	d.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return d
}

func (d *Discount) GetStartsAt() string {
	return d.Get(COLUMN_STARTS_AT)
}

func (d *Discount) GetStartsAtCarbon() *carbon.Carbon {
	startsAt := d.GetStartsAt()
	return carbon.Parse(startsAt)
}

func (d *Discount) SetStartsAt(startsAt string) DiscountInterface {
	d.Set(COLUMN_STARTS_AT, startsAt)
	return d
}

func (d *Discount) GetStatus() string {
	return d.Get(COLUMN_STATUS)
}

func (d *Discount) SetStatus(status string) DiscountInterface {
	d.Set(COLUMN_STATUS, status)
	return d
}

func (d *Discount) GetTitle() string {
	return d.Get(COLUMN_TITLE)
}

func (d *Discount) SetTitle(title string) DiscountInterface {
	d.Set(COLUMN_TITLE, title)
	return d
}

func (d *Discount) GetType() string {
	return d.Get(COLUMN_TYPE)
}

func (d *Discount) SetType(type_ string) DiscountInterface {
	d.Set(COLUMN_TYPE, type_)
	return d
}

func (d *Discount) GetUpdatedAt() string {
	return d.Get(COLUMN_UPDATED_AT)
}

func (d *Discount) GetUpdatedAtCarbon() *carbon.Carbon {
	updatedAt := d.GetUpdatedAt()
	return carbon.Parse(updatedAt)
}

func (d *Discount) SetUpdatedAt(updatedAt string) DiscountInterface {
	d.Set(COLUMN_UPDATED_AT, updatedAt)
	return d
}

// IsActive returns true if the discount status is active
func (d *Discount) IsActive() bool {
	return d.GetStatus() == DISCOUNT_STATUS_ACTIVE
}

// IsDraft returns true if the discount status is draft
func (d *Discount) IsDraft() bool {
	return d.GetStatus() == DISCOUNT_STATUS_DRAFT
}

// IsInactive returns true if the discount status is inactive
func (d *Discount) IsInactive() bool {
	return d.GetStatus() == DISCOUNT_STATUS_INACTIVE
}

// IsStarted returns true if the discount has started (starts_at <= now)
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

// IsEnded returns true if the discount has ended (ends_at <= now)
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

// IsExpired is an alias for IsEnded
func (d *Discount) IsExpired() bool {
	return d.IsEnded()
}

// IsValidNow returns true if the discount is active, started, and not ended
func (d *Discount) IsValidNow() bool {
	return d.IsActive() && d.IsStarted() && !d.IsEnded()
}
