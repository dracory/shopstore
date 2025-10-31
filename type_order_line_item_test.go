package shopstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewOrderLineItemDefaults(t *testing.T) {
	item := NewOrderLineItem()
	if item == nil {
		t.Fatal("NewOrderLineItem returned nil")
	}

	if item.Status() != ORDER_STATUS_PENDING {
		t.Fatalf("expected status %q, got %q", ORDER_STATUS_PENDING, item.Status())
	}

	if item.Title() != "" {
		t.Fatalf("expected empty title, got %q", item.Title())
	}

	if item.Memo() != "" {
		t.Fatalf("expected empty memo, got %q", item.Memo())
	}

	if item.PriceFloat() != 0 {
		t.Fatalf("expected default price 0, got %f", item.PriceFloat())
	}

	if item.QuantityInt() != 1 {
		t.Fatalf("expected default quantity 1, got %d", item.QuantityInt())
	}

	if item.OrderID() != "" {
		t.Fatalf("expected empty order id, got %q", item.OrderID())
	}

	if item.ProductID() != "" {
		t.Fatalf("expected empty product id, got %q", item.ProductID())
	}

	if item.ID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if item.CreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if item.UpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if item.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", sb.MAX_DATETIME, item.SoftDeletedAt())
	}

	metas, err := item.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if item.Meta("missing") != "" {
		t.Fatal("expected Meta for missing key to return empty string")
	}
}

func TestOrderLineItemDataTracking(t *testing.T) {
	item := &OrderLineItem{}

	item.SetTitle("Item Title").
		SetMemo("Memo").
		SetOrderID("ORDER1").
		SetProductID("PRODUCT1").
		SetStatus(ORDER_STATUS_COMPLETED)

	item.SetPriceFloat(15.75)
	item.SetQuantityInt(3)

	data := item.Data()
	if data == nil {
		t.Fatal("expected Data to return initialized map")
	}

	expected := map[string]string{
		COLUMN_TITLE:      "Item Title",
		COLUMN_MEMO:       "Memo",
		COLUMN_ORDER_ID:   "ORDER1",
		COLUMN_PRODUCT_ID: "PRODUCT1",
		COLUMN_STATUS:     ORDER_STATUS_COMPLETED,
		COLUMN_PRICE:      "15.75",
		COLUMN_QUANTITY:   "3",
	}

	for key, want := range expected {
		if got := data[key]; got != want {
			t.Fatalf("expected %s to be %q, got %q", key, want, got)
		}
	}

	changed := item.DataChanged()
	for key, want := range expected {
		if got := changed[key]; got != want {
			t.Fatalf("expected DataChanged to track %s as %q, got %q", key, want, got)
		}
	}

	item.MarkAsNotDirty()
	if len(item.DataChanged()) != 0 {
		t.Fatal("expected DataChanged to be empty after MarkAsNotDirty")
	}

	item.SetTitle("Updated")
	if item.DataChanged()[COLUMN_TITLE] != "Updated" {
		t.Fatalf("expected DataChanged to track updated title, got %q", item.DataChanged()[COLUMN_TITLE])
	}
}

func TestOrderLineItemCarbonHelpers(t *testing.T) {
	item := &OrderLineItem{}

	createdAt := "2025-08-09 10:11:12"
	updatedAt := "2025-09-10 11:12:13"
	softDeletedAt := "2025-10-11 12:13:14"

	if _, ok := item.SetCreatedAt(createdAt).(*OrderLineItem); !ok {
		t.Fatal("expected SetCreatedAt to return *OrderLineItem")
	}
	if item.CreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, item.CreatedAt())
	}
	if item.CreatedAtCarbon().ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", item.CreatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := item.SetUpdatedAt(updatedAt).(*OrderLineItem); !ok {
		t.Fatal("expected SetUpdatedAt to return *OrderLineItem")
	}
	if item.UpdatedAtCarbon().ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", item.UpdatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := item.SetSoftDeletedAt(softDeletedAt).(*OrderLineItem); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *OrderLineItem")
	}
	if item.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", item.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestOrderLineItemMetasRoundTrip(t *testing.T) {
	item := &OrderLineItem{}

	if err := item.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := item.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", metas["alpha"])
	}

	if item.Meta("alpha") != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", item.Meta("alpha"))
	}
}

