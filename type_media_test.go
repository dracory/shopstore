package shopstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewMediaDefaults(t *testing.T) {
	media := NewMedia()
	if media == nil {
		t.Fatal("NewMedia returned nil")
	}

	if media.GetStatus() != MEDIA_STATUS_DRAFT {
		t.Fatalf("expected status %q, got %q", MEDIA_STATUS_DRAFT, media.GetStatus())
	}

	if media.GetTitle() != "" {
		t.Fatalf("expected empty title, got %q", media.GetTitle())
	}

	if media.GetDescription() != "" {
		t.Fatalf("expected empty description, got %q", media.GetDescription())
	}

	if media.GetMemo() != "" {
		t.Fatalf("expected empty memo, got %q", media.GetMemo())
	}

	if media.GetID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if media.GetCreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if media.GetUpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if media.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", sb.MAX_DATETIME, media.GetSoftDeletedAt())
	}

	metas, err := media.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if media.GetMeta("missing") != "" {
		t.Fatal("expected GetMeta for missing key to return empty string")
	}
}

func TestMediaDataTracking(t *testing.T) {
	media := &Media{}

	media.SetTitle("Title").
		SetDescription("Desc").
		SetMemo("Memo").
		SetEntityID("ENTITY").
		SetURL("https://example.com").
		SetStatus(MEDIA_STATUS_ACTIVE).
		SetType("image/png")

	media.SetSequence(42)

	data := media.Data()
	if data == nil {
		t.Fatal("expected Data to return initialized map")
	}

	expected := map[string]string{
		COLUMN_TITLE:       "Title",
		COLUMN_DESCRIPTION: "Desc",
		COLUMN_MEMO:        "Memo",
		COLUMN_ENTITY_ID:   "ENTITY",
		COLUMN_MEDIA_URL:   "https://example.com",
		COLUMN_STATUS:      MEDIA_STATUS_ACTIVE,
		COLUMN_MEDIA_TYPE:  "image/png",
		COLUMN_SEQUENCE:    "42",
	}

	for key, want := range expected {
		if got := data[key]; got != want {
			t.Fatalf("expected %s to be %q, got %q", key, want, got)
		}
	}

	changed := media.DataChanged()
	for key, want := range expected {
		if got := changed[key]; got != want {
			t.Fatalf("expected DataChanged to track %s as %q, got %q", key, want, got)
		}
	}

	media.MarkAsNotDirty()
	if len(media.DataChanged()) != 0 {
		t.Fatal("expected DataChanged to be empty after MarkAsNotDirty")
	}

	media.SetTitle("Updated")
	if media.DataChanged()[COLUMN_TITLE] != "Updated" {
		t.Fatalf("expected DataChanged to track updated title, got %q", media.DataChanged()[COLUMN_TITLE])
	}
}

func TestMediaCarbonHelpers(t *testing.T) {
	media := &Media{}

	createdAt := "2025-06-01 10:20:30"
	updatedAt := "2025-06-02 11:22:33"
	softDeletedAt := "2025-07-01 00:00:00"

	if _, ok := media.SetCreatedAt(createdAt).(*Media); !ok {
		t.Fatal("expected SetCreatedAt to return *Media")
	}
	if media.GetCreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, media.GetCreatedAt())
	}
	if media.GetCreatedAtCarbon().ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", media.GetCreatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := media.SetUpdatedAt(updatedAt).(*Media); !ok {
		t.Fatal("expected SetUpdatedAt to return *Media")
	}
	if media.GetUpdatedAtCarbon().ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", media.GetUpdatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := media.SetSoftDeletedAt(softDeletedAt).(*Media); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Media")
	}
	if media.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", media.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestMediaMetasRoundTrip(t *testing.T) {
	media := &Media{}

	if err := media.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := media.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", metas["alpha"])
	}

	if media.GetMeta("alpha") != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", media.GetMeta("alpha"))
	}
}

