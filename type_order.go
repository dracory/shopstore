package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == CLASS ====================================================================

// Order represents a customer order in the shop store.
// Orders track customer purchases with status workflow, pricing, quantity management,
// soft deletion, and metadata storage. Supports various order states from pending to completed.
type Order struct {
	dataobject.DataObject
}

// == INTERFACES ===============================================================

// Compile-time interface compliance check
var _ OrderInterface = (*Order)(nil)

// == CONSTRUCTORS =============================================================

// NewOrder creates a new order with default values:
// - Status: pending
// - Quantity: 1
// - Price: 0.00 (free)
// - Memo: empty
// - CreatedAt: current UTC time
// - UpdatedAt: current UTC time
// - SoftDeletedAt: max datetime (not deleted)
// - Metas: empty map
func NewOrder() OrderInterface {
	o := (&Order{}).
		SetID(GenerateShortID()).
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

// NewOrderFromExistingData creates an order from existing data map.
// Used when hydrating from database or external sources.
func NewOrderFromExistingData(data map[string]string) OrderInterface {
	o := &Order{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

// IsAwaitingFulfillment returns true if order status is awaiting fulfillment.
func (order *Order) IsAwaitingFulfillment() bool {
	return order.GetStatus() == ORDER_STATUS_AWAITING_FULFILLMENT
}

// IsAwaitingPayment returns true if order status is awaiting payment.
func (order *Order) IsAwaitingPayment() bool {
	return order.GetStatus() == ORDER_STATUS_AWAITING_PAYMENT
}

// IsAwaitingPickup returns true if order status is awaiting pickup.
func (order *Order) IsAwaitingPickup() bool {
	return order.GetStatus() == ORDER_STATUS_AWAITING_PICKUP
}

// IsAwaitingShipment returns true if order status is awaiting shipment.
func (order *Order) IsAwaitingShipment() bool {
	return order.GetStatus() == ORDER_STATUS_AWAITING_SHIPMENT
}

// IsCancelled returns true if order status is cancelled.
func (order *Order) IsCancelled() bool {
	return order.GetStatus() == ORDER_STATUS_CANCELLED
}

// IsCompleted returns true if order status is completed.
func (order *Order) IsCompleted() bool {
	return order.GetStatus() == ORDER_STATUS_COMPLETED
}

// IsDeclined returns true if order status is declined.
func (order *Order) IsDeclined() bool {
	return order.GetStatus() == ORDER_STATUS_DECLINED
}

// IsDisputed returns true if order status is disputed.
func (order *Order) IsDisputed() bool {
	return order.GetStatus() == ORDER_STATUS_DISPUTED
}

// IsManualVerificationRequired returns true if order requires manual verification.
func (order *Order) IsManualVerificationRequired() bool {
	return order.GetStatus() == ORDER_STATUS_MANUAL_VERIFICATION_REQUIRED
}

// IsPending returns true if order status is pending.
func (order *Order) IsPending() bool {
	return order.GetStatus() == ORDER_STATUS_PENDING
}

// IsRefunded returns true if order status is refunded.
func (order *Order) IsRefunded() bool {
	return order.GetStatus() == ORDER_STATUS_REFUNDED
}

// IsShipped returns true if order status is shipped.
func (order *Order) IsShipped() bool {
	return order.GetStatus() == ORDER_STATUS_SHIPPED
}

// == GETTERS & SETTERS ========================================================

// GetCreatedAt returns the creation timestamp as a string.
func (order *Order) GetCreatedAt() string {
	return order.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
func (order *Order) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(order.GetCreatedAt(), carbon.UTC)
}

// SetCreatedAt sets the creation timestamp.
func (order *Order) SetCreatedAt(createdAt string) OrderInterface {
	order.Set(COLUMN_CREATED_AT, createdAt)
	return order
}

// GetCustomerID returns the customer ID.
func (order *Order) GetCustomerID() string {
	return order.Get(COLUMN_CUSTOMER_ID)
}

// SetCustomerID sets the customer ID.
func (order *Order) SetCustomerID(id string) OrderInterface {
	order.Set(COLUMN_CUSTOMER_ID, id)
	return order
}

// GetID returns the unique identifier.
func (order *Order) GetID() string {
	return order.Get(COLUMN_ID)
}

// SetID sets the unique identifier.
func (order *Order) SetID(id string) OrderInterface {
	order.Set(COLUMN_ID, id)
	return order
}

// GetMemo returns the internal memo.
func (order *Order) GetMemo() string {
	return order.Get(COLUMN_MEMO)
}

// SetMemo sets the internal memo.
func (order *Order) SetMemo(memo string) OrderInterface {
	order.Set(COLUMN_MEMO, memo)
	return order
}

// GetMeta returns a specific metadata value by name. Returns empty string if not found.
func (order *Order) GetMeta(name string) string {
	metas, err := order.GetMetas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// SetMeta sets a single metadata value.
func (order *Order) SetMeta(name string, value string) error {
	return order.MetasUpsert(map[string]string{name: value})
}

// GetMetas returns all metadata as a map. Returns empty map if no metas stored.
func (order *Order) GetMetas() (map[string]string, error) {
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

// SetMetas replaces all metadata with the provided map.
// Warning: this overwrites any existing metadata.
func (order *Order) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	order.Set(COLUMN_METAS, string(mapString))
	return nil
}

// MetasUpsert merges the provided metadata with existing values.
func (order *Order) MetasUpsert(metas map[string]string) error {
	currentMetas, err := order.GetMetas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return order.SetMetas(currentMetas)
}

// MetaRemove removes a single metadata entry.
func (order *Order) MetaRemove(name string) error {
	metas, err := order.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return order.SetMetas(metas)
}

// MetasRemove removes multiple metadata entries.
func (order *Order) MetasRemove(names []string) error {
	for _, name := range names {
		if err := order.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

// GetStatus returns the current status.
func (order *Order) GetStatus() string {
	return order.Get(COLUMN_STATUS)
}

// SetStatus sets the current status.
func (order *Order) SetStatus(status string) OrderInterface {
	order.Set(COLUMN_STATUS, status)
	return order
}

// GetPrice returns the price as a string.
func (order *Order) GetPrice() string {
	return order.Get(COLUMN_PRICE)
}

// SetPrice sets the price from a string.
func (order *Order) SetPrice(price string) OrderInterface {
	order.Set(COLUMN_PRICE, price)
	return order
}

// GetPriceFloat returns the price as a float64.
func (order *Order) GetPriceFloat() float64 {
	price := cast.ToFloat64(order.Get(COLUMN_PRICE))
	return price
}

// SetPriceFloat sets the price from a float64.
func (order *Order) SetPriceFloat(price float64) OrderInterface {
	order.SetPrice(cast.ToString(price))
	return order
}

// GetQuantity returns the quantity as a string.
func (order *Order) GetQuantity() string {
	return order.Get(COLUMN_QUANTITY)
}

// SetQuantity sets the quantity from a string.
func (order *Order) SetQuantity(quantity string) OrderInterface {
	order.Set(COLUMN_QUANTITY, quantity)
	return order
}

// GetQuantityInt returns the quantity as an int64.
func (order *Order) GetQuantityInt() int64 {
	quantity := cast.ToInt64(order.GetQuantity())
	return quantity
}

// SetQuantityInt sets the quantity from an int64.
func (order *Order) SetQuantityInt(quantity int64) OrderInterface {
	order.SetQuantity(cast.ToString(quantity))
	return order
}

// GetSoftDeletedAt returns the soft deletion timestamp.
func (order *Order) GetSoftDeletedAt() string {
	return order.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
func (order *Order) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(order.GetSoftDeletedAt(), carbon.UTC)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (order *Order) SetSoftDeletedAt(deletedAt string) OrderInterface {
	order.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return order
}

// GetUpdatedAt returns the last update timestamp.
func (order *Order) GetUpdatedAt() string {
	return order.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
func (order *Order) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(order.GetUpdatedAt(), carbon.UTC)
}

// SetUpdatedAt sets the last update timestamp.
func (order *Order) SetUpdatedAt(updatedAt string) OrderInterface {
	order.Set(COLUMN_UPDATED_AT, updatedAt)
	return order
}