func TestOrderLineItemMetasUpsertMergesValues(t *testing.T) {
	item := &OrderLineItem{}

	if err := item.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := item.UpsertMetas(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := item.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "updated" {
		t.Fatalf("expected alpha meta to be updated, got %q", metas["alpha"])
	}

	if metas["gamma"] != "delta" {
		t.Fatalf("expected gamma meta to be present, got %q", metas["gamma"])
	}
}

func TestOrderLineItemMetasInvalidJSON(t *testing.T) {
	item := NewOrderLineItemFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if _, err := item.Metas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if item.Meta("anything") != "" {
		t.Fatalf("expected Meta to return empty string on invalid JSON")
	}
}

func TestOrderLineItemMetasHandlesNullJSON(t *testing.T) {
	item := NewOrderLineItemFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := item.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := item.UpsertMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := item.Meta("alpha"); got != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", got)
	}
}

func TestOrderLineItemSetMetaConvenience(t *testing.T) {
	item := &OrderLineItem{}

	if err := item.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if item.Meta("key") != "value" {
		t.Fatalf("expected Meta to return %q, got %q", "value", item.Meta("key"))
	}
}

func TestOrderLineItemPriceAndQuantityHelpers(t *testing.T) {
	item := &OrderLineItem{}

	if _, ok := item.SetQuantityInt(10).(*OrderLineItem); !ok {
		t.Fatal("expected SetQuantityInt to return *OrderLineItem")
	}
	if item.QuantityInt() != 10 {
		t.Fatalf("expected QuantityInt to be 10, got %d", item.QuantityInt())
	}
	if item.Quantity() != "10" {
		t.Fatalf("expected Quantity to be \"10\", got %q", item.Quantity())
	}

	if _, ok := item.SetPriceFloat(49.95).(*OrderLineItem); !ok {
		t.Fatal("expected SetPriceFloat to return *OrderLineItem")
	}
	if item.PriceFloat() != 49.95 {
		t.Fatalf("expected PriceFloat to be 49.95, got %f", item.PriceFloat())
	}
	if item.Price() != "49.95" {
		t.Fatalf("expected Price to be \"49.95\", got %q", item.Price())
	}
}

func TestOrderLineItemSetterChainingAndGetters(t *testing.T) {
	item := &OrderLineItem{}

	if _, ok := item.SetMemo("memo").(*OrderLineItem); !ok {
		t.Fatal("expected SetMemo to return *OrderLineItem")
	}
	if item.Memo() != "memo" {
		t.Fatalf("expected Memo getter to return %q, got %q", "memo", item.Memo())
	}

	if _, ok := item.SetOrderID("ORDER1").(*OrderLineItem); !ok {
		t.Fatal("expected SetOrderID to return *OrderLineItem")
	}
	if item.OrderID() != "ORDER1" {
		t.Fatalf("expected OrderID getter to return %q, got %q", "ORDER1", item.OrderID())
	}

	if _, ok := item.SetProductID("PRODUCT1").(*OrderLineItem); !ok {
		t.Fatal("expected SetProductID to return *OrderLineItem")
	}
	if item.ProductID() != "PRODUCT1" {
		t.Fatalf("expected ProductID getter to return %q, got %q", "PRODUCT1", item.ProductID())
	}

	if _, ok := item.SetStatus(ORDER_STATUS_PARTIALLY_SHIPPED).(*OrderLineItem); !ok {
		t.Fatal("expected SetStatus to return *OrderLineItem")
	}
	if item.Status() != ORDER_STATUS_PARTIALLY_SHIPPED {
		t.Fatalf("expected Status getter to return %q, got %q", ORDER_STATUS_PARTIALLY_SHIPPED, item.Status())
	}

	if _, ok := item.SetTitle("title").(*OrderLineItem); !ok {
		t.Fatal("expected SetTitle to return *OrderLineItem")
	}
	if item.Title() != "title" {
		t.Fatalf("expected Title getter to return %q, got %q", "title", item.Title())
	}
}

func TestOrderLineItemIsSoftDeleted(t *testing.T) {
	item := &OrderLineItem{}

	if _, ok := item.SetSoftDeletedAt(sb.MAX_DATETIME).(*OrderLineItem); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *OrderLineItem")
	}
	if item.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAt to be %q, got %q", sb.MAX_DATETIME, item.SoftDeletedAt())
	}
	if item.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAtCarbon to be %q, got %q", sb.MAX_DATETIME, item.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := item.SetSoftDeletedAt("2024-01-01 00:00:00").(*OrderLineItem); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *OrderLineItem")
	}
	if item.SoftDeletedAt() != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAt to be updated, got %q", item.SoftDeletedAt())
	}
	if item.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAtCarbon to match updated value, got %q", item.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}
