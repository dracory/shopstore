package shopstore

import (
	"context"
	"strings"
	"testing"

	"github.com/dracory/neat"
)

func TestStoreProductCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = product.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.ProductCreate(ctx, product)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreProductFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Ruler").
		SetQuantityInt(1).
		SetPriceFloat(19.99).
		SetMemo("test ruler")

	err = product.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.ProductCreate(ctx, product)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	productFound, errFind := store.ProductFindByID(ctx, product.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if productFound == nil {
		t.Fatal("Product MUST NOT be nil")
	}

	if productFound.GetTitle() != "Ruler" {
		t.Fatal("Product title MUST BE 'Ruler', found: ", productFound.GetTitle())
	}

	if productFound.GetStatus() != PRODUCT_STATUS_DRAFT {
		t.Fatal("Product status MUST BE 'draft', found: ", productFound.GetStatus())
	}

	if productFound.GetQuantity() != "1" {
		t.Fatal("Product quantity MUST BE '1', found: ", productFound.GetQuantity())
	}

	if productFound.GetPrice() != "19.99" {
		t.Fatal("Product price MUST BE '19.99', found: ", productFound.GetPrice())
	}

	if productFound.GetMemo() != "test ruler" {
		t.Fatal("Product memo MUST BE 'test ruler', found: ", productFound.GetMemo())
	}

	if productFound.GetMeta("color") != "green" {
		t.Fatal("Product color meta MUST BE 'green', found: ", productFound.GetMeta("color"))
	}

	if productFound.GetMeta("size") != "xxl" {
		t.Fatal("Product size meta MUST BE 'xxl', found: ", productFound.GetMeta("xxl"))
	}

	if !strings.Contains(productFound.GetSoftDeletedAt(), MAX_DATETIME) {
		t.Fatal("Product MUST NOT be soft deleted", productFound.GetSoftDeletedAt())
	}
}

func TestStoreProductSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Ruler").
		SetQuantityInt(1).
		SetPriceFloat(19.99).
		SetMemo("test ruler")

	ctx := context.Background()
	err = store.ProductCreate(ctx, product)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if product.GetSoftDeletedAt() != MAX_DATETIME {
		t.Fatal("Product MUST NOT be soft deleted")
	}

	err = store.ProductSoftDelete(ctx, product)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	productFound, errFind := store.ProductFindByID(ctx, product.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if productFound != nil {
		t.Fatal("Product MUST be nil")
	}

	productFindWithDeleted, errFind := store.ProductList(ctx, NewProductQuery().
		SetID(product.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(productFindWithDeleted) < 1 {
		t.Fatal("Product list MUST NOT be empty")
		return
	}

	if strings.Contains(productFindWithDeleted[0].GetSoftDeletedAt(), neat.NullDateTime) {
		t.Fatal("Product MUST be soft deleted", productFindWithDeleted[0].GetSoftDeletedAt())
	}
}

func TestStoreProductList_ExcludeIDs(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create 3 products
	p1 := NewProduct().SetTitle("Product 1")
	p2 := NewProduct().SetTitle("Product 2")
	p3 := NewProduct().SetTitle("Product 3")

	err = store.ProductCreate(ctx, p1)
	if err != nil {
		t.Fatal(err)
	}
	err = store.ProductCreate(ctx, p2)
	if err != nil {
		t.Fatal(err)
	}
	err = store.ProductCreate(ctx, p3)
	if err != nil {
		t.Fatal(err)
	}

	// Test IDNotIn: exclude p1 and p2
	list, err := store.ProductList(ctx, NewProductQuery().SetIDNotIn([]string{p1.GetID(), p2.GetID()}))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 product, got %d", len(list))
	}
	if list[0].GetID() != p3.GetID() {
		t.Errorf("expected product %s, got %s", p3.GetID(), list[0].GetID())
	}
}

func TestStoreProductList_MetasIn(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create products with different meta values
	p1 := NewProduct().SetTitle("Product 1")
	err = p1.SetMetas(map[string]string{"is_featured": "featured", "color": "green"})
	if err != nil {
		t.Fatal(err)
	}

	p2 := NewProduct().SetTitle("Product 2")
	err = p2.SetMetas(map[string]string{"is_featured": "not_featured", "color": "green"})
	if err != nil {
		t.Fatal(err)
	}

	p3 := NewProduct().SetTitle("Product 3")
	err = p3.SetMetas(map[string]string{"is_featured": "featured", "color": "red"})
	if err != nil {
		t.Fatal(err)
	}

	err = store.ProductCreate(ctx, p1)
	if err != nil {
		t.Fatal(err)
	}
	err = store.ProductCreate(ctx, p2)
	if err != nil {
		t.Fatal(err)
	}
	err = store.ProductCreate(ctx, p3)
	if err != nil {
		t.Fatal(err)
	}

	// Test MetasIn: filter by single meta key-value (is_featured=featured)
	list, err := store.ProductList(ctx, NewProductQuery().SetMetasIn(map[string]string{
		"is_featured": "featured",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 featured products, got %d", len(list))
	}

	// Test MetasIn: filter by multiple meta key-values (is_featured=featured AND color=green)
	list, err = store.ProductList(ctx, NewProductQuery().SetMetasIn(map[string]string{
		"is_featured": "featured",
		"color":       "green",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 product matching both filters, got %d", len(list))
	}
	if list[0].GetID() != p1.GetID() {
		t.Errorf("expected product %s, got %s", p1.GetID(), list[0].GetID())
	}
}

func TestStoreProductList_MetasNotIn(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create products with different meta values
	p1 := NewProduct().SetTitle("Product 1")
	err = p1.SetMetas(map[string]string{"is_featured": "featured", "color": "green"})
	if err != nil {
		t.Fatal(err)
	}

	p2 := NewProduct().SetTitle("Product 2")
	err = p2.SetMetas(map[string]string{"is_featured": "not_featured", "color": "green"})
	if err != nil {
		t.Fatal(err)
	}

	p3 := NewProduct().SetTitle("Product 3")
	err = p3.SetMetas(map[string]string{"is_featured": "featured", "color": "red"})
	if err != nil {
		t.Fatal(err)
	}

	// p4 has no is_featured key at all — must be included in MetasNotIn results
	p4 := NewProduct().SetTitle("Product 4")
	err = p4.SetMetas(map[string]string{"color": "blue"})
	if err != nil {
		t.Fatal(err)
	}

	err = store.ProductCreate(ctx, p1)
	if err != nil {
		t.Fatal(err)
	}
	err = store.ProductCreate(ctx, p2)
	if err != nil {
		t.Fatal(err)
	}
	err = store.ProductCreate(ctx, p3)
	if err != nil {
		t.Fatal(err)
	}
	err = store.ProductCreate(ctx, p4)
	if err != nil {
		t.Fatal(err)
	}

	// Test MetasNotIn: exclude products where is_featured=featured
	// Should return p2 (not_featured) and p4 (key absent → NULL != value is NULL → must be included)
	list, err := store.ProductList(ctx, NewProductQuery().SetMetasNotIn(map[string]string{
		"is_featured": "featured",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 products (not featured + missing key), got %d", len(list))
	}
	gotIDs := map[string]bool{}
	for _, p := range list {
		gotIDs[p.GetID()] = true
	}
	if !gotIDs[p2.GetID()] {
		t.Errorf("expected product %s in results", p2.GetID())
	}
	if !gotIDs[p4.GetID()] {
		t.Errorf("expected product %s (missing key) in results", p4.GetID())
	}
}
