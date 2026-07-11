package shopstore

import (
	"context"
	"errors"
	"testing"
)

// -----------------------------------------------------------------------
// Order consistency checks
// -----------------------------------------------------------------------

func TestOrderDelete_BlockedByLineItems(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUST1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderCreate(ctx, order); err != nil {
		t.Fatal("unexpected error:", err)
	}

	item := NewOrderLineItem().
		SetOrderID(order.GetID()).
		SetProductID("PROD1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderLineItemCreate(ctx, item); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDelete(ctx, order)
	if !errors.Is(err, ErrOrderHasActiveLineItems) {
		t.Fatalf("expected ErrOrderHasActiveLineItems, got: %v", err)
	}

	found, errFind := store.OrderFindByID(ctx, order.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("order must still exist after blocked delete")
	}
}

func TestOrderDeleteByID_BlockedByLineItems(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUST1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderCreate(ctx, order); err != nil {
		t.Fatal("unexpected error:", err)
	}

	item := NewOrderLineItem().
		SetOrderID(order.GetID()).
		SetProductID("PROD1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderLineItemCreate(ctx, item); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDeleteByID(ctx, order.GetID())
	if !errors.Is(err, ErrOrderHasActiveLineItems) {
		t.Fatalf("expected ErrOrderHasActiveLineItems, got: %v", err)
	}

	found, errFind := store.OrderFindByID(ctx, order.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("order must still exist after blocked delete")
	}
}

func TestOrderDelete_BlockedByMedia(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUST1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderCreate(ctx, order); err != nil {
		t.Fatal("unexpected error:", err)
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID(order.GetID()).
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/img.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDelete(ctx, order)
	if !errors.Is(err, ErrOrderHasActiveMedia) {
		t.Fatalf("expected ErrOrderHasActiveMedia, got: %v", err)
	}

	found, errFind := store.OrderFindByID(ctx, order.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("order must still exist after blocked delete")
	}
}

func TestOrderSoftDelete_BlockedByLineItems(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUST1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderCreate(ctx, order); err != nil {
		t.Fatal("unexpected error:", err)
	}

	item := NewOrderLineItem().
		SetOrderID(order.GetID()).
		SetProductID("PROD1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderLineItemCreate(ctx, item); err != nil {
		t.Fatal("unexpected error:", err)
	}

	originalSoftDeletedAt := order.GetSoftDeletedAt()

	err = store.OrderSoftDelete(ctx, order)
	if !errors.Is(err, ErrOrderHasActiveLineItems) {
		t.Fatalf("expected ErrOrderHasActiveLineItems, got: %v", err)
	}

	if order.GetSoftDeletedAt() != originalSoftDeletedAt {
		t.Fatal("SoftDeletedAt must be unchanged when consistency check fails")
	}

	found, errFind := store.OrderFindByID(ctx, order.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("order must still exist after blocked soft delete")
	}
}

func TestOrderSoftDeleteByID_BlockedByLineItems(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUST1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderCreate(ctx, order); err != nil {
		t.Fatal("unexpected error:", err)
	}

	item := NewOrderLineItem().
		SetOrderID(order.GetID()).
		SetProductID("PROD1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderLineItemCreate(ctx, item); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderSoftDeleteByID(ctx, order.GetID())
	if !errors.Is(err, ErrOrderHasActiveLineItems) {
		t.Fatalf("expected ErrOrderHasActiveLineItems, got: %v", err)
	}

	found, errFind := store.OrderFindByID(ctx, order.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("order must still exist after blocked soft delete")
	}
}

func TestOrderDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent order, got:", err)
	}
}

func TestOrderSoftDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderSoftDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent order, got:", err)
	}
}

// -----------------------------------------------------------------------
// Product consistency checks
// -----------------------------------------------------------------------

func TestProductDelete_BlockedByVariants(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	parent := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Parent Product").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, parent); err != nil {
		t.Fatal("unexpected error:", err)
	}

	variant := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Variant Product").
		SetParentID(parent.GetID()).
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, variant); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductDelete(ctx, parent)
	if !errors.Is(err, ErrProductHasActiveVariants) {
		t.Fatalf("expected ErrProductHasActiveVariants, got: %v", err)
	}

	found, errFind := store.ProductFindByID(ctx, parent.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("product must still exist after blocked delete")
	}
}

