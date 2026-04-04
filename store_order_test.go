package shopstore

import (
	"context"
	"strings"
	"testing"

	"github.com/dracory/sb"
)

func TestStoreOderCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = order.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.OrderCreate(ctx, order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreOderDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()
	err = store.OrderCreate(ctx, order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDelete(ctx, order)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderFound, err := store.OrderFindByID(ctx, order.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if orderFound != nil {
		t.Fatal("expected nil order")
	}
}

func TestStoreOderDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()
	err = store.OrderCreate(ctx, order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDeleteByID(ctx, order.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderFound, err := store.OrderFindByID(ctx, order.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if orderFound != nil {
		t.Fatal("expected nil order")
	}
}

func TestStoreOrderFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99).
		SetMemo("test memo")

	err = order.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.OrderCreate(ctx, order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderFound, errFind := store.OrderFindByID(ctx, order.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderFound == nil {
		t.Fatal("Order MUST NOT be nil")
	}

	if orderFound.GetCustomerID() != "CUSTOMER01_ID" {
		t.Fatal("Order user id MUST BE 'CUSTOMER01_ID', found: ", orderFound.GetCustomerID())
	}

	if orderFound.GetStatus() != ORDER_STATUS_PENDING {
		t.Fatal("Order status MUST BE 'pending', found: ", orderFound.GetStatus())
	}

	if orderFound.GetQuantity() != "1" {
		t.Fatal("Order quantity MUST BE '1', found: ", orderFound.GetQuantity())
	}

	if orderFound.GetPrice() != "19.99" {
		t.Fatal("Order price MUST BE '19.99', found: ", orderFound.GetPrice())
	}

	if orderFound.GetMemo() != "test memo" {
		t.Fatal("Order memo MUST BE 'test memo', found: ", orderFound.GetMemo())
	}

	if orderFound.GetMeta("color") != "green" {
		t.Fatal("Order color meta MUST BE 'green', found: ", orderFound.GetMeta("color"))
	}

	if orderFound.GetMeta("size") != "xxl" {
		t.Fatal("Order size meta MUST BE 'xxl', found: ", orderFound.GetMeta("xxl"))
	}

	if !strings.Contains(orderFound.GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Order MUST NOT be soft deleted", orderFound.GetSoftDeletedAt())
	}
}

func TestStoreOrderSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("USER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()
	err = store.OrderCreate(ctx, order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderSoftDeleteByID(ctx, order.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if order.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Order MUST NOT be soft deleted")
	}

	orderFound, errFind := store.OrderFindByID(ctx, order.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderFound != nil {
		t.Fatal("Order MUST be nil")
	}

	orderFindWithDeleted, errFind := store.OrderList(ctx, NewOrderQuery().
		SetID(order.GetID()).
		SetSoftDeletedIncluded(true))

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(orderFindWithDeleted) < 1 {
		t.Fatal("Order list MUST NOT be empty")
		return
	}

	if strings.Contains(orderFindWithDeleted[0].GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Order MUST be soft deleted", orderFindWithDeleted[0].GetSoftDeletedAt())
	}

}

func TestStoreOrderUpdate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()
	err = store.OrderCreate(ctx, order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	order.SetMemo("test memo")

	err = store.OrderUpdate(ctx, order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderFound, errFind := store.OrderFindByID(ctx, order.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderFound == nil {
		t.Fatal("Order MUST NOT be nil")
	}

	if orderFound.GetMemo() != "test memo" {
		t.Fatal("Order memo MUST BE 'test memo', found: ", orderFound.GetMemo())
	}
}

func TestStoreOrderLineItemCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()
	err = store.OrderLineItemCreate(ctx, orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreOrderLineItemDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()
	err = store.OrderLineItemCreate(ctx, orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderLineItemDelete(ctx, orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderLineItemFound, errFind := store.OrderLineItemFindByID(ctx, orderLineItem.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderLineItemFound != nil {
		t.Fatal("OrderLineItem MUST be nil")
	}
}

func TestStoreOrderLineItemFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()
	err = store.OrderLineItemCreate(ctx, orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderLineItemFound, errFind := store.OrderLineItemFindByID(ctx, orderLineItem.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderLineItemFound == nil {
		t.Fatal("OrderLineItem MUST NOT be nil")
	}
}

func TestStoreOrderLineItemList(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	orderLineItem2 := NewOrderLineItem().
		SetOrderID("ORDER02_ID").
		SetProductID("PRODUCT02_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()

	err = store.OrderLineItemCreate(ctx, orderLineItem)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderLineItemCreate(ctx, orderLineItem2)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderLineItemsFound, errFind := store.OrderLineItemList(ctx, NewOrderLineItemQuery().
		SetOrderID(orderLineItem.GetOrderID()))

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(orderLineItemsFound) != 1 {
		t.Fatal("OrderLineItem MUST NOT be nil")
	}
}

func TestStoreOrderLineItemSoftDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	ctx := context.Background()
	err = store.OrderLineItemCreate(ctx, orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderLineItemSoftDeleteByID(ctx, orderLineItem.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderLineItemFound, errFind := store.OrderLineItemFindByID(ctx, orderLineItem.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderLineItemFound != nil {
		t.Fatal("OrderLineItem MUST be nil")
	}

	orderLineItems, errFind := store.OrderLineItemList(ctx, NewOrderLineItemQuery().
		SetOrderID(orderLineItem.GetOrderID()).
		SetSoftDeletedIncluded(true))

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(orderLineItems) != 1 {
		t.Fatal("OrderLineItem MUST be deleted")
	}
}
