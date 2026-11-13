package shopstore

import (
	"context"
	"testing"
)

// TestProductCategoryCRUD verifies create/list/set/delete round-trips.
func TestProductCategoryCRUD(t *testing.T) {
	store, _ := initStore(":memory:")
	ctx := context.Background()

	p := NewProduct().SetTitle("P1").SetStatus(PRODUCT_STATUS_ACTIVE)
	if err := store.ProductCreate(ctx, p); err != nil {
		t.Fatalf("ProductCreate: %v", err)
	}
	c1 := NewCategory().SetTitle("C1").SetStatus(CATEGORY_STATUS_ACTIVE)
	c2 := NewCategory().SetTitle("C2").SetStatus(CATEGORY_STATUS_ACTIVE)
	for _, c := range []CategoryInterface{c1, c2} {
		if err := store.CategoryCreate(ctx, c); err != nil {
			t.Fatalf("CategoryCreate: %v", err)
		}
	}

	// Set two categories for the product
	if err := store.ProductCategorySet(ctx, p.GetID(), []string{c1.GetID(), c2.GetID()}); err != nil {
		t.Fatalf("ProductCategorySet: %v", err)
	}

	ids, err := store.ProductCategoryIDs(ctx, p.GetID())
	if err != nil {
		t.Fatalf("ProductCategoryIDs: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 category IDs, got %v", ids)
	}

	count, err := store.ProductCategoryCount(ctx, NewProductCategoryQuery().SetProductID(p.GetID()))
	if err != nil {
		t.Fatalf("ProductCategoryCount: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count=2, got %d", count)
	}

	// Replace with a single category — stale link must be removed
	if err := store.ProductCategorySet(ctx, p.GetID(), []string{c2.GetID()}); err != nil {
		t.Fatalf("ProductCategorySet replace: %v", err)
	}
	ids, err = store.ProductCategoryIDs(ctx, p.GetID())
	if err != nil {
		t.Fatalf("ProductCategoryIDs: %v", err)
	}
	if len(ids) != 1 || ids[0] != c2.GetID() {
		t.Fatalf("expected only c2 linked, got %v", ids)
	}

	// Clear all links
	if err := store.ProductCategorySet(ctx, p.GetID(), nil); err != nil {
		t.Fatalf("ProductCategorySet clear: %v", err)
	}
	ids, err = store.ProductCategoryIDs(ctx, p.GetID())
	if err != nil {
		t.Fatalf("ProductCategoryIDs: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected no links, got %v", ids)
	}

	// Direct create/find/delete
	link := NewProductCategory().SetProductID(p.GetID()).SetCategoryID(c1.GetID())
	if err := store.ProductCategoryCreate(ctx, link); err != nil {
		t.Fatalf("ProductCategoryCreate: %v", err)
	}
	found, err := store.ProductCategoryFindByID(ctx, link.GetID())
	if err != nil || found == nil {
		t.Fatalf("ProductCategoryFindByID: %v %v", found, err)
	}
	if found.GetProductID() != p.GetID() || found.GetCategoryID() != c1.GetID() {
		t.Fatalf("unexpected link: %v", found.Data())
	}
	if err := store.ProductCategoryDeleteByID(ctx, link.GetID()); err != nil {
		t.Fatalf("ProductCategoryDeleteByID: %v", err)
	}
	found, _ = store.ProductCategoryFindByID(ctx, link.GetID())
	if found != nil {
		t.Fatal("expected nil after delete")
	}
}