func TestProductDelete_BlockedByLineItems(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Product").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	item := NewOrderLineItem().
		SetOrderID("ORDER1").
		SetProductID(product.GetID()).
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderLineItemCreate(ctx, item); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductDelete(ctx, product)
	if !errors.Is(err, ErrProductHasActiveLineItems) {
		t.Fatalf("expected ErrProductHasActiveLineItems, got: %v", err)
	}

	found, errFind := store.ProductFindByID(ctx, product.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("product must still exist after blocked delete")
	}
}

func TestProductDelete_BlockedByMedia(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Product").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID(product.GetID()).
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/img.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductDelete(ctx, product)
	if !errors.Is(err, ErrProductHasActiveMedia) {
		t.Fatalf("expected ErrProductHasActiveMedia, got: %v", err)
	}

	found, errFind := store.ProductFindByID(ctx, product.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("product must still exist after blocked delete")
	}
}

func TestProductSoftDelete_BlockedByVariants(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	parent := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Parent Product").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, parent); err != nil {
		t.Fatal("unexpected error:", err)
	}

	variant := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Variant Product").
		SetParentID(parent.GetID()).
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, variant); err != nil {
		t.Fatal("unexpected error:", err)
	}

	originalSoftDeletedAt := parent.GetSoftDeletedAt()

	err = store.ProductSoftDelete(ctx, parent)
	if !errors.Is(err, ErrProductHasActiveVariants) {
		t.Fatalf("expected ErrProductHasActiveVariants, got: %v", err)
	}

	if parent.GetSoftDeletedAt() != originalSoftDeletedAt {
		t.Fatal("SoftDeletedAt must be unchanged when consistency check fails")
	}

	found, errFind := store.ProductFindByID(ctx, parent.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("product must still exist after blocked soft delete")
	}
}

func TestProductSoftDeleteByID_BlockedByVariants(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	parent := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Parent Product").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, parent); err != nil {
		t.Fatal("unexpected error:", err)
	}

	variant := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Variant Product").
		SetParentID(parent.GetID()).
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, variant); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductSoftDeleteByID(ctx, parent.GetID())
	if !errors.Is(err, ErrProductHasActiveVariants) {
		t.Fatalf("expected ErrProductHasActiveVariants, got: %v", err)
	}

	found, errFind := store.ProductFindByID(ctx, parent.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("product must still exist after blocked soft delete")
	}
}

func TestProductDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent product, got:", err)
	}
}

func TestProductSoftDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductSoftDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent product, got:", err)
	}
}

// -----------------------------------------------------------------------
// Category consistency checks
// -----------------------------------------------------------------------

func TestCategoryDelete_BlockedByChildren(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	parent := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("Parent Category")

	if err := store.CategoryCreate(ctx, parent); err != nil {
		t.Fatal("unexpected error:", err)
	}

	child := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("Child Category").
		SetParentID(parent.GetID())

	if err := store.CategoryCreate(ctx, child); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategoryDelete(ctx, parent)
	if !errors.Is(err, ErrCategoryHasActiveChildren) {
		t.Fatalf("expected ErrCategoryHasActiveChildren, got: %v", err)
	}

	found, errFind := store.CategoryFindByID(ctx, parent.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("category must still exist after blocked delete")
	}
}

