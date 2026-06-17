package shopstore

import (
	"context"
	"strings"
	"testing"

	"github.com/dromara/carbon/v2"
	_ "modernc.org/sqlite"
)

func TestStoreCategoryCount(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	count, err := store.CategoryCount(ctx, NewCategoryQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	for i := 0; i < 3; i++ {
		cat := NewCategory().SetStatus(CATEGORY_STATUS_ACTIVE).SetTitle("CAT")
		if err := store.CategoryCreate(ctx, cat); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err = store.CategoryCount(ctx, NewCategoryQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 3 {
		t.Fatalf("expected count 3, got %d", count)
	}
}

func TestStoreDiscountCount(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	count, err := store.DiscountCount(ctx, NewDiscountQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	for i := 0; i < 2; i++ {
		d := NewDiscount().SetStatus(DISCOUNT_STATUS_ACTIVE).SetTitle("DISC")
		if err := store.DiscountCreate(ctx, d); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err = store.DiscountCount(ctx, NewDiscountQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestStoreMediaCount(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	count, err := store.MediaCount(ctx, NewMediaQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	for i := 0; i < 4; i++ {
		m := NewMedia().SetStatus(MEDIA_STATUS_ACTIVE).SetEntityID("E1").SetTitle("M").SetURL("http://x").SetType(MEDIA_TYPE_IMAGE_JPG).SetSequence(1)
		if err := store.MediaCreate(ctx, m); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err = store.MediaCount(ctx, NewMediaQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 4 {
		t.Fatalf("expected count 4, got %d", count)
	}
}

func TestStoreOrderCount(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	count, err := store.OrderCount(ctx, NewOrderQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	for i := 0; i < 5; i++ {
		o := NewOrder().SetStatus(ORDER_STATUS_PENDING).SetCustomerID("C1")
		if err := store.OrderCreate(ctx, o); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err = store.OrderCount(ctx, NewOrderQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 5 {
		t.Fatalf("expected count 5, got %d", count)
	}
}

func TestStoreOrderLineItemCount(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	count, err := store.OrderLineItemCount(ctx, NewOrderLineItemQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	for i := 0; i < 2; i++ {
		item := NewOrderLineItem().SetStatus(ORDER_STATUS_PENDING).SetOrderID("O1").SetProductID("P1").SetQuantityInt(1).SetPriceFloat(10)
		if err := store.OrderLineItemCreate(ctx, item); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err = store.OrderLineItemCount(ctx, NewOrderLineItemQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestStoreProductCount(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	count, err := store.ProductCount(ctx, NewProductQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	for i := 0; i < 3; i++ {
		p := NewProduct().SetStatus(PRODUCT_STATUS_ACTIVE).SetTitle("P").SetQuantityInt(1).SetPriceFloat(1)
		if err := store.ProductCreate(ctx, p); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err = store.ProductCount(ctx, NewProductQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 3 {
		t.Fatalf("expected count 3, got %d", count)
	}
}

func TestStoreProductDelete(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	product := NewProduct().SetStatus(PRODUCT_STATUS_DRAFT).SetTitle("Del").SetQuantityInt(1).SetPriceFloat(1)
	ctx := context.Background()
	if err := store.ProductCreate(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ProductDelete(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.ProductFindByID(ctx, product.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found != nil {
		t.Fatal("Product MUST be nil after hard delete")
	}
}

func TestStoreProductDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	product := NewProduct().SetStatus(PRODUCT_STATUS_DRAFT).SetTitle("Del").SetQuantityInt(1).SetPriceFloat(1)
	ctx := context.Background()
	if err := store.ProductCreate(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ProductDeleteByID(ctx, product.GetID()); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.ProductFindByID(ctx, product.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found != nil {
		t.Fatal("Product MUST be nil after hard delete")
	}
}

func TestStoreProductSoftDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	product := NewProduct().SetStatus(PRODUCT_STATUS_DRAFT).SetTitle("Del").SetQuantityInt(1).SetPriceFloat(1)
	ctx := context.Background()
	if err := store.ProductCreate(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ProductSoftDeleteByID(ctx, product.GetID()); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.ProductFindByID(ctx, product.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found != nil {
		t.Fatal("Product MUST be nil after soft delete")
	}
}

func TestStoreProductVariantList(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	parent := NewProduct().
		SetStatus(PRODUCT_STATUS_ACTIVE).
		SetTitle("Parent").
		SetQuantityInt(10).
		SetPriceFloat(100)
	if err := store.ProductCreate(ctx, parent); err != nil {
		t.Fatal("unexpected error:", err)
	}

	variant := NewProduct().
		SetStatus(PRODUCT_STATUS_ACTIVE).
		SetTitle("Variant").
		SetParentID(parent.GetID()).
		SetQuantityInt(5).
		SetPriceFloat(100)
	if err := store.ProductCreate(ctx, variant); err != nil {
		t.Fatal("unexpected error:", err)
	}

	variants, err := store.ProductVariantList(ctx, parent.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(variants) != 1 {
		t.Fatalf("expected 1 variant, got %d", len(variants))
	}
	if variants[0].GetID() != variant.GetID() {
		t.Fatal("unexpected variant id")
	}
}

func TestStoreProductIsParent(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	product := NewProduct().SetStatus(PRODUCT_STATUS_ACTIVE).SetTitle("P").SetQuantityInt(1).SetPriceFloat(1)
	if err := store.ProductCreate(ctx, product); err != nil {
		t.Fatal("unexpected error:", err)
	}

	isParent, err := store.ProductIsParent(ctx, product.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	// NewProduct sets an empty variant matrix schema, so HasVariantMatrixSchema returns true
	if !isParent {
		t.Fatal("expected true because NewProduct sets empty variant matrix schema")
	}
}

func TestStoreProductGetParent(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	parent := NewProduct().SetStatus(PRODUCT_STATUS_ACTIVE).SetTitle("Parent").SetQuantityInt(1).SetPriceFloat(1)
	if err := store.ProductCreate(ctx, parent); err != nil {
		t.Fatal("unexpected error:", err)
	}

	variant := NewProduct().SetStatus(PRODUCT_STATUS_ACTIVE).SetTitle("Var").SetParentID(parent.GetID()).SetQuantityInt(1).SetPriceFloat(1)
	if err := store.ProductCreate(ctx, variant); err != nil {
		t.Fatal("unexpected error:", err)
	}

	foundParent, err := store.ProductGetParent(ctx, variant.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if foundParent == nil {
		t.Fatal("parent MUST NOT be nil")
	}
	if foundParent.GetID() != parent.GetID() {
		t.Fatal("unexpected parent id")
	}

	_, err = store.ProductGetParent(ctx, parent.GetID())
	if err == nil {
		t.Fatal("expected error for non-variant product")
	}
}

func TestStoreDiscountFindByCode(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_ACTIVE).
		SetTitle("CODE_TEST").
		SetCode("SPECIAL10")
	if err := store.DiscountCreate(ctx, discount); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.DiscountFindByCode(ctx, "SPECIAL10")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found == nil {
		t.Fatal("discount MUST NOT be nil")
	}
	if found.GetCode() != "SPECIAL10" {
		t.Fatal("unexpected discount code")
	}

	notFound, err := store.DiscountFindByCode(ctx, "MISSING")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if notFound != nil {
		t.Fatal("expected nil for missing code")
	}
}

func TestStoreCategoryListWithOffsetOrderBy(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		cat := NewCategory().SetStatus(CATEGORY_STATUS_ACTIVE).SetTitle("CAT")
		if err := store.CategoryCreate(ctx, cat); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	list, err := store.CategoryList(ctx, NewCategoryQuery().
		SetStatus(CATEGORY_STATUS_ACTIVE).
		SetOffset(2).
		SetLimit(10).
		SetOrderBy(COLUMN_CREATED_AT).
		SetSortDirection("asc"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 items after offset, got %d", len(list))
	}
}

func TestStoreCategoryListWithTitleLike(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	cat1 := NewCategory().SetStatus(CATEGORY_STATUS_ACTIVE).SetTitle("Alpha Category")
	cat2 := NewCategory().SetStatus(CATEGORY_STATUS_ACTIVE).SetTitle("Beta Category")
	if err := store.CategoryCreate(ctx, cat1); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.CategoryCreate(ctx, cat2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.CategoryList(ctx, NewCategoryQuery().
		SetTitleLike("Alpha").
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
	if list[0].GetTitle() != "Alpha Category" {
		t.Fatal("unexpected title")
	}
}

func TestStoreCategoryListWithIDIn(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	cat1 := NewCategory().SetStatus(CATEGORY_STATUS_ACTIVE).SetTitle("C1")
	cat2 := NewCategory().SetStatus(CATEGORY_STATUS_ACTIVE).SetTitle("C2")
	if err := store.CategoryCreate(ctx, cat1); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.CategoryCreate(ctx, cat2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.CategoryList(ctx, NewCategoryQuery().
		SetIDIn([]string{cat1.GetID(), cat2.GetID()}).
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list))
	}
}

func TestStoreCategoryListWithStatusIn(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	cat1 := NewCategory().SetStatus(CATEGORY_STATUS_ACTIVE).SetTitle("C1")
	cat2 := NewCategory().SetStatus(CATEGORY_STATUS_DRAFT).SetTitle("C2")
	if err := store.CategoryCreate(ctx, cat1); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.CategoryCreate(ctx, cat2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.CategoryList(ctx, NewCategoryQuery().
		SetStatus(CATEGORY_STATUS_ACTIVE).
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
	if list[0].GetStatus() != CATEGORY_STATUS_ACTIVE {
		t.Fatal("expected active status")
	}
}

func TestStoreDiscountListWithCodeAndType(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	d1 := NewDiscount().SetStatus(DISCOUNT_STATUS_ACTIVE).SetTitle("D1").SetCode("CODE1").SetType(DISCOUNT_TYPE_AMOUNT).SetAmount(10)
	d2 := NewDiscount().SetStatus(DISCOUNT_STATUS_ACTIVE).SetTitle("D2").SetCode("CODE2").SetType(DISCOUNT_TYPE_PERCENT).SetAmount(5)
	if err := store.DiscountCreate(ctx, d1); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.DiscountCreate(ctx, d2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.DiscountList(ctx, NewDiscountQuery().
		SetCode("CODE1").
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}

	list, err = store.DiscountList(ctx, NewDiscountQuery().
		SetType(DISCOUNT_TYPE_PERCENT).
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
}

func TestStoreDiscountListWithDateRanges(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	d := NewDiscount().
		SetStatus(DISCOUNT_STATUS_ACTIVE).
		SetTitle("D1").
		SetStartsAt("2023-01-01 00:00:00").
		SetEndsAt("2023-12-31 23:59:59")
	if err := store.DiscountCreate(ctx, d); err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Use actual stored created_at value for query
	createdAt := d.GetCreatedAt()
	startsAt := d.GetStartsAt()
	endsAt := d.GetEndsAt()

	list, err := store.DiscountList(ctx, NewDiscountQuery().
		SetCreatedAtGte(createdAt).
		SetCreatedAtLte(createdAt).
		SetStartsAtGte(startsAt).
		SetStartsAtLte(startsAt).
		SetEndsAtGte(endsAt).
		SetEndsAtLte(endsAt).
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
}

func TestStoreOrderListWithCustomerIDAndStatusIn(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	o1 := NewOrder().SetStatus(ORDER_STATUS_PENDING).SetCustomerID("C1")
	o2 := NewOrder().SetStatus(ORDER_STATUS_COMPLETED).SetCustomerID("C1")
	o3 := NewOrder().SetStatus(ORDER_STATUS_PENDING).SetCustomerID("C2")
	if err := store.OrderCreate(ctx, o1); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.OrderCreate(ctx, o2); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.OrderCreate(ctx, o3); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.OrderList(ctx, NewOrderQuery().
		SetCustomerID("C1").
		SetStatusIn([]string{ORDER_STATUS_PENDING, ORDER_STATUS_COMPLETED}).
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list))
	}
}

func TestStoreOrderListWithCreatedAtRange(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	now := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	o := NewOrder().SetStatus(ORDER_STATUS_PENDING).SetCustomerID("C1")
	if err := store.OrderCreate(ctx, o); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.OrderList(ctx, NewOrderQuery().
		SetCreatedAtGte(now).
		SetCreatedAtLte(now).
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
}

func TestStoreOrderLineItemListWithOrderIDInAndProductID(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	item1 := NewOrderLineItem().SetStatus(ORDER_STATUS_PENDING).SetOrderID("O1").SetProductID("P1").SetQuantityInt(1).SetPriceFloat(10)
	item2 := NewOrderLineItem().SetStatus(ORDER_STATUS_PENDING).SetOrderID("O2").SetProductID("P1").SetQuantityInt(2).SetPriceFloat(20)
	if err := store.OrderLineItemCreate(ctx, item1); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.OrderLineItemCreate(ctx, item2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.OrderLineItemList(ctx, NewOrderLineItemQuery().
		SetOrderIDIn([]string{"O1", "O2"}).
		SetProductID("P1").
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list))
	}
}

func TestStoreMediaListWithEntityIDAndType(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	m1 := NewMedia().SetStatus(MEDIA_STATUS_ACTIVE).SetEntityID("E1").SetTitle("M1").SetURL("http://x").SetType(MEDIA_TYPE_IMAGE_JPG).SetSequence(1)
	m2 := NewMedia().SetStatus(MEDIA_STATUS_ACTIVE).SetEntityID("E1").SetTitle("M2").SetURL("http://y").SetType(MEDIA_TYPE_VIDEO_MP4).SetSequence(2)
	if err := store.MediaCreate(ctx, m1); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.MediaCreate(ctx, m2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	// First list all media to verify records exist
	all, err := store.MediaList(ctx, NewMediaQuery().SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 total media, got %d", len(all))
	}
	for i, m := range all {
		t.Logf("media %d: entity_id=%q type=%q status=%q soft_deleted=%q", i, m.GetEntityID(), m.GetType(), m.GetStatus(), m.GetSoftDeletedAt())
	}

	list, err := store.MediaList(ctx, NewMediaQuery().
		SetEntityID("E1").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
}

func TestStoreProductListWithTitleLike(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	p1 := NewProduct().SetStatus(PRODUCT_STATUS_ACTIVE).SetTitle("Awesome Widget").SetQuantityInt(1).SetPriceFloat(10)
	p2 := NewProduct().SetStatus(PRODUCT_STATUS_ACTIVE).SetTitle("Boring Widget").SetQuantityInt(1).SetPriceFloat(10)
	if err := store.ProductCreate(ctx, p1); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.ProductCreate(ctx, p2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.ProductList(ctx, NewProductQuery().
		SetTitleLike("Awesome").
		SetLimit(10))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
	if list[0].GetTitle() != "Awesome Widget" {
		t.Fatal("unexpected title")
	}
}

func TestNewStoreMissingCategoryTableName(t *testing.T) {
	db, err := initDB(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	_, err = NewStore(NewStoreOptions{
		DB: db,
	})
	if err == nil {
		t.Fatal("expected error for missing CategoryTableName")
	}
	if !strings.Contains(err.Error(), "CategoryTableName") {
		t.Fatal("expected error message to mention CategoryTableName")
	}
}

func TestNewStoreMissingDB(t *testing.T) {
	_, err := NewStore(NewStoreOptions{
		CategoryTableName:      "shop_category",
		DiscountTableName:      "shop_discount",
		MediaTableName:         "shop_media",
		OrderTableName:         "shop_order",
		OrderLineItemTableName: "shop_order_line_item",
		ProductTableName:       "shop_product",
	})
	if err == nil {
		t.Fatal("expected error for nil DB")
	}
	if !strings.Contains(err.Error(), "DB") {
		t.Fatal("expected error message to mention DB")
	}
}

func TestStoreEnableDebug(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Should not panic
	store.EnableDebug(true)
	store.EnableDebug(false)
}

func TestStoreMigrateDown(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	if err := store.MigrateDown(ctx); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.MigrateUp(ctx); err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreCategoryDeleteNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategoryDelete(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil category")
	}
}

func TestStoreDiscountDeleteNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.DiscountDelete(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil discount")
	}
}

func TestStoreMediaDeleteNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MediaDelete(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil media")
	}
}

func TestStoreOrderDeleteNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDelete(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil order")
	}
}

func TestStoreProductDeleteNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductDelete(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil product")
	}
}

func TestStoreProductDeleteByIDEmpty(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductDeleteByID(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestStoreCategorySoftDeleteNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategorySoftDelete(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil category")
	}
}

func TestStoreDiscountSoftDeleteNil(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.DiscountSoftDelete(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil discount")
	}
}
