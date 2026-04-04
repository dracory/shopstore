package shopstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewCategoryDefaults(t *testing.T) {
	category := NewCategory()
	if category == nil {
		t.Fatal("NewCategory returned nil")
	}

	if category.GetStatus() != CATEGORY_STATUS_DRAFT {
		t.Fatalf("expected status %q, got %q", CATEGORY_STATUS_DRAFT, category.GetStatus())
	}

	if category.GetParentID() != "" {
		t.Fatalf("expected empty parent ID, got %q", category.GetParentID())
	}

	if category.GetDescription() != "" {
		t.Fatalf("expected empty description, got %q", category.GetDescription())
	}

	if category.GetMemo() != "" {
		t.Fatalf("expected empty memo, got %q", category.GetMemo())
	}

	if category.GetID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if category.GetCreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if category.GetUpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if category.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at to be %q, got %q", sb.MAX_DATETIME, category.GetSoftDeletedAt())
	}

	metas, err := category.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if category.GetMeta("missing") != "" {
		t.Fatal("expected GetMeta for missing key to return empty string")
	}
}

func TestCategorySetMetasRoundTrip(t *testing.T) {
	category := &Category{}

	if err := category.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := category.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if got := metas["alpha"]; got != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", got)
	}

	if got := category.GetMeta("alpha"); got != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", got)
	}
}

func TestCategoryUpsertMetasMergesValues(t *testing.T) {
	category := &Category{}

	if err := category.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := category.MetasUpsert(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := category.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if got := metas["alpha"]; got != "updated" {
		t.Fatalf("expected alpha meta to be updated, got %q", got)
	}

	if got := metas["gamma"]; got != "delta" {
		t.Fatalf("expected gamma meta to be present, got %q", got)
	}
}

func TestCategoryMetaHandlesInvalidJSON(t *testing.T) {
	category := NewCategoryFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if _, err := category.GetMetas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if got := category.GetMeta("anything"); got != "" {
		t.Fatalf("expected GetMeta to return empty string on invalid JSON, got %q", got)
	}
}

func TestCategoryMetasHandlesNullJSON(t *testing.T) {
	category := NewCategoryFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := category.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := category.MetasUpsert(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := category.GetMeta("alpha"); got != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", got)
	}
}

func TestCategoryStatusPredicates(t *testing.T) {
	category := &Category{}

	if _, ok := category.SetStatus(CATEGORY_STATUS_ACTIVE).(*Category); !ok {
		t.Fatal("expected SetStatus to return *Category")
	}

	if !category.IsActive() {
		t.Fatal("expected category to be active")
	}

	if category.IsDraft() {
		t.Fatal("expected category not to be draft when active")
	}

	if _, ok := category.SetStatus(CATEGORY_STATUS_DRAFT).(*Category); !ok {
		t.Fatal("expected SetStatus to return *Category")
	}

	if !category.IsDraft() {
		t.Fatal("expected category to be draft")
	}

	if _, ok := category.SetStatus(CATEGORY_STATUS_INACTIVE).(*Category); !ok {
		t.Fatal("expected SetStatus to return *Category")
	}

	if !category.IsInactive() {
		t.Fatal("expected category to be inactive")
	}
}

func TestCategoryIsSoftDeleted(t *testing.T) {
	category := &Category{}

	if _, ok := category.SetSoftDeletedAt(sb.MAX_DATETIME).(*Category); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Category")
	}

	if category.IsSoftDeleted() {
		t.Fatal("expected category not to be soft deleted when timestamp is MAX_DATETIME")
	}

	if _, ok := category.SetSoftDeletedAt("2024-01-01 00:00:00").(*Category); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Category")
	}

	if !category.IsSoftDeleted() {
		t.Fatal("expected category to be soft deleted when timestamp differs from MAX_DATETIME")
	}
}

func TestCategoryDataTracking(t *testing.T) {
	category := &Category{}

	category.SetTitle("Title").
		SetDescription("Desc").
		SetMemo("Memo").
		SetParentID("parent")

	data := category.Data()
	if data == nil {
		t.Fatal("expected Data to return initialized map")
	}

	if data[COLUMN_TITLE] != "Title" {
		t.Fatalf("expected title to be %q, got %q", "Title", data[COLUMN_TITLE])
	}

	changed := category.DataChanged()
	if len(changed) != 4 {
		t.Fatalf("expected four changed entries, got %d", len(changed))
	}

	category.MarkAsNotDirty()
	if len(category.DataChanged()) != 0 {
		t.Fatal("expected DataChanged to be empty after MarkAsNotDirty")
	}

	category.SetTitle("Updated")
	if category.DataChanged()[COLUMN_TITLE] != "Updated" {
		t.Fatalf("expected DataChanged to track updated title, got %q", category.DataChanged()[COLUMN_TITLE])
	}
}