func TestCategoryDelete_BlockedByMedia(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("Category")

	if err := store.CategoryCreate(ctx, category); err != nil {
		t.Fatal("unexpected error:", err)
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID(category.GetID()).
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/img.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategoryDelete(ctx, category)
	if !errors.Is(err, ErrCategoryHasActiveMedia) {
		t.Fatalf("expected ErrCategoryHasActiveMedia, got: %v", err)
	}

	found, errFind := store.CategoryFindByID(ctx, category.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("category must still exist after blocked delete")
	}
}

func TestCategorySoftDelete_BlockedByChildren(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	parent := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("Parent Category")

	if err := store.CategoryCreate(ctx, parent); err != nil {
		t.Fatal("unexpected error:", err)
	}

	child := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("Child Category").
		SetParentID(parent.GetID())

	if err := store.CategoryCreate(ctx, child); err != nil {
		t.Fatal("unexpected error:", err)
	}

	originalSoftDeletedAt := parent.GetSoftDeletedAt()

	err = store.CategorySoftDelete(ctx, parent)
	if !errors.Is(err, ErrCategoryHasActiveChildren) {
		t.Fatalf("expected ErrCategoryHasActiveChildren, got: %v", err)
	}

	if parent.GetSoftDeletedAt() != originalSoftDeletedAt {
		t.Fatal("SoftDeletedAt must be unchanged when consistency check fails")
	}

	found, errFind := store.CategoryFindByID(ctx, parent.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("category must still exist after blocked soft delete")
	}
}

func TestCategorySoftDeleteByID_BlockedByChildren(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	parent := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("Parent Category")

	if err := store.CategoryCreate(ctx, parent); err != nil {
		t.Fatal("unexpected error:", err)
	}

	child := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("Child Category").
		SetParentID(parent.GetID())

	if err := store.CategoryCreate(ctx, child); err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategorySoftDeleteByID(ctx, parent.GetID())
	if !errors.Is(err, ErrCategoryHasActiveChildren) {
		t.Fatalf("expected ErrCategoryHasActiveChildren, got: %v", err)
	}

	found, errFind := store.CategoryFindByID(ctx, parent.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found == nil {
		t.Fatal("category must still exist after blocked soft delete")
	}
}

func TestCategoryDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategoryDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent category, got:", err)
	}
}

func TestCategorySoftDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategorySoftDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent category, got:", err)
	}
}

// -----------------------------------------------------------------------
// Media delete / soft delete idempotency
// -----------------------------------------------------------------------

func TestMediaDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MediaDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent media, got:", err)
	}
}

func TestMediaSoftDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MediaSoftDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent media, got:", err)
	}
}

// -----------------------------------------------------------------------
// Discount delete / soft delete idempotency
// -----------------------------------------------------------------------

func TestDiscountSoftDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.DiscountSoftDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent discount, got:", err)
	}
}

func TestDiscountSoftDeleteByID_EmptyIDReturnsError(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.DiscountSoftDeleteByID(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty discount id, got nil")
	}
}

// -----------------------------------------------------------------------
// OrderLineItem soft delete idempotency
// -----------------------------------------------------------------------

func TestOrderLineItemSoftDeleteByID_NotFoundReturnsNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderLineItemSoftDeleteByID(context.Background(), "nonexistent-id")
	if err != nil {
		t.Fatal("expected nil error for non-existent order line item, got:", err)
	}
}

func TestOrderLineItemSoftDeleteByID_EmptyIDReturnsError(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderLineItemSoftDeleteByID(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty order line item id, got nil")
	}
}

// -----------------------------------------------------------------------
// Normal delete succeeds when no active children exist
// -----------------------------------------------------------------------

func TestOrderDelete_SucceedsWhenNoChildren(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUST1").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.OrderCreate(ctx, order); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.OrderDelete(ctx, order); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, errFind := store.OrderFindByID(ctx, order.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found != nil {
		t.Fatal("order must be nil after successful delete")
	}
}

func TestProductDelete_SucceedsWhenNoChildren(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Product").
		SetQuantityInt(1).
		SetPriceFloat(10.00)

	if err := store.ProductCreate(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ProductDelete(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, errFind := store.ProductFindByID(ctx, product.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found != nil {
		t.Fatal("product must be nil after successful delete")
	}
}

func TestCategoryDelete_SucceedsWhenNoChildren(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("Category")

	if err := store.CategoryCreate(ctx, category); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.CategoryDelete(ctx, category); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, errFind := store.CategoryFindByID(ctx, category.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if found != nil {
		t.Fatal("category must be nil after successful delete")
	}
}
