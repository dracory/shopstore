package shopstore

import "testing"

func TestNewSettingDefaults(t *testing.T) {
	setting := NewSetting()
	if setting == nil {
		t.Fatal("NewSetting returned nil")
	}

	if setting.GetID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if setting.GetKey() != "" {
		t.Fatalf("expected empty key, got %q", setting.GetKey())
	}

	if setting.GetValue() != "" {
		t.Fatalf("expected empty value, got %q", setting.GetValue())
	}

	if setting.GetSoftDeletedAt() != MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", MAX_DATETIME, setting.GetSoftDeletedAt())
	}

	if setting.GetCreatedAt() == "" {
		t.Fatal("expected created at to be non-empty")
	}

	if setting.GetUpdatedAt() == "" {
		t.Fatal("expected updated at to be non-empty")
	}
}

func TestSettingSettersAndGetters(t *testing.T) {
	setting := NewSetting()

	setting.SetKey("currency").SetValue("GBP")

	if setting.GetKey() != "currency" {
		t.Fatalf("expected key %q, got %q", "currency", setting.GetKey())
	}

	if setting.GetValue() != "GBP" {
		t.Fatalf("expected value %q, got %q", "GBP", setting.GetValue())
	}

	setting.SetCreatedAt("2026-01-01 10:00:00")
	if setting.GetCreatedAt() != "2026-01-01 10:00:00" {
		t.Fatalf("expected created at %q, got %q", "2026-01-01 10:00:00", setting.GetCreatedAt())
	}
	if setting.GetCreatedAtCarbon() == nil {
		t.Fatal("expected created at carbon to be non-nil")
	}

	setting.SetUpdatedAt("2026-01-02 10:00:00")
	if setting.GetUpdatedAt() != "2026-01-02 10:00:00" {
		t.Fatalf("expected updated at %q, got %q", "2026-01-02 10:00:00", setting.GetUpdatedAt())
	}
	if setting.GetUpdatedAtCarbon() == nil {
		t.Fatal("expected updated at carbon to be non-nil")
	}

	setting.SetSoftDeletedAt("2026-01-03 10:00:00")
	if setting.GetSoftDeletedAt() != "2026-01-03 10:00:00" {
		t.Fatalf("expected soft deleted at %q, got %q", "2026-01-03 10:00:00", setting.GetSoftDeletedAt())
	}
	if setting.GetSoftDeletedAtCarbon() == nil {
		t.Fatal("expected soft deleted at carbon to be non-nil")
	}
	if !setting.IsSoftDeleted() {
		t.Fatal("expected IsSoftDeleted to be true")
	}
}

func TestNewSettingFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "test-id",
		COLUMN_SETTING_KEY:     "currency",
		COLUMN_SETTING_VALUE:   "EUR",
		COLUMN_CREATED_AT:      "2026-01-01 00:00:00",
		COLUMN_UPDATED_AT:      "2026-01-01 00:00:00",
		COLUMN_SOFT_DELETED_AT: MAX_DATETIME,
	}

	setting := NewSettingFromExistingData(data)
	if setting == nil {
		t.Fatal("NewSettingFromExistingData returned nil")
	}

	if setting.GetID() != "test-id" {
		t.Fatalf("expected id %q, got %q", "test-id", setting.GetID())
	}
	if setting.GetKey() != "currency" {
		t.Fatalf("expected key %q, got %q", "currency", setting.GetKey())
	}
	if setting.GetValue() != "EUR" {
		t.Fatalf("expected value %q, got %q", "EUR", setting.GetValue())
	}
}
