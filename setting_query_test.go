package shopstore

import "testing"

func TestSettingQueryValidate(t *testing.T) {
	query := NewSettingQuery()
	if err := query.Validate(); err != nil {
		t.Fatalf("expected valid query, got error: %v", err)
	}

	if err := NewSettingQuery().SetID("").SetLimit(10).Validate(); err == nil {
		t.Fatal("expected error for empty id")
	}

	if err := NewSettingQuery().SetKey("").Validate(); err == nil {
		t.Fatal("expected error for empty key")
	}

	if err := NewSettingQuery().SetKeyLike("").Validate(); err == nil {
		t.Fatal("expected error for empty key_like")
	}

	if err := NewSettingQuery().SetLimit(0).Validate(); err == nil {
		t.Fatal("expected error for non-positive limit")
	}

	if err := NewSettingQuery().SetOffset(-1).Validate(); err == nil {
		t.Fatal("expected error for negative offset")
	}

	if err := NewSettingQuery().SetOrderBy("").Validate(); err == nil {
		t.Fatal("expected error for empty order_by")
	}

	if err := NewSettingQuery().SetSortDirection("").Validate(); err == nil {
		t.Fatal("expected error for empty sort_direction")
	}

	if err := NewSettingQuery().SetCreatedAtGte("").Validate(); err == nil {
		t.Fatal("expected error for empty created_at_gte")
	}

	if err := NewSettingQuery().SetCreatedAtLte("").Validate(); err == nil {
		t.Fatal("expected error for empty created_at_lte")
	}
}

func TestSettingQueryFluentAccessors(t *testing.T) {
	query := NewSettingQuery().
		SetID("id-1").
		SetIDIn([]string{"a", "b"}).
		SetCreatedAtGte("2026-01-01 00:00:00").
		SetCreatedAtLte("2026-12-31 23:59:59").
		SetKey("currency").
		SetKeyLike("curr").
		SetLimit(5).
		SetOffset(10).
		SetOrderBy(COLUMN_SETTING_KEY).
		SetSortDirection("asc").
		SetCountOnly(true).
		SetSoftDeletedIncluded(true).
		SetColumns([]string{COLUMN_ID})

	if !query.HasID() || query.ID() != "id-1" {
		t.Fatal("expected id to be set")
	}
	if !query.HasIDIn() || len(query.IDIn()) != 2 {
		t.Fatal("expected id_in to be set")
	}
	if !query.HasCreatedAtGte() || query.CreatedAtGte() != "2026-01-01 00:00:00" {
		t.Fatal("expected created_at_gte to be set")
	}
	if !query.HasCreatedAtLte() || query.CreatedAtLte() != "2026-12-31 23:59:59" {
		t.Fatal("expected created_at_lte to be set")
	}
	if !query.HasKey() || query.Key() != "currency" {
		t.Fatal("expected key to be set")
	}
	if !query.HasKeyLike() || query.KeyLike() != "curr" {
		t.Fatal("expected key_like to be set")
	}
	if !query.HasLimit() || query.Limit() != 5 {
		t.Fatal("expected limit to be set")
	}
	if !query.HasOffset() || query.Offset() != 10 {
		t.Fatal("expected offset to be set")
	}
	if !query.HasOrderBy() || query.OrderBy() != COLUMN_SETTING_KEY {
		t.Fatal("expected order_by to be set")
	}
	if !query.HasSortDirection() || query.SortDirection() != "asc" {
		t.Fatal("expected sort_direction to be set")
	}
	if !query.IsCountOnly() {
		t.Fatal("expected count_only to be true")
	}
	if !query.SoftDeletedIncluded() {
		t.Fatal("expected soft_deleted_included to be true")
	}
	if len(query.Columns()) != 1 {
		t.Fatal("expected columns to be set")
	}
}
