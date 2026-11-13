package shopstore

import (
	"context"
	"testing"
)

func TestStoreSettingSetGetRoundTrip(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	value, err := store.SettingGet(ctx, "currency", "USD")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if value != "USD" {
		t.Fatalf("expected default %q, got %q", "USD", value)
	}

	if err := store.SettingSet(ctx, "currency", "GBP"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err = store.SettingGet(ctx, "currency", "USD")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if value != "GBP" {
		t.Fatalf("expected %q, got %q", "GBP", value)
	}

	// Update existing key
	if err := store.SettingSet(ctx, "currency", "EUR"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err = store.SettingGet(ctx, "currency", "USD")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if value != "EUR" {
		t.Fatalf("expected %q, got %q", "EUR", value)
	}

	count, err := store.SettingCount(ctx, NewSettingQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 setting, got %d", count)
	}
}

func TestStoreSettingCRUD(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	setting := NewSetting().SetKey("shop_name").SetValue("Test Shop")
	if err := store.SettingCreate(ctx, setting); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.SettingFindByKey(ctx, "shop_name")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found == nil {
		t.Fatal("expected setting to be found")
	}
	if found.GetValue() != "Test Shop" {
		t.Fatalf("expected value %q, got %q", "Test Shop", found.GetValue())
	}

	foundByID, err := store.SettingFindByID(ctx, setting.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if foundByID == nil || foundByID.GetKey() != "shop_name" {
		t.Fatal("expected setting to be found by id")
	}

	setting.SetValue("Renamed Shop")
	if err := store.SettingUpdate(ctx, setting); err != nil {
		t.Fatal("unexpected error:", err)
	}

	updated, err := store.SettingFindByKey(ctx, "shop_name")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if updated == nil {
		t.Fatal("expected updated setting to be found")
	}
	if updated.GetValue() != "Renamed Shop" {
		t.Fatalf("expected value %q, got %q", "Renamed Shop", updated.GetValue())
	}

	list, err := store.SettingList(ctx, NewSettingQuery().SetKeyLike("shop"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(list))
	}

	if err := store.SettingDeleteByKey(ctx, "shop_name"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	deleted, err := store.SettingFindByKey(ctx, "shop_name")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if deleted != nil {
		t.Fatal("expected setting to be deleted")
	}
}

func TestStoreSettingDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	setting := NewSetting().SetKey("del").SetValue("x")
	if err := store.SettingCreate(ctx, setting); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.SettingDelete(ctx, setting); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.SettingFindByID(ctx, setting.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found != nil {
		t.Fatal("expected setting to be deleted")
	}

	setting2 := NewSetting().SetKey("del2").SetValue("y")
	if err := store.SettingCreate(ctx, setting2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.SettingDeleteByID(ctx, setting2.GetID()); err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, _ := store.SettingCount(ctx, nil)
	if count != 0 {
		t.Fatalf("expected 0 settings, got %d", count)
	}
}

func TestStoreSettingTableNameAndListFilters(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store.SettingTableName() != "shop_setting" {
		t.Fatalf("expected table name %q, got %q", "shop_setting", store.SettingTableName())
	}

	ctx := context.Background()

	setting := NewSetting().SetKey("dated").SetValue("v")
	if err := store.SettingCreate(ctx, setting); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.SettingList(ctx, NewSettingQuery().
		SetCreatedAtGte("2000-01-01 00:00:00").
		SetCreatedAtLte("2999-01-01 00:00:00").
		SetOrderBy(COLUMN_SETTING_KEY).
		SetSortDirection("asc").
		SetLimit(10).
		SetOffset(0))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(list))
	}

	list, err = store.SettingList(ctx, NewSettingQuery().
		SetCreatedAtGte("2999-01-01 00:00:00"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 settings, got %d", len(list))
	}

	list, err = store.SettingList(ctx, NewSettingQuery().
		SetCreatedAtLte("2000-01-01 00:00:00"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 settings, got %d", len(list))
	}
}

func TestStoreSettingSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	setting := NewSetting().SetKey("temp").SetValue("x")
	if err := store.SettingCreate(ctx, setting); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.SettingSoftDeleteByID(ctx, setting.GetID()); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.SettingFindByKey(ctx, "temp")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found != nil {
		t.Fatal("expected soft-deleted setting to be excluded")
	}

	list, err := store.SettingList(ctx, NewSettingQuery().SetSoftDeletedIncluded(true))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 soft-deleted setting, got %d", len(list))
	}
}
