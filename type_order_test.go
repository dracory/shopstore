package shopstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewOrderDefaults(t *testing.T) {
	order := NewOrder()
	if order == nil {
		t.Fatal("NewOrder returned nil")
	}

	if order.GetStatus() != ORDER_STATUS_PENDING {
		t.Fatalf("expected status %q, got %q", ORDER_STATUS_PENDING, order.GetStatus())
	}

	if order.GetPriceFloat() != 0 {
		t.Fatalf("expected default price 0, got %f", order.GetPriceFloat())
	}

	if order.GetQuantityInt() != 1 {
		t.Fatalf("expected default quantity 1, got %d", order.GetQuantityInt())
	}

	if order.GetMemo() != "" {
		t.Fatalf("expected empty memo, got %q", order.GetMemo())
	}

	if order.GetCustomerID() != "" {
		t.Fatalf("expected empty customer id, got %q", order.GetCustomerID())
	}

	if order.GetID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if order.GetCreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if order.GetUpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if order.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", sb.MAX_DATETIME, order.GetSoftDeletedAt())
	}

	metas, err := order.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if order.GetMeta("missing") != "" {
		t.Fatal("expected GetMeta for missing key to return empty string")
	}
}

func TestOrderDataTracking(t *testing.T) {
	order := &Order{}

	order.SetCustomerID("CUST1").
		SetMemo("Memo").
		SetStatus(ORDER_STATUS_COMPLETED)

	order.SetPriceFloat(19.99)
	order.SetQuantityInt(4)

	data := order.Data()
	if data == nil {
		t.Fatal("expected Data to return initialized map")
	}

	expected := map[string]string{
		COLUMN_CUSTOMER_ID: "CUST1",
		COLUMN_MEMO:        "Memo",
		COLUMN_STATUS:      ORDER_STATUS_COMPLETED,
		COLUMN_PRICE:       "19.99",
		COLUMN_QUANTITY:    "4",
	}

	for key, want := range expected {
		if got := data[key]; got != want {
			t.Fatalf("expected %s to be %q, got %q", key, want, got)
		}
	}

	changed := order.DataChanged()
	for key, want := range expected {
		if got := changed[key]; got != want {
			t.Fatalf("expected DataChanged to track %s as %q, got %q", key, want, got)
		}
	}

	order.MarkAsNotDirty()
	if len(order.DataChanged()) != 0 {
		t.Fatal("expected DataChanged to be empty after MarkAsNotDirty")
	}

	order.SetMemo("Updated")
	if order.DataChanged()[COLUMN_MEMO] != "Updated" {
		t.Fatalf("expected DataChanged to track updated memo, got %q", order.DataChanged()[COLUMN_MEMO])
	}
}

func TestOrderCarbonHelpers(t *testing.T) {
	order := &Order{}

	createdAt := "2025-04-05 06:07:08"
	updatedAt := "2025-05-06 07:08:09"
	softDeletedAt := "2026-01-01 00:00:00"

	if _, ok := order.SetCreatedAt(createdAt).(*Order); !ok {
		t.Fatal("expected SetCreatedAt to return *Order")
	}
	if order.GetCreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, order.GetCreatedAt())
	}
	if order.GetCreatedAtCarbon().ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", order.GetCreatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := order.SetUpdatedAt(updatedAt).(*Order); !ok {
		t.Fatal("expected SetUpdatedAt to return *Order")
	}
	if order.GetUpdatedAtCarbon().ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", order.GetUpdatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := order.SetSoftDeletedAt(softDeletedAt).(*Order); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Order")
	}
	if order.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", order.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestOrderMetasRoundTrip(t *testing.T) {
	order := &Order{}

	if err := order.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := order.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", metas["alpha"])
	}

	if order.GetMeta("alpha") != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", order.GetMeta("alpha"))
	}
}