func TestMediaUpsertMetasMergesValues(t *testing.T) {
	media := &Media{}

	if err := media.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := media.MetasUpsert(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := media.GetMetas()
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

func TestMediaMetasInvalidJSON(t *testing.T) {
	media := NewMediaFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if _, err := media.GetMetas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if media.GetMeta("anything") != "" {
		t.Fatalf("expected GetMeta to return empty string on invalid JSON")
	}
}

func TestMediaMetasHandlesNullJSON(t *testing.T) {
	media := NewMediaFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := media.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := media.MetasUpsert(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := media.GetMeta("alpha"); got != "beta" {
		t.Fatalf("expected GetMeta helper to return %q, got %q", "beta", got)
	}
}

func TestMediaSetMetaConvenience(t *testing.T) {
	media := &Media{}

	if err := media.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if media.GetMeta("key") != "value" {
		t.Fatalf("expected GetMeta to return %q, got %q", "value", media.GetMeta("key"))
	}
}

func TestMediaSetterChainingAndGetters(t *testing.T) {
	media := &Media{}

	if _, ok := media.SetDescription("desc").(*Media); !ok {
		t.Fatal("expected SetDescription to return *Media")
	}
	if media.GetDescription() != "desc" {
		t.Fatalf("expected GetDescription getter to return %q, got %q", "desc", media.GetDescription())
	}

	if _, ok := media.SetMemo("memo").(*Media); !ok {
		t.Fatal("expected SetMemo to return *Media")
	}
	if media.GetMemo() != "memo" {
		t.Fatalf("expected GetMemo getter to return %q, got %q", "memo", media.GetMemo())
	}

	if _, ok := media.SetStatus(MEDIA_STATUS_DRAFT).(*Media); !ok {
		t.Fatal("expected SetStatus to return *Media")
	}
	if media.GetStatus() != MEDIA_STATUS_DRAFT {
		t.Fatalf("expected GetStatus getter to return %q, got %q", MEDIA_STATUS_DRAFT, media.GetStatus())
	}

	if _, ok := media.SetTitle("title").(*Media); !ok {
		t.Fatal("expected SetTitle to return *Media")
	}
	if media.GetTitle() != "title" {
		t.Fatalf("expected GetTitle getter to return %q, got %q", "title", media.GetTitle())
	}

	if _, ok := media.SetType("image/jpg").(*Media); !ok {
		t.Fatal("expected SetType to return *Media")
	}
	if media.GetType() != "image/jpg" {
		t.Fatalf("expected GetType getter to return %q, got %q", "image/jpg", media.GetType())
	}

	if _, ok := media.SetEntityID("entity").(*Media); !ok {
		t.Fatal("expected SetEntityID to return *Media")
	}
	if media.GetEntityID() != "entity" {
		t.Fatalf("expected GetEntityID getter to return %q, got %q", "entity", media.GetEntityID())
	}

	if _, ok := media.SetURL("https://example.com").(*Media); !ok {
		t.Fatal("expected SetURL to return *Media")
	}
	if media.GetURL() != "https://example.com" {
		t.Fatalf("expected GetURL getter to return %q, got %q", "https://example.com", media.GetURL())
	}

	if _, ok := media.SetSequence(7).(*Media); !ok {
		t.Fatal("expected SetSequence to return *Media")
	}
	if media.GetSequence() != 7 {
		t.Fatalf("expected GetSequence getter to return %d, got %d", 7, media.GetSequence())
	}
}

func TestMediaIsSoftDeleted(t *testing.T) {
	media := &Media{}

	if _, ok := media.SetSoftDeletedAt(sb.MAX_DATETIME).(*Media); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Media")
	}
	if media.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAt to be %q, got %q", sb.MAX_DATETIME, media.GetSoftDeletedAt())
	}
	if media.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAtCarbon to be %q, got %q", sb.MAX_DATETIME, media.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := media.SetSoftDeletedAt("2024-01-01 00:00:00").(*Media); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Media")
	}
	if media.GetSoftDeletedAt() != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAt to be updated, got %q", media.GetSoftDeletedAt())
	}
	if media.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAtCarbon to match updated value, got %q", media.GetSoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestMediaMetaRemove(t *testing.T) {
	media := &Media{}

	if err := media.SetMetas(map[string]string{"key1": "value1", "key2": "value2"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := media.MetaRemove("key1"); err != nil {
		t.Fatalf("unexpected error removing meta: %v", err)
	}

	if media.GetMeta("key1") != "" {
		t.Fatalf("expected key1 to be removed, got %q", media.GetMeta("key1"))
	}

	if media.GetMeta("key2") != "value2" {
		t.Fatalf("expected key2 to still exist, got %q", media.GetMeta("key2"))
	}
}

func TestMediaMetaRemoveNonExistent(t *testing.T) {
	media := &Media{}

	if err := media.SetMetas(map[string]string{"key1": "value1"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := media.MetaRemove("nonexistent"); err != nil {
		t.Fatalf("unexpected error removing non-existent meta: %v", err)
	}

	if media.GetMeta("key1") != "value1" {
		t.Fatalf("expected key1 to still exist, got %q", media.GetMeta("key1"))
	}
}

func TestMediaMetasRemove(t *testing.T) {
	media := &Media{}

	if err := media.SetMetas(map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := media.MetasRemove([]string{"key1", "key2"}); err != nil {
		t.Fatalf("unexpected error removing metas: %v", err)
	}

	if media.GetMeta("key1") != "" {
		t.Fatalf("expected key1 to be removed, got %q", media.GetMeta("key1"))
	}

	if media.GetMeta("key2") != "" {
		t.Fatalf("expected key2 to be removed, got %q", media.GetMeta("key2"))
	}

	if media.GetMeta("key3") != "value3" {
		t.Fatalf("expected key3 to still exist, got %q", media.GetMeta("key3"))
	}
}

func TestMediaMetasRemoveEmptySlice(t *testing.T) {
	media := &Media{}

	if err := media.SetMetas(map[string]string{"key1": "value1"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	if err := media.MetasRemove([]string{}); err != nil {
		t.Fatalf("unexpected error removing empty slice: %v", err)
	}

	if media.GetMeta("key1") != "value1" {
		t.Fatalf("expected key1 to still exist, got %q", media.GetMeta("key1"))
	}
}

func TestMediaMetaRemoveErrorPropagation(t *testing.T) {
	media := NewMediaFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if err := media.MetaRemove("key"); err == nil {
		t.Fatal("expected error when removing meta with invalid JSON")
	}
}

func TestMediaMetasRemoveErrorPropagation(t *testing.T) {
	media := NewMediaFromExistingData(map[string]string{
		COLUMN_METAS: "{invalid",
	})

	if err := media.MetasRemove([]string{"key"}); err == nil {
		t.Fatal("expected error when removing metas with invalid JSON")
	}
}
