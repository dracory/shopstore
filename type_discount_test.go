package shopstore

import (
	"strings"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewDiscountDefaults(t *testing.T) {
	discount := NewDiscount()
	if discount == nil {
		t.Fatal("NewDiscount returned nil")
	}

	if discount.Status() != DISCOUNT_STATUS_DRAFT {
		t.Fatalf("expected status %q, got %q", DISCOUNT_STATUS_DRAFT, discount.Status())
	}

	if discount.Type() != DISCOUNT_TYPE_PERCENT {
		t.Fatalf("expected type %q, got %q", DISCOUNT_TYPE_PERCENT, discount.Type())
	}

	if discount.Title() != "" {
		t.Fatalf("expected empty title, got %q", discount.Title())
	}

	if discount.Description() != "" {
		t.Fatalf("expected empty description, got %q", discount.Description())
	}

	if discount.Memo() != "" {
		t.Fatalf("expected empty memo, got %q", discount.Memo())
	}

	if discount.Amount() != 0 {
		t.Fatalf("expected amount 0, got %f", discount.Amount())
	}

	if discount.Code() == "" {
		t.Fatal("expected generated code to be non-empty")
	}

	if discount.StartsAt() != sb.NULL_DATETIME {
		t.Fatalf("expected starts at %q, got %q", sb.NULL_DATETIME, discount.StartsAt())
	}

	if discount.EndsAt() != sb.NULL_DATETIME {
		t.Fatalf("expected ends at %q, got %q", sb.NULL_DATETIME, discount.EndsAt())
	}

	if discount.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", sb.MAX_DATETIME, discount.SoftDeletedAt())
	}

	if discount.ID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if discount.CreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if discount.UpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	metas, err := discount.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if discount.Meta("missing") != "" {
		t.Fatal("expected Meta for missing key to return empty string")
	}
}

func TestNewDiscountCodeUsesCrockfordAlphabet(t *testing.T) {
	discount := NewDiscount()
	code := strings.ToUpper(discount.Code())

	const allowed = "BCDFGHJKLMNPQRSTVWXYZ23456789"
	for _, r := range code {
		if !strings.ContainsRune(allowed, r) {
			t.Fatalf("generated code %q contains disallowed character %q", code, r)
		}
	}
}

func TestDiscountDataTracking(t *testing.T) {
	discount := &Discount{}

	discount.SetTitle("Title").
		SetDescription("Desc").
		SetMemo("Memo").
		SetCode("CODE").
		SetStatus(DISCOUNT_STATUS_ACTIVE).
		SetType(DISCOUNT_TYPE_AMOUNT)

	if _, ok := discount.SetAmount(12.34).(*Discount); !ok {
		t.Fatal("expected SetAmount to return *Discount")
	}

	data := discount.Data()
	if data == nil {
		t.Fatal("expected Data to return initialized map")
	}

	expected := map[string]string{
		COLUMN_TITLE:       "Title",
		COLUMN_DESCRIPTION: "Desc",
		COLUMN_MEMO:        "Memo",
		COLUMN_CODE:        "CODE",
		COLUMN_STATUS:      DISCOUNT_STATUS_ACTIVE,
		COLUMN_TYPE:        DISCOUNT_TYPE_AMOUNT,
		COLUMN_AMOUNT:      "12.34",
	}

	for key, want := range expected {
		if got := data[key]; got != want {
			t.Fatalf("expected %s to be %q, got %q", key, want, got)
		}
	}

	changed := discount.DataChanged()
	for key, want := range expected {
		if got := changed[key]; got != want {
			t.Fatalf("expected DataChanged to track %s as %q, got %q", key, want, got)
		}
	}

	discount.MarkAsNotDirty()
	if len(discount.DataChanged()) != 0 {
		t.Fatal("expected DataChanged to be empty after MarkAsNotDirty")
	}

	discount.SetTitle("Updated")
	if discount.DataChanged()[COLUMN_TITLE] != "Updated" {
		t.Fatalf("expected DataChanged to track updated title, got %q", discount.DataChanged()[COLUMN_TITLE])
	}
}

func TestDiscountCarbonHelpers(t *testing.T) {
	discount := &Discount{}

	createdAt := "2025-01-01 00:00:00"
	updatedAt := "2025-02-02 12:34:56"
	startsAt := "2025-03-03 09:00:00"
	endsAt := "2025-04-04 18:30:00"
	softDeletedAt := "2026-05-05 00:00:00"

	if _, ok := discount.SetCreatedAt(createdAt).(*Discount); !ok {
		t.Fatal("expected SetCreatedAt to return *Discount")
	}
	if discount.CreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, discount.CreatedAt())
	}
	if discount.CreatedAtCarbon().ToDateTimeString() != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", discount.CreatedAtCarbon().ToDateTimeString())
	}

	if _, ok := discount.SetUpdatedAt(updatedAt).(*Discount); !ok {
		t.Fatal("expected SetUpdatedAt to return *Discount")
	}
	if discount.UpdatedAtCarbon().ToDateTimeString() != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", discount.UpdatedAtCarbon().ToDateTimeString())
	}

	if _, ok := discount.SetStartsAt(startsAt).(*Discount); !ok {
		t.Fatal("expected SetStartsAt to return *Discount")
	}
	if discount.StartsAtCarbon().ToDateTimeString(carbon.UTC) != startsAt {
		t.Fatalf("expected StartsAtCarbon to match input, got %q", discount.StartsAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := discount.SetEndsAt(endsAt).(*Discount); !ok {
		t.Fatal("expected SetEndsAt to return *Discount")
	}
	if discount.EndsAtCarbon().ToDateTimeString(carbon.UTC) != endsAt {
		t.Fatalf("expected EndsAtCarbon to match input, got %q", discount.EndsAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := discount.SetSoftDeletedAt(softDeletedAt).(*Discount); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Discount")
	}
	if discount.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", discount.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestDiscountMetasRoundTrip(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := discount.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if got := metas["alpha"]; got != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", got)
	}

	if got := discount.Meta("alpha"); got != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", got)
	}
}

