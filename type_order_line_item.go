package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == CLASS ====================================================================

// OrderLineItem represents an individual product within an order.
// Line items track their own pricing, quantity, product association, and status.
// Supports soft deletion and metadata storage.
type OrderLineItem struct {
	dataobject.DataObject
}

// == INTERFACES ===============================================================

// Compile-time interface compliance check
var _ OrderLineItemInterface = (*OrderLineItem)(nil)

// == CONSTRUCTORS =============================================================

// NewOrderLineItem creates a new order line item with default values:
// - Status: pending
// - Title: empty
// - Quantity: 1
// - Price: 0.00 (free)
// - Memo: empty
// - CreatedAt: current UTC time
// - UpdatedAt: current UTC time
// - SoftDeletedAt: max datetime (not deleted)
// - Metas: empty map
func NewOrderLineItem() OrderLineItemInterface {
	o := (&OrderLineItem{}).
		SetID(GenerateShortID()).
		SetStatus(ORDER_STATUS_PENDING).
		SetTitle("").
		SetQuantityInt(1). // By default 1
		SetPriceFloat(0).  // Free. By default
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	_ = o.SetMetas(map[string]string{})

	return o
}

// NewOrderLineItemFromExistingData creates an order line item from existing data map.
// Used when hydrating from database or external sources.
func NewOrderLineItemFromExistingData(data map[string]string) OrderLineItemInterface {
	o := &OrderLineItem{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

// GetCreatedAt returns the creation timestamp as a string.
func (o *OrderLineItem) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
func (o *OrderLineItem) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}

// SetCreatedAt sets the creation timestamp.
func (o *OrderLineItem) SetCreatedAt(createdAt string) OrderLineItemInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// GetID returns the unique identifier.
func (o *OrderLineItem) GetID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the unique identifier.
func (o *OrderLineItem) SetID(id string) OrderLineItemInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// GetMemo returns the internal memo.
func (o *OrderLineItem) GetMemo() string {
	return o.Get(COLUMN_MEMO)
}

// SetMemo sets the internal memo.
func (o *OrderLineItem) SetMemo(memo string) OrderLineItemInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// GetMetas returns all metadata as a map. Returns empty map if no metas stored.
func (o *OrderLineItem) GetMetas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

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

// GetMeta returns a specific metadata value by name. Returns empty string if not found.
func (o *OrderLineItem) GetMeta(name string) string {
	metas, err := o.GetMetas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// SetMeta sets a single metadata value.
func (o *OrderLineItem) SetMeta(name string, value string) error {
	return o.MetasUpsert(map[string]string{name: value})
}

// SetMetas replaces all metadata with the provided map.
// Warning: this overwrites any existing metadata.
func (o *OrderLineItem) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	o.Set(COLUMN_METAS, string(mapString))
	return nil
}

// MetasUpsert merges the provided metadata with existing values.
func (o *OrderLineItem) MetasUpsert(metas map[string]string) error {
	currentMetas, err := o.GetMetas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

// MetaRemove removes a single metadata entry.
func (o *OrderLineItem) MetaRemove(name string) error {
	metas, err := o.GetMetas()
	if err != nil {
		return err
	}
	delete(metas, name)
	return o.SetMetas(metas)
}

// MetasRemove removes multiple metadata entries.
func (o *OrderLineItem) MetasRemove(names []string) error {
	for _, name := range names {
		if err := o.MetaRemove(name); err != nil {
			return err
		}
	}
	return nil
}

// GetOrderID returns the associated order ID.
func (o *OrderLineItem) GetOrderID() string {
	return o.Get(COLUMN_ORDER_ID)
}

// SetOrderID sets the associated order ID.
func (o *OrderLineItem) SetOrderID(orderID string) OrderLineItemInterface {
	o.Set(COLUMN_ORDER_ID, orderID)
	return o
}

// GetPrice returns the price as a string.
func (o *OrderLineItem) GetPrice() string {
	return o.Get(COLUMN_PRICE)
}

// SetPrice sets the price from a string.
func (o *OrderLineItem) SetPrice(price string) OrderLineItemInterface {
	o.Set(COLUMN_PRICE, price)
	return o
}

// GetPriceFloat returns the price as a float64.
func (o *OrderLineItem) GetPriceFloat() float64 {
	price := o.GetPrice()
	priceFloat := cast.ToFloat64(price)
	return priceFloat
}

// SetPriceFloat sets the price from a float64.
func (o *OrderLineItem) SetPriceFloat(price float64) OrderLineItemInterface {
	o.SetPrice(cast.ToString(price))
	return o
}

// GetProductID returns the associated product ID.
func (o *OrderLineItem) GetProductID() string {
	return o.Get(COLUMN_PRODUCT_ID)
}

// SetProductID sets the associated product ID.
func (o *OrderLineItem) SetProductID(productID string) OrderLineItemInterface {
	o.Set(COLUMN_PRODUCT_ID, productID)
	return o
}

// GetQuantity returns the quantity as a string.
func (o *OrderLineItem) GetQuantity() string {
	return o.Get(COLUMN_QUANTITY)
}

// SetQuantity sets the quantity from a string.
func (o *OrderLineItem) SetQuantity(quantity string) OrderLineItemInterface {
	o.Set(COLUMN_QUANTITY, quantity)
	return o
}

// GetQuantityInt returns the quantity as an int64.
func (o *OrderLineItem) GetQuantityInt() int64 {
	quantity := o.GetQuantity()
	quantityInt := cast.ToInt64(quantity)
	return quantityInt
}

// SetQuantityInt sets the quantity from an int64.
func (o *OrderLineItem) SetQuantityInt(quantity int64) OrderLineItemInterface {
	o.SetQuantity(cast.ToString(quantity))
	return o
}

// GetSoftDeletedAt returns the soft deletion timestamp.
func (o *OrderLineItem) GetSoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
func (o *OrderLineItem) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetSoftDeletedAt(), carbon.UTC)
}

