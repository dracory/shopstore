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

	if !strings.Contains(productFound.GetSoftDeletedAt(), sb.MAX_DATETIME) {
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

	if product.GetSoftDeletedAt() != sb.MAX_DATETIME {
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

	if strings.Contains(productFindWithDeleted[0].GetSoftDeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Product MUST be soft deleted", productFindWithDeleted[0].GetSoftDeletedAt())
	}
}
