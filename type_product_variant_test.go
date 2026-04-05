package shopstore

import (
	"testing"
)

// == VARIANT DIMENSION TESTS =================================================

func TestProductSetVariantDimensionsSimpleArray(t *testing.T) {
	product := &Product{}

	err := product.SetVariantDimensions([]string{"color", "size"})
	if err != nil {
		t.Fatalf("unexpected error setting variant dimensions: %v", err)
	}

	if !product.HasVariantDimensions() {
		t.Fatal("expected HasVariantDimensions to return true")
	}

	names, err := product.GetVariantDimensionNames()
	if err != nil {
		t.Fatalf("unexpected error getting dimension names: %v", err)
	}

	if len(names) != 2 {
		t.Fatalf("expected 2 dimension names, got %d", len(names))
	}

	if names[0] != "color" || names[1] != "size" {
		t.Fatalf("expected [color, size], got %v", names)
	}
}

func TestProductSetVariantDimensionsStructured(t *testing.T) {
	product := &Product{}

	dims := []VariantDimension{
		{Name: "color", Required: true, Options: []string{"red", "blue", "black"}},
		{Name: "size", Required: true, Options: []string{"8", "9", "10", "11"}},
	}

	err := product.SetVariantDimensions(dims)
	if err != nil {
		t.Fatalf("unexpected error setting variant dimensions: %v", err)
	}

	retrievedDims, err := product.GetVariantDimensions()
	if err != nil {
		t.Fatalf("unexpected error getting variant dimensions: %v", err)
	}

	if len(retrievedDims) != 2 {
		t.Fatalf("expected 2 dimensions, got %d", len(retrievedDims))
	}

	if retrievedDims[0].Name != "color" || !retrievedDims[0].Required {
		t.Fatalf("expected first dimension to be color with required=true, got %+v", retrievedDims[0])
	}

	if len(retrievedDims[0].Options) != 3 {
		t.Fatalf("expected 3 color options, got %d", len(retrievedDims[0].Options))
	}
}

func TestProductSetVariantDimensionsOnVariant(t *testing.T) {
	product := &Product{}
	product.SetParentID("parent123")

	err := product.SetVariantDimensions([]string{"color"})
	if err == nil {
		t.Fatal("expected error when setting dimensions on a variant")
	}
}

func TestProductHasVariantDimensionsEmpty(t *testing.T) {
	product := &Product{}

	if product.HasVariantDimensions() {
		t.Fatal("expected HasVariantDimensions to return false for empty product")
	}

	dims, err := product.GetVariantDimensions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dims) != 0 {
		t.Fatalf("expected empty dimensions, got %d", len(dims))
	}

	names, err := product.GetVariantDimensionNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(names) != 0 {
		t.Fatalf("expected empty dimension names, got %d", len(names))
	}
}

func TestProductGetVariantDimensionsNull(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_VARIANT_MATRIX_SCHEMA: "null",
	})

	dims, err := product.GetVariantDimensions()
	if err != nil {
		t.Fatalf("unexpected error with null JSON: %v", err)
	}

	if len(dims) != 0 {
		t.Fatalf("expected empty dimensions for null JSON, got %d", len(dims))
	}

	if product.HasVariantDimensions() {
		t.Fatal("expected HasVariantDimensions to return false for null")
	}
}

func TestProductGetVariantDimensionsInvalidJSON(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_VARIANT_MATRIX_SCHEMA: "{invalid",
	})

	_, err := product.GetVariantDimensions()
	if err == nil {
		t.Fatal("expected error when parsing invalid JSON")
	}

	_, err = product.GetVariantDimensionNames()
	if err == nil {
		t.Fatal("expected error when getting dimension names with invalid JSON")
	}
}

// == PARENT/VARIANT RELATIONSHIP TESTS =======================================

func TestProductIsVariant(t *testing.T) {
	product := &Product{}

	if product.IsVariant() {
		t.Fatal("expected IsVariant to return false for product without parent")
	}

	product.SetParentID("parent123")

	if !product.IsVariant() {
		t.Fatal("expected IsVariant to return true after setting parent ID")
	}
}

func TestProductIsParent(t *testing.T) {
	product := &Product{}

	if product.IsParent() {
		t.Fatal("expected IsParent to return false for product without dimensions")
	}

	err := product.SetVariantDimensions([]string{"color", "size"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !product.IsParent() {
		t.Fatal("expected IsParent to return true after setting dimensions")
	}
}

func TestProductGetParentID(t *testing.T) {
	product := &Product{}

	if product.GetParentID() != "" {
		t.Fatalf("expected empty parent ID by default, got %q", product.GetParentID())
	}

	product.SetParentID("parent456")

	if product.GetParentID() != "parent456" {
		t.Fatalf("expected parent ID parent456, got %q", product.GetParentID())
	}
}

func TestProductSetParentIDFluentInterface(t *testing.T) {
	product := &Product{}

	result := product.SetParentID("parent789")
	if result != product {
		t.Fatal("expected SetParentID to return the same product for fluent interface")
	}
}