func TestOrderMetasUpsertMergesValues(t *testing.T) {
	order := &Order{}

	if err := order.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := order.MetasUpsert(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := order.GetMetas()
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

func TestOrderMetasInvalidJSON(t *testing.T) {
	order := NewOrderFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if _, err := order.GetMetas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if order.GetMeta("anything") != "" {
		t.Fatalf("expected GetMeta to return empty string on invalid JSON")
	}
}

func TestOrderMetasHandlesNullJSON(t *testing.T) {
	order := NewOrderFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := order.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := order.MetasUpsert(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := order.GetMeta("alpha"); got != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", got)
	}
}

func TestOrderSetMetaConvenience(t *testing.T) {
	order := &Order{}

	if err := order.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if order.GetMeta("key") != "value" {
		t.Fatalf("expected GetMeta to return %q, got %q", "value", order.GetMeta("key"))
	}
}

func TestOrderPriceAndQuantityHelpers(t *testing.T) {
	order := &Order{}

	if _, ok := order.SetQuantityInt(10).(*Order); !ok {
		t.Fatal("expected SetQuantityInt to return *Order")
	}
	if order.GetQuantityInt() != 10 {
		t.Fatalf("expected QuantityInt to be 10, got %d", order.GetQuantityInt())
	}
	if order.GetQuantity() != "10" {
		t.Fatalf("expected Quantity to be \"10\", got %q", order.GetQuantity())
	}

	if _, ok := order.SetPriceFloat(49.95).(*Order); !ok {
		t.Fatal("expected SetPriceFloat to return *Order")
	}
	if order.GetPriceFloat() != 49.95 {
		t.Fatalf("expected PriceFloat to be 49.95, got %f", order.GetPriceFloat())
	}
	if order.GetPrice() != "49.95" {
		t.Fatalf("expected Price to be \"49.95\", got %q", order.GetPrice())
	}
}

func TestOrderStatusPredicates(t *testing.T) {
	order := &Order{}

	cases := []struct {
		status  string
		checker func(OrderInterface) bool
	}{
		{ORDER_STATUS_AWAITING_FULFILLMENT, OrderInterface.IsAwaitingFulfillment},
		{ORDER_STATUS_AWAITING_PAYMENT, OrderInterface.IsAwaitingPayment},
		{ORDER_STATUS_AWAITING_PICKUP, OrderInterface.IsAwaitingPickup},
		{ORDER_STATUS_AWAITING_SHIPMENT, OrderInterface.IsAwaitingShipment},
		{ORDER_STATUS_CANCELLED, OrderInterface.IsCancelled},
		{ORDER_STATUS_COMPLETED, OrderInterface.IsCompleted},
		{ORDER_STATUS_DECLINED, OrderInterface.IsDeclined},
		{ORDER_STATUS_DISPUTED, OrderInterface.IsDisputed},
		{ORDER_STATUS_MANUAL_VERIFICATION_REQUIRED, OrderInterface.IsManualVerificationRequired},
		{ORDER_STATUS_PENDING, OrderInterface.IsPending},
		{ORDER_STATUS_REFUNDED, OrderInterface.IsRefunded},
		{ORDER_STATUS_SHIPPED, OrderInterface.IsShipped},
	}

	for _, tc := range cases {
		order.SetStatus(tc.status)
		if !tc.checker(order) {
			t.Fatalf("expected predicate for %q to return true", tc.status)
		}
	}
}

func TestOrderSetterChainingAndGetters(t *testing.T) {
	order := &Order{}

	if _, ok := order.SetCustomerID("cust").(*Order); !ok {
		t.Fatal("expected SetCustomerID to return *Order")
	}
	if order.GetCustomerID() != "cust" {
		t.Fatalf("expected GetCustomerID getter to return %q, got %q", "cust", order.GetCustomerID())
	}

	if _, ok := order.SetMemo("memo").(*Order); !ok {
		t.Fatal("expected SetMemo to return *Order")
	}
	if order.GetMemo() != "memo" {
		t.Fatalf("expected GetMemo getter to return %q, got %q", "memo", order.GetMemo())
	}

	if _, ok := order.SetStatus(ORDER_STATUS_REFUNDED).(*Order); !ok {
		t.Fatal("expected SetStatus to return *Order")
	}
	if order.GetStatus() != ORDER_STATUS_REFUNDED {
		t.Fatalf("expected GetStatus getter to return %q, got %q", ORDER_STATUS_REFUNDED, order.GetStatus())
	}
}

func TestOrderIsSoftDeleted(t *testing.T) {
	order := &Order{}

	if _, ok := order.SetSoftDeletedAt(sb.MAX_DATETIME).(*Order); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Order")
	}
	if order.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAt to be %q, got %q", sb.MAX_DATETIME, order.GetSoftDeletedAt())
	}
	if order.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAtCarbon to be %q, got %q", sb.MAX_DATETIME, order.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := order.SetSoftDeletedAt("2024-01-01 00:00:00").(*Order); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Order")
	}
	if order.GetSoftDeletedAt() != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAt to be updated, got %q", order.GetSoftDeletedAt())
	}
	if order.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAtCarbon to match updated value, got %q", order.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestOrderMetaRemove(t *testing.T) {
	order := &Order{}

	if err := order.SetMetas(map[string]string{"key1": "value1", "key2": "value2"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := order.MetaRemove("key1"); err != nil {
		t.Fatalf("unexpected error removing meta: %v", err)
	}

	if order.GetMeta("key1") != "" {
		t.Fatalf("expected key1 to be removed, got %q", order.GetMeta("key1"))
	}

	if order.GetMeta("key2") != "value2" {
		t.Fatalf("expected key2 to still exist, got %q", order.GetMeta("key2"))
	}
}

func TestOrderMetaRemoveNonExistent(t *testing.T) {
	order := &Order{}

	if err := order.SetMetas(map[string]string{"key1": "value1"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := order.MetaRemove("nonexistent"); err != nil {
		t.Fatalf("unexpected error removing non-existent meta: %v", err)
	}

	if order.GetMeta("key1") != "value1" {
		t.Fatalf("expected key1 to still exist, got %q", order.GetMeta("key1"))
	}
}

func TestOrderMetasRemove(t *testing.T) {
	order := &Order{}

	if err := order.SetMetas(map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := order.MetasRemove([]string{"key1", "key2"}); err != nil {
		t.Fatalf("unexpected error removing metas: %v", err)
	}

	if order.GetMeta("key1") != "" {
		t.Fatalf("expected key1 to be removed, got %q", order.GetMeta("key1"))
	}

	if order.GetMeta("key2") != "" {
		t.Fatalf("expected key2 to be removed, got %q", order.GetMeta("key2"))
	}

	if order.GetMeta("key3") != "value3" {
		t.Fatalf("expected key3 to still exist, got %q", order.GetMeta("key3"))
	}
}

func TestOrderMetasRemoveEmptySlice(t *testing.T) {
	order := &Order{}

	if err := order.SetMetas(map[string]string{"key1": "value1"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := order.MetasRemove([]string{}); err != nil {
		t.Fatalf("unexpected error removing empty slice: %v", err)
	}

	if order.GetMeta("key1") != "value1" {
		t.Fatalf("expected key1 to still exist, got %q", order.GetMeta("key1"))
	}
}

func TestOrderMetaRemoveErrorPropagation(t *testing.T) {
	order := NewOrderFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if err := order.MetaRemove("key"); err == nil {
		t.Fatal("expected error when removing meta with invalid JSON")
	}
}

func TestOrderMetasRemoveErrorPropagation(t *testing.T) {
	order := NewOrderFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if err := order.MetasRemove([]string{"key"}); err == nil {
		t.Fatal("expected error when removing metas with invalid JSON")
	}
}
