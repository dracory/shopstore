package shopstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == CLASS ====================================================================

type OrderLineItem struct {
	dataobject.DataObject
}

var _ OrderLineItemInterface = (*OrderLineItem)(nil)

// == CONSTRUCTORS =============================================================

func NewOrderLineItem() OrderLineItemInterface {
	o := (&OrderLineItem{}).
		SetID(uid.HumanUid()).
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

func NewOrderLineItemFromExistingData(data map[string]string) OrderLineItemInterface {
	o := &OrderLineItem{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

func (o *OrderLineItem) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *OrderLineItem) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *OrderLineItem) SetCreatedAt(createdAt string) OrderLineItemInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *OrderLineItem) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *OrderLineItem) SetID(id string) OrderLineItemInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *OrderLineItem) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *OrderLineItem) SetMemo(memo string) OrderLineItemInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *OrderLineItem) Metas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

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

func (o *OrderLineItem) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *OrderLineItem) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *OrderLineItem) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}
	o.Set(COLUMN_METAS, string(mapString))
	return nil
}

func (o *OrderLineItem) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *OrderLineItem) OrderID() string {
	return o.Get(COLUMN_ORDER_ID)
}

func (o *OrderLineItem) SetOrderID(orderID string) OrderLineItemInterface {
	o.Set(COLUMN_ORDER_ID, orderID)
	return o
}

func (o *OrderLineItem) Price() string {
	return o.Get(COLUMN_PRICE)
}

func (o *OrderLineItem) SetPrice(price string) OrderLineItemInterface {
	o.Set(COLUMN_PRICE, price)
	return o
}

func (o *OrderLineItem) PriceFloat() float64 {
	price := o.Price()
	priceFloat := cast.ToFloat64(price)
	return priceFloat
}

func (o *OrderLineItem) SetPriceFloat(price float64) OrderLineItemInterface {
	o.SetPrice(cast.ToString(price))
	return o
}

func (o *OrderLineItem) ProductID() string {
	return o.Get(COLUMN_PRODUCT_ID)
}

func (o *OrderLineItem) SetProductID(productID string) OrderLineItemInterface {
	o.Set(COLUMN_PRODUCT_ID, productID)
	return o
}

func (o *OrderLineItem) Quantity() string {
	return o.Get(COLUMN_QUANTITY)
}

func (o *OrderLineItem) SetQuantity(quantity string) OrderLineItemInterface {
	o.Set(COLUMN_QUANTITY, quantity)
	return o
}

func (o *OrderLineItem) QuantityInt() int64 {
	quantity := o.Quantity()
	quantityInt := cast.ToInt64(quantity)
	return quantityInt
}

func (o *OrderLineItem) SetQuantityInt(quantity int64) OrderLineItemInterface {
	o.SetQuantity(cast.ToString(quantity))
	return o
}

func (o *OrderLineItem) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *OrderLineItem) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *OrderLineItem) SetSoftDeletedAt(deletedAt string) OrderLineItemInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return o
}

func (o *OrderLineItem) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *OrderLineItem) SetStatus(status string) OrderLineItemInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *OrderLineItem) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *OrderLineItem) SetTitle(title string) OrderLineItemInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *OrderLineItem) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *OrderLineItem) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

func (o *OrderLineItem) SetUpdatedAt(updatedAt string) OrderLineItemInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
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
