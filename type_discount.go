package shopstore

import (
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/utils"
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
	code := uid.Timestamp()

	d := (&Discount{}).
		SetID(uid.HumanUid()).
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

func NewDiscountFromExistingData(data map[string]string) DiscountInterface {
	o := &Discount{}
	o.Hydrate(data)
	return o
}

// == METHODS ================================================================

// == SETTERS AND GETTERS ====================================================

func (d *Discount) Amount() float64 {
	amountStr := d.Get(COLUMN_AMOUNT)
	amount, err := utils.ToFloat(amountStr)

	if err != nil {
		return 0
	}

	return amount
}

func (d *Discount) SetAmount(amount float64) DiscountInterface {
	amountStr := utils.ToString(amount)
	d.Set(COLUMN_AMOUNT, amountStr)
	return d
}

func (d *Discount) Code() string {
	return d.Get(COLUMN_CODE)
}

func (d *Discount) SetCode(code string) DiscountInterface {
	d.Set(COLUMN_CODE, code)
	return d
}

func (d *Discount) CreatedAt() string {
	return d.Get(COLUMN_CREATED_AT)
}

func (d *Discount) CreatedAtCarbon() *carbon.Carbon {
	createdAt := d.CreatedAt()
	return carbon.Parse(createdAt)
}

func (d *Discount) SetCreatedAt(createdAt string) DiscountInterface {
	d.Set(COLUMN_CREATED_AT, createdAt)
	return d
}

func (d *Discount) Description() string {
	return d.Get(COLUMN_DESCRIPTION)
}

func (d *Discount) SetDescription(description string) DiscountInterface {
	d.Set(COLUMN_DESCRIPTION, description)
	return d
}

func (d *Discount) EndsAt() string {
	return d.Get(COLUMN_ENDS_AT)
}

func (d *Discount) EndsAtCarbon() *carbon.Carbon {
	endsAt := d.EndsAt()
	return carbon.Parse(endsAt)
}

func (d *Discount) SetEndsAt(endsAt string) DiscountInterface {
	d.Set(COLUMN_ENDS_AT, endsAt)
	return d
}

// ID returns the ID of the exam
func (o *Discount) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the ID of the exam
func (o *Discount) SetID(id string) DiscountInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (d *Discount) Memo() string {
	return d.Get(COLUMN_MEMO)
}

func (d *Discount) SetMemo(memo string) DiscountInterface {
	d.Set(COLUMN_MEMO, memo)
	return d
}

func (d *Discount) Meta(name string) string {
	metas, err := d.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (d *Discount) MetaRemove(name string) error {
	metas, err := d.Metas()

	if err != nil {
		return err
	}

	delete(metas, name)

	return d.SetMetas(metas)
}

func (d *Discount) SetMeta(name string, value string) error {
	return d.MetasUpsert(map[string]string{name: value})
}

func (d *Discount) Metas() (map[string]string, error) {
	metasStr := d.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
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
	currentMetas, err := d.Metas()

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
	mapString, err := utils.ToJSON(metas)

	if err != nil {
		return err
	}

	d.Set(COLUMN_METAS, mapString)

	return nil
}

func (d *Discount) SoftDeletedAt() string {
	return d.Get(COLUMN_SOFT_DELETED_AT)
}

func (d *Discount) SoftDeletedAtCarbon() *carbon.Carbon {
	deletedAt := d.SoftDeletedAt()
	return carbon.Parse(deletedAt)
}

func (d *Discount) SetSoftDeletedAt(deletedAt string) DiscountInterface {
	d.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return d
}

func (d *Discount) StartsAt() string {
	return d.Get(COLUMN_STARTS_AT)
}

func (d *Discount) StartsAtCarbon() *carbon.Carbon {
	startsAt := d.StartsAt()
	return carbon.Parse(startsAt)
}

func (d *Discount) SetStartsAt(startsAt string) DiscountInterface {
	d.Set(COLUMN_STARTS_AT, startsAt)
	return d
}

func (d *Discount) Status() string {
	return d.Get(COLUMN_STATUS)
}

func (d *Discount) SetStatus(status string) DiscountInterface {
	d.Set(COLUMN_STATUS, status)
	return d
}

func (d *Discount) Title() string {
	return d.Get(COLUMN_TITLE)
}

func (d *Discount) SetTitle(title string) DiscountInterface {
	d.Set(COLUMN_TITLE, title)
	return d
}

func (d *Discount) Type() string {
	return d.Get(COLUMN_TYPE)
}

func (d *Discount) SetType(type_ string) DiscountInterface {
	d.Set(COLUMN_TYPE, type_)
	return d
}

func (d *Discount) UpdatedAt() string {
	return d.Get(COLUMN_UPDATED_AT)
}

func (d *Discount) UpdatedAtCarbon() *carbon.Carbon {
	updatedAt := d.UpdatedAt()
	return carbon.Parse(updatedAt)
}

func (d *Discount) SetUpdatedAt(updatedAt string) DiscountInterface {
	d.Set(COLUMN_UPDATED_AT, updatedAt)
	return d
}