func TestCategoryCarbonHelpers(t *testing.T) {
	category := &Category{}

	createdAt := "2025-10-31 12:34:56"
	updatedAt := "2025-11-01 08:15:30"
	softDeletedAt := "2025-12-01 00:00:00"

	if _, ok := category.SetCreatedAt(createdAt).(*Category); !ok {
		t.Fatal("expected SetCreatedAt to return *Category")
	}
	if category.GetCreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, category.GetCreatedAt())
	}
	createdCarbon := category.GetCreatedAtCarbon()
	if createdCarbon == nil {
		t.Fatal("expected GetCreatedAtCarbon to return value")
	}
	if createdCarbon.ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected GetCreatedAtCarbon to match input, got %q", createdCarbon.ToDateTimeString(carbon.UTC))
	}

	if _, ok := category.SetUpdatedAt(updatedAt).(*Category); !ok {
		t.Fatal("expected SetUpdatedAt to return *Category")
	}
	updatedCarbon := category.GetUpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Fatal("expected GetUpdatedAtCarbon to return value")
	}
	if updatedCarbon.ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected GetUpdatedAtCarbon to match input, got %q", updatedCarbon.ToDateTimeString(carbon.UTC))
	}

	if _, ok := category.SetSoftDeletedAt(softDeletedAt).(*Category); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Category")
	}
	softDeletedCarbon := category.GetSoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Fatal("expected GetSoftDeletedAtCarbon to return value")
	}
	if softDeletedCarbon.ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected GetSoftDeletedAtCarbon to match input, got %q", softDeletedCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestCategorySetterChainingAndGetters(t *testing.T) {
	category := &Category{}

	if _, ok := category.SetDescription("desc").(*Category); !ok {
		t.Fatal("expected SetDescription to return *Category")
	}
	if category.GetDescription() != "desc" {
		t.Fatalf("expected GetDescription getter to return %q, got %q", "desc", category.GetDescription())
	}

	if _, ok := category.SetMemo("memo").(*Category); !ok {
		t.Fatal("expected SetMemo to return *Category")
	}
	if category.GetMemo() != "memo" {
		t.Fatalf("expected GetMemo getter to return %q, got %q", "memo", category.GetMemo())
	}

	if _, ok := category.SetParentID("parent").(*Category); !ok {
		t.Fatal("expected SetParentID to return *Category")
	}
	if category.GetParentID() != "parent" {
		t.Fatalf("expected GetParentID getter to return %q, got %q", "parent", category.GetParentID())
	}

	if _, ok := category.SetTitle("title").(*Category); !ok {
		t.Fatal("expected SetTitle to return *Category")
	}
	if category.GetTitle() != "title" {
		t.Fatalf("expected GetTitle getter to return %q, got %q", "title", category.GetTitle())
	}

	if _, ok := category.SetStatus(CATEGORY_STATUS_ACTIVE).(*Category); !ok {
		t.Fatal("expected SetStatus to return *Category")
	}
	if category.GetStatus() != CATEGORY_STATUS_ACTIVE {
		t.Fatalf("expected GetStatus getter to return %q, got %q", CATEGORY_STATUS_ACTIVE, category.GetStatus())
	}
}

func TestCategorySetMetaConvenience(t *testing.T) {
	category := &Category{}

	if err := category.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if got := category.GetMeta("key"); got != "value" {
		t.Fatalf("expected GetMeta to return %q, got %q", "value", got)
	}
}

func TestCategoryMetaRemove(t *testing.T) {
	category := &Category{}

	if err := category.SetMetas(map[string]string{"key1": "value1", "key2": "value2"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := category.MetaRemove("key1"); err != nil {
		t.Fatalf("unexpected error removing meta: %v", err)
	}

	if category.GetMeta("key1") != "" {
		t.Fatalf("expected key1 to be removed, got %q", category.GetMeta("key1"))
	}

	if category.GetMeta("key2") != "value2" {
		t.Fatalf("expected key2 to still exist, got %q", category.GetMeta("key2"))
	}
}

func TestCategoryMetaRemoveNonExistent(t *testing.T) {
	category := &Category{}

	if err := category.SetMetas(map[string]string{"key1": "value1"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := category.MetaRemove("nonexistent"); err != nil {
		t.Fatalf("unexpected error removing non-existent meta: %v", err)
	}

	if category.GetMeta("key1") != "value1" {
		t.Fatalf("expected key1 to still exist, got %q", category.GetMeta("key1"))
	}
}

func TestCategoryMetasRemove(t *testing.T) {
	category := &Category{}

	if err := category.SetMetas(map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := category.MetasRemove([]string{"key1", "key2"}); err != nil {
		t.Fatalf("unexpected error removing metas: %v", err)
	}

	if category.GetMeta("key1") != "" {
		t.Fatalf("expected key1 to be removed, got %q", category.GetMeta("key1"))
	}

	if category.GetMeta("key2") != "" {
		t.Fatalf("expected key2 to be removed, got %q", category.GetMeta("key2"))
	}

	if category.GetMeta("key3") != "value3" {
		t.Fatalf("expected key3 to still exist, got %q", category.GetMeta("key3"))
	}
}

func TestCategoryMetasRemoveEmptySlice(t *testing.T) {
	category := &Category{}

	if err := category.SetMetas(map[string]string{"key1": "value1"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := category.MetasRemove([]string{}); err != nil {
		t.Fatalf("unexpected error removing empty slice: %v", err)
	}

	if category.GetMeta("key1") != "value1" {
		t.Fatalf("expected key1 to still exist, got %q", category.GetMeta("key1"))
	}
}

func TestCategoryMetaRemoveErrorPropagation(t *testing.T) {
	category := NewCategoryFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if err := category.MetaRemove("key"); err == nil {
		t.Fatal("expected error when removing meta with invalid JSON")
	}
}

func TestCategoryMetasRemoveErrorPropagation(t *testing.T) {
	category := NewCategoryFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if err := category.MetasRemove([]string{"key"}); err == nil {
		t.Fatal("expected error when removing metas with invalid JSON")
	}
}
