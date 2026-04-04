package shopstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewProductDefaults(t *testing.T) {
	product := NewProduct()
	if product == nil {
		t.Fatal("NewProduct returned nil")
	}

	if product.GetStatus() != PRODUCT_STATUS_DRAFT {
		t.Fatalf("expected status %q, got %q", PRODUCT_STATUS_DRAFT, product.GetStatus())
	}

	if product.GetTitle() != "" {
		t.Fatalf("expected empty title, got %q", product.GetTitle())
	}

	if product.GetDescription() != "" {
		t.Fatalf("expected empty description, got %q", product.GetDescription())
	}

	if product.GetShortDescription() != "" {
		t.Fatalf("expected empty short description, got %q", product.GetShortDescription())
	}

	if product.GetQuantityInt() != 0 {
		t.Fatalf("expected default quantity 0, got %d", product.GetQuantityInt())
	}

	if product.GetPriceFloat() != 0 {
		t.Fatalf("expected default price 0, got %f", product.GetPriceFloat())
	}

	if product.GetMemo() != "" {
		t.Fatalf("expected empty memo, got %q", product.GetMemo())
	}

	if product.GetID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if product.GetCreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if product.GetUpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if product.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", sb.MAX_DATETIME, product.GetSoftDeletedAt())
	}

	metas, err := product.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if product.GetMeta("missing") != "" {
		t.Fatal("expected GetMeta for missing key to return empty string")
	}
}

func TestProductDataTracking(t *testing.T) {
	product := &Product{}

	product.SetTitle("Title").
		SetDescription("Desc").
		SetShortDescription("Short").
		SetMemo("Memo").
		SetStatus(PRODUCT_STATUS_ACTIVE)

	product.SetPriceFloat(19.99)
	product.SetQuantityInt(5)

	data := product.Data()
	if data == nil {
		t.Fatal("expected Data to return initialized map")
	}

	expected := map[string]string{
		COLUMN_TITLE:             "Title",
		COLUMN_DESCRIPTION:       "Desc",
		COLUMN_SHORT_DESCRIPTION: "Short",
		COLUMN_MEMO:              "Memo",
		COLUMN_STATUS:            PRODUCT_STATUS_ACTIVE,
		COLUMN_PRICE:             "19.99",
		COLUMN_QUANTITY:          "5",
	}

	for key, want := range expected {
		if got := data[key]; got != want {
			t.Fatalf("expected %s to be %q, got %q", key, want, got)
		}
	}

	changed := product.DataChanged()
	for key, want := range expected {
		if got := changed[key]; got != want {
			t.Fatalf("expected DataChanged to track %s as %q, got %q", key, want, got)
		}
	}

	product.MarkAsNotDirty()
	if len(product.DataChanged()) != 0 {
		t.Fatal("expected DataChanged to be empty after MarkAsNotDirty")
	}

	product.SetTitle("Updated")
	if product.DataChanged()[COLUMN_TITLE] != "Updated" {
		t.Fatalf("expected DataChanged to track updated title, got %q", product.DataChanged()[COLUMN_TITLE])
	}
}

func TestProductCarbonHelpers(t *testing.T) {
	product := &Product{}

	createdAt := "2025-01-02 03:04:05"
	updatedAt := "2025-02-03 04:05:06"
	softDeletedAt := "2025-03-04 05:06:07"

	if _, ok := product.SetCreatedAt(createdAt).(*Product); !ok {
		t.Fatal("expected SetCreatedAt to return *Product")
	}
	if product.GetCreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, product.GetCreatedAt())
	}
	if product.GetCreatedAtCarbon().ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", product.GetCreatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := product.SetUpdatedAt(updatedAt).(*Product); !ok {
		t.Fatal("expected SetUpdatedAt to return *Product")
	}
	if product.GetUpdatedAtCarbon().ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", product.GetUpdatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := product.SetSoftDeletedAt(softDeletedAt).(*Product); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Product")
	}
	if product.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", product.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestProductMetasRoundTrip(t *testing.T) {
	product := &Product{}

	if err := product.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := product.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", metas["alpha"])
	}

	if product.GetMeta("alpha") != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", product.GetMeta("alpha"))
	}
}

