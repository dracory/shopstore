package shopstore

import (
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
)

// == CLASS ====================================================================

type Order struct {
	dataobject.DataObject
}

var _ OrderInterface = (*Order)(nil)

// == CONSTRUCTORS =============================================================

func NewOrder() OrderInterface {
	o := (&Order{}).
		SetID(uid.HumanUid()).
		SetStatus(ORDER_STATUS_PENDING).
		SetQuantityInt(1). // By default 1
		SetPriceFloat(0).  // Free. By default
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	_ = o.SetMetas(map[string]string{})

	return o
}

func NewOrderFromExistingData(data map[string]string) *Order {
	o := &Order{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

func (order *Order) IsAwaitingFulfillment() bool {
	return order.Status() == ORDER_STATUS_AWAITING_FULFILLMENT
}

func (order *Order) IsAwaitingPayment() bool {
	return order.Status() == ORDER_STATUS_AWAITING_PAYMENT
}

func (order *Order) IsAwaitingPickup() bool {
	return order.Status() == ORDER_STATUS_AWAITING_PICKUP
}

func (order *Order) IsAwaitingShipment() bool {
	return order.Status() == ORDER_STATUS_AWAITING_SHIPMENT
}

func (order *Order) IsCancelled() bool {
	return order.Status() == ORDER_STATUS_CANCELLED
}

func (order *Order) IsCompleted() bool {
	return order.Status() == ORDER_STATUS_COMPLETED
}

func (order *Order) IsDeclined() bool {
	return order.Status() == ORDER_STATUS_DECLINED
}

func (order *Order) IsDisputed() bool {
	return order.Status() == ORDER_STATUS_DISPUTED
}

func (order *Order) IsManualVerificationRequired() bool {
	return order.Status() == ORDER_STATUS_MANUAL_VERIFICATION_REQUIRED
}

func (order *Order) IsPending() bool {
	return order.Status() == ORDER_STATUS_PENDING
}

func (order *Order) IsRefunded() bool {
	return order.Status() == ORDER_STATUS_REFUNDED
}

func (order *Order) IsShipped() bool {
	return order.Status() == ORDER_STATUS_SHIPPED
}

// == GETTERS & SETTERS ========================================================

func (order *Order) CreatedAt() string {
	return order.Get(COLUMN_CREATED_AT)
}

func (order *Order) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(order.CreatedAt(), carbon.UTC)
}

func (order *Order) SetCreatedAt(createdAt string) OrderInterface {
	order.Set(COLUMN_CREATED_AT, createdAt)
	return order
}

func (order *Order) CustomerID() string {
	return order.Get(COLUMN_CUSTOMER_ID)
}

func (order *Order) SetCustomerID(id string) OrderInterface {
	order.Set(COLUMN_CUSTOMER_ID, id)
	return order
}

func (order *Order) ID() string {
	return order.Get(COLUMN_ID)
}

func (order *Order) SetID(id string) OrderInterface {
	order.Set(COLUMN_ID, id)
	return order
}

func (order *Order) Memo() string {
	return order.Get(COLUMN_MEMO)
}

func (order *Order) SetMemo(memo string) OrderInterface {
	order.Set(COLUMN_MEMO, memo)
	return order
}

func (order *Order) Metas() (map[string]string, error) {
	metasStr := order.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
}

func (order *Order) Meta(name string) string {
	metas, err := order.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (order *Order) SetMeta(name string, value string) error {
	return order.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (order *Order) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}
	order.Set(COLUMN_METAS, mapString)
	return nil
}

func (order *Order) UpsertMetas(metas map[string]string) error {
	currentMetas, err := order.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return order.SetMetas(currentMetas)
}

func (order *Order) Status() string {
	return order.Get(COLUMN_STATUS)
}

func (order *Order) SetStatus(status string) OrderInterface {
	order.Set(COLUMN_STATUS, status)
	return order
}

func (order *Order) Price() string {
	return order.Get(COLUMN_PRICE)
}

func (order *Order) SetPrice(price string) OrderInterface {
	order.Set(COLUMN_PRICE, price)
	return order
}

func (order *Order) PriceFloat() float64 {
	price, _ := utils.ToFloat(order.Get(COLUMN_PRICE))
	return price
}

func (order *Order) SetPriceFloat(price float64) OrderInterface {
	order.SetPrice(utils.ToString(price))
	return order
}

func (order *Order) Quantity() string {
	return order.Get(COLUMN_QUANTITY)
}

func (order *Order) SetQuantity(quantity string) OrderInterface {
	order.Set(COLUMN_QUANTITY, quantity)
	return order
}

func (order *Order) QuantityInt() int64 {
	quantity, _ := utils.ToInt(order.Quantity())
	return quantity
}

func (order *Order) SetQuantityInt(quantity int64) OrderInterface {
	order.SetQuantity(utils.ToString(quantity))
	return order
}

func (order *Order) SoftDeletedAt() string {
	return order.Get(COLUMN_SOFT_DELETED_AT)
}

func (order *Order) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(order.SoftDeletedAt(), carbon.UTC)
}

func (order *Order) SetSoftDeletedAt(deletedAt string) OrderInterface {
	order.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return order
}

func (order *Order) UpdatedAt() string {
	return order.Get(COLUMN_UPDATED_AT)
}

func (order *Order) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(order.UpdatedAt(), carbon.UTC)
}

func (order *Order) SetUpdatedAt(updatedAt string) OrderInterface {
	order.Set(COLUMN_UPDATED_AT, updatedAt)
	return order
}