// SetSoftDeletedAt sets the soft deletion timestamp.
func (o *OrderLineItem) SetSoftDeletedAt(deletedAt string) OrderLineItemInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return o
}

// GetStatus returns the current status.
func (o *OrderLineItem) GetStatus() string {
	return o.Get(COLUMN_STATUS)
}

// SetStatus sets the current status.
func (o *OrderLineItem) SetStatus(status string) OrderLineItemInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

// GetTitle returns the line item title.
func (o *OrderLineItem) GetTitle() string {
	return o.Get(COLUMN_TITLE)
}

// SetTitle sets the line item title.
func (o *OrderLineItem) SetTitle(title string) OrderLineItemInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

// GetUpdatedAt returns the last update timestamp.
func (o *OrderLineItem) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
func (o *OrderLineItem) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetUpdatedAt(), carbon.UTC)
}

// SetUpdatedAt sets the last update timestamp.
func (o *OrderLineItem) SetUpdatedAt(updatedAt string) OrderLineItemInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// IsActive returns true if the order line item is in an active state.
func (o *OrderLineItem) IsActive() bool {
	status := o.GetStatus()
	return status == ORDER_STATUS_AWAITING_FULFILLMENT ||
		status == ORDER_STATUS_AWAITING_PAYMENT ||
		status == ORDER_STATUS_AWAITING_PICKUP ||
		status == ORDER_STATUS_AWAITING_SHIPMENT ||
		status == ORDER_STATUS_PENDING ||
		status == ORDER_STATUS_PARTIALLY_SHIPPED ||
		status == ORDER_STATUS_SHIPPED
}

// IsCancelled returns true if the order line item is cancelled.
func (o *OrderLineItem) IsCancelled() bool {
	return o.GetStatus() == ORDER_STATUS_CANCELLED
}

// IsCompleted returns true if the order line item is completed.
func (o *OrderLineItem) IsCompleted() bool {
	return o.GetStatus() == ORDER_STATUS_COMPLETED
}

// IsDraft returns true if the order line item is in draft/pending state.
func (o *OrderLineItem) IsDraft() bool {
	return o.GetStatus() == ORDER_STATUS_PENDING
}

// HasQuantity returns true if the quantity is greater than 0.
func (o *OrderLineItem) HasQuantity() bool {
	return o.GetQuantityInt() > 0
}

// IsFree returns true if the price is less than or equal to 0.
func (o *OrderLineItem) IsFree() bool {
	return o.GetPriceFloat() <= 0
}

// type LineItem struct {
// 	ID       string
// 	OrdeID   string
// 	Name     string
// 	Price    float64
// 	Quantity int64
// }

// func (order *Order) LineItemAdd(lineItem LineItem) {
// 	order.lineItems = append(order.lineItems, lineItem)
// }

// func (order *Order) LineItemList() []LineItem {
// 	return order.lineItems
// }

// func (order *Order) LineItemRemove(lineItemId string) error {

// 	index := order.findLineItemIndex(lineItemId)
// 	if index != -1 {
// 		// As I'd like to keep the items ordered, in Golang we have to shift all of the elements at
// 		// the right of the one being deleted, by one to the left.
// 		order.lineItems = append(order.lineItems[:index], order.lineItems[index+1:]...)
// 	}

// 	return nil
// }

// func (order *Order) LineItemsRemoveAll() {
// 	order.lineItems = []LineItem{}
// }

// func (order *Order) findLineItemIndex(lineItemId string) int {
// 	for index, lineItem := range order.lineItems {
// 		if lineItemId == lineItem.ID {
// 			return index
// 		}
// 	}

// 	return -1
// }
