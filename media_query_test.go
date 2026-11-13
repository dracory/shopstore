package shopstore

import (
	"context"
	"testing"
)

func TestMediaQueryValidate(t *testing.T) {
	query := NewMediaQuery()
	if err := query.Validate(); err != nil {
		t.Fatalf("expected valid query, got error: %v", err)
	}

	if err := NewMediaQuery().SetID("").Validate(); err == nil {
		t.Fatal("expected error for empty id")
	}

	if err := NewMediaQuery().SetIDIn([]string{}).Validate(); err == nil {
		t.Fatal("expected error for empty id_in")
	}

	if err := NewMediaQuery().SetEntityID("").Validate(); err == nil {
		t.Fatal("expected error for empty entity_id")
	}

	if err := NewMediaQuery().SetStatus("").Validate(); err == nil {
		t.Fatal("expected error for empty status")
	}

	if err := NewMediaQuery().SetTitleLike("").Validate(); err == nil {
		t.Fatal("expected error for empty title_like")
	}

	if err := NewMediaQuery().SetOrderBy("").Validate(); err == nil {
		t.Fatal("expected error for empty order_by")
	}

	if err := NewMediaQuery().SetSortDirection("").Validate(); err == nil {
		t.Fatal("expected error for empty sort_direction")
	}

	if err := NewMediaQuery().SetOffset(-1).Validate(); err == nil {
		t.Fatal("expected error for negative offset")
	}

	if err := NewMediaQuery().SetLimit(-1).Validate(); err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestMediaQueryFluentAccessors(t *testing.T) {
	query := NewMediaQuery().
		SetID("id-1").
		SetIDIn([]string{"a", "b"}).
		SetEntityID("entity-1").
		SetStatus(MEDIA_STATUS_ACTIVE).
		SetTitleLike("photo").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetLimit(5).
		SetOffset(10).
		SetOrderBy(COLUMN_TITLE).
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
	if !query.HasEntityID() || query.EntityID() != "entity-1" {
		t.Fatal("expected entity_id to be set")
	}
	if !query.HasStatus() || query.Status() != MEDIA_STATUS_ACTIVE {
		t.Fatal("expected status to be set")
	}
	if !query.HasTitleLike() || query.TitleLike() != "photo" {
		t.Fatal("expected title_like to be set")
	}
	if !query.HasType() || query.Type() != MEDIA_TYPE_IMAGE_JPG {
		t.Fatal("expected type to be set")
	}
	if !query.HasLimit() || query.Limit() != 5 {
		t.Fatal("expected limit to be set")
	}
	if !query.HasOffset() || query.Offset() != 10 {
		t.Fatal("expected offset to be set")
	}
	if !query.HasOrderBy() || query.OrderBy() != COLUMN_TITLE {
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

func TestStoreMediaListFilters(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	media := NewMedia().
		SetTitle("Summer photo").
		SetStatus(MEDIA_STATUS_ACTIVE).
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetURL("https://example.com/summer.jpg").
		SetSequence(1).
		SetEntityID("entity-1")
	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatal("unexpected error:", err)
	}

	other := NewMedia().
		SetTitle("Winter video").
		SetStatus(MEDIA_STATUS_DRAFT).
		SetType(MEDIA_TYPE_VIDEO_MP4).
		SetURL("https://example.com/winter.mp4").
		SetSequence(2).
		SetEntityID("entity-2")
	if err := store.MediaCreate(ctx, other); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.MediaList(ctx, NewMediaQuery().SetStatus(MEDIA_STATUS_ACTIVE))
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 active media, got %d (err=%v)", len(list), err)
	}

	list, err = store.MediaList(ctx, NewMediaQuery().SetTitleLike("Summer"))
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 media matching title, got %d (err=%v)", len(list), err)
	}

	list, err = store.MediaList(ctx, NewMediaQuery().SetIDIn([]string{media.GetID(), other.GetID()}).
		SetType(MEDIA_TYPE_VIDEO_MP4))
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 video media, got %d (err=%v)", len(list), err)
	}

	list, err = store.MediaList(ctx, NewMediaQuery().
		SetOrderBy(COLUMN_TITLE).
		SetSortDirection("asc").
		SetLimit(1).
		SetOffset(1))
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 media on second page, got %d (err=%v)", len(list), err)
	}

	count, err := store.MediaCount(ctx, NewMediaQuery().SetEntityID("entity-1"))
	if err != nil || count != 1 {
		t.Fatalf("expected count 1, got %d (err=%v)", count, err)
	}
}