func TestProductMetasUpsertMergesValues(t *testing.T) {
	product := &Product{}

	if err := product.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := product.MetasUpsert(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := product.GetMetas()
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

func TestProductMetasHandlesNullJSON(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := product.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := product.MetasUpsert(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := product.GetMeta("alpha"); got != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", got)
	}
}

func TestProductMetasInvalidJSON(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if _, err := product.GetMetas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if product.GetMeta("anything") != "" {
		t.Fatalf("expected GetMeta to return empty string on invalid JSON")
	}
}

func TestProductSetMetaConvenience(t *testing.T) {
	product := &Product{}

	if err := product.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if product.GetMeta("key") != "value" {
		t.Fatalf("expected GetMeta to return %q, got %q", "value", product.GetMeta("key"))
	}
}

func TestProductQuantityAndPriceHelpers(t *testing.T) {
	product := &Product{}

	if _, ok := product.SetQuantityInt(10).(*Product); !ok {
		t.Fatal("expected SetQuantityInt to return *Product")
	}
	if product.GetQuantityInt() != 10 {
		t.Fatalf("expected QuantityInt to be 10, got %d", product.GetQuantityInt())
	}
	if product.GetQuantity() != "10" {
		t.Fatalf("expected Quantity to be \"10\", got %q", product.GetQuantity())
	}

	if _, ok := product.SetPriceFloat(49.95).(*Product); !ok {
		t.Fatal("expected SetPriceFloat to return *Product")
	}
	if product.GetPriceFloat() != 49.95 {
		t.Fatalf("expected PriceFloat to be 49.95, got %f", product.GetPriceFloat())
	}
	if product.GetPrice() != "49.95" {
		t.Fatalf("expected Price to be \"49.95\", got %q", product.GetPrice())
	}

	if product.IsFree() {
		t.Fatal("expected product with positive price not to be free")
	}

	product.SetPriceFloat(0)
	if !product.IsFree() {
		t.Fatal("expected product with zero price to be free")
	}
}

func TestProductStatusPredicates(t *testing.T) {
	product := &Product{}

	product.SetStatus(PRODUCT_STATUS_ACTIVE)
	if !product.IsActive() {
		t.Fatal("expected product to be active")
	}
	if product.IsDraft() {
		t.Fatal("expected product not to be draft when active")
	}

	product.SetStatus(PRODUCT_STATUS_DISABLED)
	if !product.IsDisabled() {
		t.Fatal("expected product to be disabled")
	}

	product.SetStatus(PRODUCT_STATUS_DRAFT)
	if !product.IsDraft() {
		t.Fatal("expected product to be draft")
	}
}

func TestProductIsSoftDeleted(t *testing.T) {
	product := &Product{}

	product.SetSoftDeletedAt(sb.MAX_DATETIME)
	if product.IsSoftDeleted() {
		t.Fatal("expected product not to be soft deleted with MAX_DATETIME")
	}

	product.SetSoftDeletedAt("2000-01-01 00:00:00")
	if !product.IsSoftDeleted() {
		t.Fatal("expected product to be soft deleted for past timestamp")
	}
}

func TestProductSlug(t *testing.T) {
	product := &Product{}
	product.SetTitle("  Hello World!  ")

	if slug := product.Slug(); slug != "hello-world" {
		t.Fatalf("expected slug to be %q, got %q", "hello-world", slug)
	}
}

func TestProductMetaRemove(t *testing.T) {
	product := &Product{}

	if err := product.SetMetas(map[string]string{"key1": "value1", "key2": "value2"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := product.MetaRemove("key1"); err != nil {
		t.Fatalf("unexpected error removing meta: %v", err)
	}

	if product.GetMeta("key1") != "" {
		t.Fatalf("expected key1 to be removed, got %q", product.GetMeta("key1"))
	}

	if product.GetMeta("key2") != "value2" {
		t.Fatalf("expected key2 to still exist, got %q", product.GetMeta("key2"))
	}
}

func TestProductMetaRemoveNonExistent(t *testing.T) {
	product := &Product{}

	if err := product.SetMetas(map[string]string{"key1": "value1"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := product.MetaRemove("nonexistent"); err != nil {
		t.Fatalf("unexpected error removing non-existent meta: %v", err)
	}

	if product.GetMeta("key1") != "value1" {
		t.Fatalf("expected key1 to still exist, got %q", product.GetMeta("key1"))
	}
}

func TestProductMetasRemove(t *testing.T) {
	product := &Product{}

	if err := product.SetMetas(map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := product.MetasRemove([]string{"key1", "key2"}); err != nil {
		t.Fatalf("unexpected error removing metas: %v", err)
	}

	if product.GetMeta("key1") != "" {
		t.Fatalf("expected key1 to be removed, got %q", product.GetMeta("key1"))
	}

	if product.GetMeta("key2") != "" {
		t.Fatalf("expected key2 to be removed, got %q", product.GetMeta("key2"))
	}

	if product.GetMeta("key3") != "value3" {
		t.Fatalf("expected key3 to still exist, got %q", product.GetMeta("key3"))
	}
}

func TestProductMetasRemoveEmptySlice(t *testing.T) {
	product := &Product{}

	if err := product.SetMetas(map[string]string{"key1": "value1"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := product.MetasRemove([]string{}); err != nil {
		t.Fatalf("unexpected error removing empty slice: %v", err)
	}

	if product.GetMeta("key1") != "value1" {
		t.Fatalf("expected key1 to still exist, got %q", product.GetMeta("key1"))
	}
}

func TestProductMetaRemoveErrorPropagation(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if err := product.MetaRemove("key"); err == nil {
		t.Fatal("expected error when removing meta with invalid JSON")
	}
}

func TestProductMetasRemoveErrorPropagation(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if err := product.MetasRemove([]string{"key"}); err == nil {
		t.Fatal("expected error when removing metas with invalid JSON")
	}
}

func TestProductHasStock(t *testing.T) {
	product := &Product{}

	// Default quantity is 0
	if product.HasStock() {
		t.Fatal("expected product with quantity 0 not to have stock")
	}

	if !product.IsOutOfStock() {
		t.Fatal("expected product with quantity 0 to be out of stock")
	}

	product.SetQuantityInt(1)
	if !product.HasStock() {
		t.Fatal("expected product with quantity 1 to have stock")
	}

	if product.IsOutOfStock() {
		t.Fatal("expected product with quantity 1 not to be out of stock")
	}

	product.SetQuantityInt(100)
	if !product.HasStock() {
		t.Fatal("expected product with quantity 100 to have stock")
	}
}

func TestProductIsOutOfStock(t *testing.T) {
	product := &Product{}

	// Default quantity is 0
	if !product.IsOutOfStock() {
		t.Fatal("expected product with quantity 0 to be out of stock")
	}

	product.SetQuantityInt(-5)
	if !product.IsOutOfStock() {
		t.Fatal("expected product with negative quantity to be out of stock")
	}

	product.SetQuantityInt(10)
	if product.IsOutOfStock() {
		t.Fatal("expected product with positive quantity not to be out of stock")
	}
}

func TestProductIsPaid(t *testing.T) {
	product := &Product{}

	// Default price is 0 (free)
	if product.IsPaid() {
		t.Fatal("expected product with price 0 not to be paid")
	}

	product.SetPriceFloat(0.01)
	if !product.IsPaid() {
		t.Fatal("expected product with price 0.01 to be paid")
	}

	product.SetPriceFloat(99.99)
	if !product.IsPaid() {
		t.Fatal("expected product with price 99.99 to be paid")
	}
}
