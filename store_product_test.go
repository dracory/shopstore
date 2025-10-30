package shopstore

import (
	"context"
	"strings"
	"testing"

	"github.com/dracory/sb"
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

	productFound, errFind := store.ProductFindByID(ctx, product.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if productFound == nil {
		t.Fatal("Product MUST NOT be nil")
	}

	if productFound.Title() != "Ruler" {
		t.Fatal("Product title MUST BE 'Ruler', found: ", productFound.Title())
	}

	if productFound.Status() != PRODUCT_STATUS_DRAFT {
		t.Fatal("Product status MUST BE 'draft', found: ", productFound.Status())
	}

	if productFound.Quantity() != "1" {
		t.Fatal("Product quantity MUST BE '1', found: ", productFound.Quantity())
	}

	if productFound.Price() != "19.9900" {
		t.Fatal("Product price MUST BE '19.9900', found: ", productFound.Price())
	}

	if productFound.Memo() != "test ruler" {
		t.Fatal("Product memo MUST BE 'test ruler', found: ", productFound.Memo())
	}

	if productFound.Meta("color") != "green" {
		t.Fatal("Product color meta MUST BE 'green', found: ", productFound.Meta("color"))
	}

	if productFound.Meta("size") != "xxl" {
		t.Fatal("Product size meta MUST BE 'xxl', found: ", productFound.Meta("xxl"))
	}

	if !strings.Contains(productFound.SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Product MUST NOT be soft deleted", productFound.SoftDeletedAt())
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

	if product.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Product MUST NOT be soft deleted")
	}

	err = store.ProductSoftDelete(ctx, product)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	productFound, errFind := store.ProductFindByID(ctx, product.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if productFound != nil {
		t.Fatal("Product MUST be nil")
	}

	productFindWithDeleted, errFind := store.ProductList(ctx, NewProductQuery().
		SetID(product.ID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(productFindWithDeleted) < 1 {
		t.Fatal("Product list MUST NOT be empty")
		return
	}

	if strings.Contains(productFindWithDeleted[0].SoftDeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Product MUST be soft deleted", productFound.SoftDeletedAt())
	}

}
