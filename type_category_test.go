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

	if category.Status() != CATEGORY_STATUS_DRAFT {
		t.Fatalf("expected status %q, got %q", CATEGORY_STATUS_DRAFT, category.Status())
	}

	if category.ParentID() != "" {
		t.Fatalf("expected empty parent ID, got %q", category.ParentID())
	}

	if category.Description() != "" {
		t.Fatalf("expected empty description, got %q", category.Description())
	}

	if category.Memo() != "" {
		t.Fatalf("expected empty memo, got %q", category.Memo())
	}

	if category.ID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if category.CreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if category.UpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if category.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at to be %q, got %q", sb.MAX_DATETIME, category.SoftDeletedAt())
	}

	metas, err := category.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if category.Meta("missing") != "" {
		t.Fatal("expected Meta for missing key to return empty string")
	}
}

func TestCategorySetMetasRoundTrip(t *testing.T) {
	category := &Category{}

	if err := category.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := category.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if got := metas["alpha"]; got != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", got)
	}

	if got := category.Meta("alpha"); got != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", got)
	}
}

func TestCategoryUpsertMetasMergesValues(t *testing.T) {
	category := &Category{}

	if err := category.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := category.UpsertMetas(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := category.Metas()
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

	if _, err := category.Metas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if got := category.Meta("anything"); got != "" {
		t.Fatalf("expected Meta to return empty string on invalid JSON, got %q", got)
	}
}

func TestCategoryMetasHandlesNullJSON(t *testing.T) {
	category := NewCategoryFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := category.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := category.UpsertMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := category.Meta("alpha"); got != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", got)
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
	if category.CreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, category.CreatedAt())
	}
	createdCarbon := category.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Fatal("expected CreatedAtCarbon to return value")
	}
	if createdCarbon.ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", createdCarbon.ToDateTimeString(carbon.UTC))
	}

	if _, ok := category.SetUpdatedAt(updatedAt).(*Category); !ok {
		t.Fatal("expected SetUpdatedAt to return *Category")
	}
	updatedCarbon := category.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Fatal("expected UpdatedAtCarbon to return value")
	}
	if updatedCarbon.ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", updatedCarbon.ToDateTimeString(carbon.UTC))
	}

	if _, ok := category.SetSoftDeletedAt(softDeletedAt).(*Category); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Category")
	}
	softDeletedCarbon := category.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Fatal("expected SoftDeletedAtCarbon to return value")
	}
	if softDeletedCarbon.ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", softDeletedCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestCategorySetterChainingAndGetters(t *testing.T) {
	category := &Category{}

	if _, ok := category.SetDescription("desc").(*Category); !ok {
		t.Fatal("expected SetDescription to return *Category")
	}
	if category.Description() != "desc" {
		t.Fatalf("expected Description getter to return %q, got %q", "desc", category.Description())
	}

	if _, ok := category.SetMemo("memo").(*Category); !ok {
		t.Fatal("expected SetMemo to return *Category")
	}
	if category.Memo() != "memo" {
		t.Fatalf("expected Memo getter to return %q, got %q", "memo", category.Memo())
	}

	if _, ok := category.SetParentID("parent").(*Category); !ok {
		t.Fatal("expected SetParentID to return *Category")
	}
	if category.ParentID() != "parent" {
		t.Fatalf("expected ParentID getter to return %q, got %q", "parent", category.ParentID())
	}

	if _, ok := category.SetTitle("title").(*Category); !ok {
		t.Fatal("expected SetTitle to return *Category")
	}
	if category.Title() != "title" {
		t.Fatalf("expected Title getter to return %q, got %q", "title", category.Title())
	}

	if _, ok := category.SetStatus(CATEGORY_STATUS_ACTIVE).(*Category); !ok {
		t.Fatal("expected SetStatus to return *Category")
	}
	if category.Status() != CATEGORY_STATUS_ACTIVE {
		t.Fatalf("expected Status getter to return %q, got %q", CATEGORY_STATUS_ACTIVE, category.Status())
	}
}

func TestCategorySetMetaConvenience(t *testing.T) {
	category := &Category{}

	if err := category.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if got := category.Meta("key"); got != "value" {
		t.Fatalf("expected Meta to return %q, got %q", "value", got)
	}
}