func TestDiscountMetasUpsertMergesValues(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := discount.MetasUpsert(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := discount.Metas()
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

func TestDiscountMetaRemove(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMetas(map[string]string{"alpha": "beta", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := discount.MetaRemove("alpha"); err != nil {
		t.Fatalf("unexpected error removing meta: %v", err)
	}

	if discount.Meta("alpha") != "" {
		t.Fatal("expected removed meta to return empty string")
	}

	metas, err := discount.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if _, exists := metas["alpha"]; exists {
		t.Fatal("expected alpha meta to be removed from stored metas")
	}

	if metas["gamma"] != "delta" {
		t.Fatalf("expected gamma meta to remain, got %q", metas["gamma"])
	}
}

func TestDiscountMetasRemoveList(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMetas(map[string]string{"alpha": "beta", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := discount.MetasRemove([]string{"alpha", "gamma"}); err != nil {
		t.Fatalf("unexpected error removing metas: %v", err)
	}

	metas, err := discount.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected all metas to be removed, got %v", metas)
	}
}

func TestDiscountMetasInvalidJSON(t *testing.T) {
	discount := NewDiscountFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if _, err := discount.Metas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if got := discount.Meta("anything"); got != "" {
		t.Fatalf("expected Meta to return empty string on invalid JSON, got %q", got)
	}
}

func TestDiscountMetasHandlesNullJSON(t *testing.T) {
	discount := NewDiscountFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := discount.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := discount.MetasUpsert(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := discount.Meta("alpha"); got != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", got)
	}
}

func TestDiscountSetMetaConvenience(t *testing.T) {
	discount := &Discount{}

	if err := discount.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if got := discount.Meta("key"); got != "value" {
		t.Fatalf("expected Meta to return %q, got %q", "value", got)
	}
}

func TestDiscountSetterChainingAndGetters(t *testing.T) {
	discount := &Discount{}

	if _, ok := discount.SetDescription("desc").(*Discount); !ok {
		t.Fatal("expected SetDescription to return *Discount")
	}
	if discount.Description() != "desc" {
		t.Fatalf("expected Description getter to return %q, got %q", "desc", discount.Description())
	}

	if _, ok := discount.SetMemo("memo").(*Discount); !ok {
		t.Fatal("expected SetMemo to return *Discount")
	}
	if discount.Memo() != "memo" {
		t.Fatalf("expected Memo getter to return %q, got %q", "memo", discount.Memo())
	}

	if _, ok := discount.SetTitle("title").(*Discount); !ok {
		t.Fatal("expected SetTitle to return *Discount")
	}
	if discount.Title() != "title" {
		t.Fatalf("expected Title getter to return %q, got %q", "title", discount.Title())
	}

	if _, ok := discount.SetStatus(DISCOUNT_STATUS_INACTIVE).(*Discount); !ok {
		t.Fatal("expected SetStatus to return *Discount")
	}
	if discount.Status() != DISCOUNT_STATUS_INACTIVE {
		t.Fatalf("expected Status getter to return %q, got %q", DISCOUNT_STATUS_INACTIVE, discount.Status())
	}

	if _, ok := discount.SetType(DISCOUNT_TYPE_AMOUNT).(*Discount); !ok {
		t.Fatal("expected SetType to return *Discount")
	}
	if discount.Type() != DISCOUNT_TYPE_AMOUNT {
		t.Fatalf("expected Type getter to return %q, got %q", DISCOUNT_TYPE_AMOUNT, discount.Type())
	}

	start := "2024-01-01 00:00:00"
	end := "2024-12-31 23:59:59"

	if _, ok := discount.SetStartsAt(start).(*Discount); !ok {
		t.Fatal("expected SetStartsAt to return *Discount")
	}
	if discount.StartsAt() != start {
		t.Fatalf("expected StartsAt getter to return %q, got %q", start, discount.StartsAt())
	}

	if _, ok := discount.SetEndsAt(end).(*Discount); !ok {
		t.Fatal("expected SetEndsAt to return *Discount")
	}
	if discount.EndsAt() != end {
		t.Fatalf("expected EndsAt getter to return %q, got %q", end, discount.EndsAt())
	}
}
