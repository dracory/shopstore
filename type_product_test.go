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

	if product.Status() != PRODUCT_STATUS_DRAFT {
		t.Fatalf("expected status %q, got %q", PRODUCT_STATUS_DRAFT, product.Status())
	}

	if product.Title() != "" {
		t.Fatalf("expected empty title, got %q", product.Title())
	}

	if product.Description() != "" {
		t.Fatalf("expected empty description, got %q", product.Description())
	}

	if product.ShortDescription() != "" {
		t.Fatalf("expected empty short description, got %q", product.ShortDescription())
	}

	if product.QuantityInt() != 0 {
		t.Fatalf("expected default quantity 0, got %d", product.QuantityInt())
	}

	if product.PriceFloat() != 0 {
		t.Fatalf("expected default price 0, got %f", product.PriceFloat())
	}

	if product.Memo() != "" {
		t.Fatalf("expected empty memo, got %q", product.Memo())
	}

	if product.ID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if product.CreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if product.UpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if product.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", sb.MAX_DATETIME, product.SoftDeletedAt())
	}

	metas, err := product.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if product.Meta("missing") != "" {
		t.Fatal("expected Meta for missing key to return empty string")
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
	if product.CreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, product.CreatedAt())
	}
	if product.CreatedAtCarbon().ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", product.CreatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := product.SetUpdatedAt(updatedAt).(*Product); !ok {
		t.Fatal("expected SetUpdatedAt to return *Product")
	}
	if product.UpdatedAtCarbon().ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", product.UpdatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := product.SetSoftDeletedAt(softDeletedAt).(*Product); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Product")
	}
	if product.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", product.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestProductMetasRoundTrip(t *testing.T) {
	product := &Product{}

	if err := product.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := product.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", metas["alpha"])
	}

	if product.Meta("alpha") != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", product.Meta("alpha"))
	}
}

func TestProductMetasUpsertMergesValues(t *testing.T) {
	product := &Product{}

	if err := product.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := product.UpsertMetas(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := product.Metas()
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

func TestProductMetasInvalidJSON(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if _, err := product.Metas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if product.Meta("anything") != "" {
		t.Fatalf("expected Meta to return empty string on invalid JSON")
	}
}

func TestProductSetMetaConvenience(t *testing.T) {
	product := &Product{}

	if err := product.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if product.Meta("key") != "value" {
		t.Fatalf("expected Meta to return %q, got %q", "value", product.Meta("key"))
	}
}

func TestProductQuantityAndPriceHelpers(t *testing.T) {
	product := &Product{}

	if _, ok := product.SetQuantityInt(10).(*Product); !ok {
		t.Fatal("expected SetQuantityInt to return *Product")
	}
	if product.QuantityInt() != 10 {
		t.Fatalf("expected QuantityInt to be 10, got %d", product.QuantityInt())
	}
	if product.Quantity() != "10" {
		t.Fatalf("expected Quantity to be \"10\", got %q", product.Quantity())
	}

	if _, ok := product.SetPriceFloat(49.95).(*Product); !ok {
		t.Fatal("expected SetPriceFloat to return *Product")
	}
	if product.PriceFloat() != 49.95 {
		t.Fatalf("expected PriceFloat to be 49.95, got %f", product.PriceFloat())
	}
	if product.Price() != "49.95" {
		t.Fatalf("expected Price to be \"49.95\", got %q", product.Price())
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
