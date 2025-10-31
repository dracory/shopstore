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

	if media.Status() != MEDIA_STATUS_DRAFT {
		t.Fatalf("expected status %q, got %q", MEDIA_STATUS_DRAFT, media.Status())
	}

	if media.Title() != "" {
		t.Fatalf("expected empty title, got %q", media.Title())
	}

	if media.Description() != "" {
		t.Fatalf("expected empty description, got %q", media.Description())
	}

	if media.Memo() != "" {
		t.Fatalf("expected empty memo, got %q", media.Memo())
	}

	if media.ID() == "" {
		t.Fatal("expected generated ID to be non-empty")
	}

	if media.CreatedAt() == "" {
		t.Fatal("expected created at to be set")
	}

	if media.UpdatedAt() == "" {
		t.Fatal("expected updated at to be set")
	}

	if media.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at %q, got %q", sb.MAX_DATETIME, media.SoftDeletedAt())
	}

	metas, err := media.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected no metas by default, got %v", metas)
	}

	if media.Meta("missing") != "" {
		t.Fatal("expected Meta for missing key to return empty string")
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
	if media.CreatedAt() != createdAt {
		t.Fatalf("expected CreatedAt to be %q, got %q", createdAt, media.CreatedAt())
	}
	if media.CreatedAtCarbon().ToDateTimeString(carbon.UTC) != createdAt {
		t.Fatalf("expected CreatedAtCarbon to match input, got %q", media.CreatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := media.SetUpdatedAt(updatedAt).(*Media); !ok {
		t.Fatal("expected SetUpdatedAt to return *Media")
	}
	if media.UpdatedAtCarbon().ToDateTimeString(carbon.UTC) != updatedAt {
		t.Fatalf("expected UpdatedAtCarbon to match input, got %q", media.UpdatedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := media.SetSoftDeletedAt(softDeletedAt).(*Media); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Media")
	}
	if media.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != softDeletedAt {
		t.Fatalf("expected SoftDeletedAtCarbon to match input, got %q", media.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}

func TestMediaMetasRoundTrip(t *testing.T) {
	media := &Media{}

	if err := media.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := media.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if metas["alpha"] != "beta" {
		t.Fatalf("expected meta to be %q, got %q", "beta", metas["alpha"])
	}

	if media.Meta("alpha") != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", media.Meta("alpha"))
	}
}

func TestMediaUpsertMetasMergesValues(t *testing.T) {
	media := &Media{}

	if err := media.SetMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error setting initial metas: %v", err)
	}

	if err := media.UpsertMetas(map[string]string{"alpha": "updated", "gamma": "delta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	metas, err := media.Metas()
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

	if _, err := media.Metas(); err == nil {
		t.Fatal("expected error when parsing invalid metas JSON")
	}

	if media.Meta("anything") != "" {
		t.Fatalf("expected Meta to return empty string on invalid JSON")
	}
}

func TestMediaMetasHandlesNullJSON(t *testing.T) {
	media := NewMediaFromExistingData(map[string]string{
		COLUMN_METAS: "null",
	})

	metas, err := media.Metas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}

	if len(metas) != 0 {
		t.Fatalf("expected empty metas map for null JSON, got %v", metas)
	}

	if err := media.UpsertMetas(map[string]string{"alpha": "beta"}); err != nil {
		t.Fatalf("unexpected error upserting metas: %v", err)
	}

	if got := media.Meta("alpha"); got != "beta" {
		t.Fatalf("expected Meta helper to return %q, got %q", "beta", got)
	}
}

func TestMediaSetMetaConvenience(t *testing.T) {
	media := &Media{}

	if err := media.SetMeta("key", "value"); err != nil {
		t.Fatalf("unexpected error from SetMeta: %v", err)
	}

	if media.Meta("key") != "value" {
		t.Fatalf("expected Meta to return %q, got %q", "value", media.Meta("key"))
	}
}

func TestMediaSetterChainingAndGetters(t *testing.T) {
	media := &Media{}

	if _, ok := media.SetDescription("desc").(*Media); !ok {
		t.Fatal("expected SetDescription to return *Media")
	}
	if media.Description() != "desc" {
		t.Fatalf("expected Description getter to return %q, got %q", "desc", media.Description())
	}

	if _, ok := media.SetMemo("memo").(*Media); !ok {
		t.Fatal("expected SetMemo to return *Media")
	}
	if media.Memo() != "memo" {
		t.Fatalf("expected Memo getter to return %q, got %q", "memo", media.Memo())
	}

	if _, ok := media.SetStatus(MEDIA_STATUS_DRAFT).(*Media); !ok {
		t.Fatal("expected SetStatus to return *Media")
	}
	if media.Status() != MEDIA_STATUS_DRAFT {
		t.Fatalf("expected Status getter to return %q, got %q", MEDIA_STATUS_DRAFT, media.Status())
	}

	if _, ok := media.SetTitle("title").(*Media); !ok {
		t.Fatal("expected SetTitle to return *Media")
	}
	if media.Title() != "title" {
		t.Fatalf("expected Title getter to return %q, got %q", "title", media.Title())
	}

	if _, ok := media.SetType("image/jpg").(*Media); !ok {
		t.Fatal("expected SetType to return *Media")
	}
	if media.Type() != "image/jpg" {
		t.Fatalf("expected Type getter to return %q, got %q", "image/jpg", media.Type())
	}

	if _, ok := media.SetEntityID("entity").(*Media); !ok {
		t.Fatal("expected SetEntityID to return *Media")
	}
	if media.EntityID() != "entity" {
		t.Fatalf("expected EntityID getter to return %q, got %q", "entity", media.EntityID())
	}

	if _, ok := media.SetURL("https://example.com").(*Media); !ok {
		t.Fatal("expected SetURL to return *Media")
	}
	if media.URL() != "https://example.com" {
		t.Fatalf("expected URL getter to return %q, got %q", "https://example.com", media.URL())
	}

	if _, ok := media.SetSequence(7).(*Media); !ok {
		t.Fatal("expected SetSequence to return *Media")
	}
	if media.Sequence() != 7 {
		t.Fatalf("expected Sequence getter to return %d, got %d", 7, media.Sequence())
	}
}

func TestMediaIsSoftDeleted(t *testing.T) {
	media := &Media{}

	if _, ok := media.SetSoftDeletedAt(sb.MAX_DATETIME).(*Media); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Media")
	}
	if media.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAt to be %q, got %q", sb.MAX_DATETIME, media.SoftDeletedAt())
	}
	if media.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != sb.MAX_DATETIME {
		t.Fatalf("expected SoftDeletedAtCarbon to be %q, got %q", sb.MAX_DATETIME, media.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}

	if _, ok := media.SetSoftDeletedAt("2024-01-01 00:00:00").(*Media); !ok {
		t.Fatal("expected SetSoftDeletedAt to return *Media")
	}
	if media.SoftDeletedAt() != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAt to be updated, got %q", media.SoftDeletedAt())
	}
	if media.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC) != "2024-01-01 00:00:00" {
		t.Fatalf("expected SoftDeletedAtCarbon to match updated value, got %q", media.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC))
	}
}
