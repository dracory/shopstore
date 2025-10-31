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

	if order.Status() != ORDER_STATUS_PENDING {
		t.Fatalf("expected status %q, got %q", ORDER_STATUS_PENDING, order.Status())
	}

	if order.PriceFloat() != 0 {
		t.Fatalf("expected default price 0, got %f", order.PriceFloat())
	}

	if order.QuantityInt() != 1 {
		t.Fatalf("expected default quantity 1, got %d", order.QuantityInt())
	}

	if order.Memo() != "" {
		t.Fatalf("expected empty memo, got %q", order.Memo())
	}

	if order.CustomerID() != "" {
		t.Fatalf("expected empty customer id, got %q", order.CustomerID())
	}

	if order.ID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if order.CreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if order.UpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if order.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", sb.MAX_DATETIME, order.SoftDeletedAt())
	}

	metas, err := order.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if order.Meta("missing") != "" {
		t.Fatal("expected Meta for missing key to return empty string")
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
	if order.CreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, order.CreatedAt())
	}
	if order.CreatedAtCarbon().ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", order.CreatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := order.SetUpdatedAt(updatedAt).(*Order); !ok {
		t.Fatal("expected SetUpdatedAt to return *Order")
	}
	if order.UpdatedAtCarbon().ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", order.UpdatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := order.SetSoftDeletedAt(softDeletedAt).(*Order); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Order")
	}
	if order.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", order.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestOrderMetasRoundTrip(t *testing.T) {
	order := &Order{}

	if err := order.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := order.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", metas["alpha"])
	}

	if order.Meta("alpha") != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", order.Meta("alpha"))
	}
}

func TestOrderMetasUpsertMergesValues(t *testing.T) {
	order := &Order{}

	if err := order.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := order.UpsertMetas(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := order.Metas()
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

	if _, err := order.Metas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if order.Meta("anything") != "" {
		t.Fatalf("expected Meta to return empty string on invalid JSON")
	}
}

func TestOrderSetMetaConvenience(t *testing.T) {
	order := &Order{}

	if err := order.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if order.Meta("key") != "value" {
		t.Fatalf("expected Meta to return %q, got %q", "value", order.Meta("key"))
	}
}

func TestOrderPriceAndQuantityHelpers(t *testing.T) {
	order := &Order{}

	if _, ok := order.SetQuantityInt(10).(*Order); !ok {
		t.Fatal("expected SetQuantityInt to return *Order")
	}
	if order.QuantityInt() != 10 {
		t.Fatalf("expected QuantityInt to be 10, got %d", order.QuantityInt())
	}
	if order.Quantity() != "10" {
		t.Fatalf("expected Quantity to be \"10\", got %q", order.Quantity())
	}

	if _, ok := order.SetPriceFloat(49.95).(*Order); !ok {
		t.Fatal("expected SetPriceFloat to return *Order")
	}
	if order.PriceFloat() != 49.95 {
		t.Fatalf("expected PriceFloat to be 49.95, got %f", order.PriceFloat())
	}
	if order.Price() != "49.95" {
		t.Fatalf("expected Price to be \"49.95\", got %q", order.Price())
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
	if order.CustomerID() != "cust" {
		t.Fatalf("expected CustomerID getter to return %q, got %q", "cust", order.CustomerID())
	}

	if _, ok := order.SetMemo("memo").(*Order); !ok {
		t.Fatal("expected SetMemo to return *Order")
	}
	if order.Memo() != "memo" {
		t.Fatalf("expected Memo getter to return %q, got %q", "memo", order.Memo())
	}

	if _, ok := order.SetStatus(ORDER_STATUS_REFUNDED).(*Order); !ok {
		t.Fatal("expected SetStatus to return *Order")
	}
	if order.Status() != ORDER_STATUS_REFUNDED {
		t.Fatalf("expected Status getter to return %q, got %q", ORDER_STATUS_REFUNDED, order.Status())
	}
}

func TestOrderIsSoftDeleted(t *testing.T) {
	order := &Order{}

	if _, ok := order.SetSoftDeletedAt(sb.MAX_DATETIME).(*Order); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Order")
	}
	if order.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAt to be %q, got %q", sb.MAX_DATETIME, order.SoftDeletedAt())
	}
	if order.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAtCarbon to be %q, got %q", sb.MAX_DATETIME, order.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := order.SetSoftDeletedAt("2024-01-01 00:00:00").(*Order); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Order")
	}
	if order.SoftDeletedAt() != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAt to be updated, got %q", order.SoftDeletedAt())
	}
	if order.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAtCarbon to match updated value, got %q", order.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}
