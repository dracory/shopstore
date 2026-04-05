package shopstore

import (
	"testing"
)

// == VARIANT MATRIX SCHEMA TESTS =============================================

func TestProductSetVariantMatrixSchema(t *testing.T) {
	product := &Product{}

	schema := VariantMatrixSchema{
		Name:     "color",
		Required: true,
		Options:  []string{"red", "blue", "black"},
	}

	err := product.SetVariantMatrixSchema(schema)
	if err != nil {
		t.Fatalf("unexpected error setting variant matrix schema: %v", err)
	}

	if !product.HasVariantMatrixSchema() {
		t.Fatal("expected HasVariantMatrixSchema to return true")
	}

	retrievedSchema, err := product.GetVariantMatrixSchema()
	if err != nil {
		t.Fatalf("unexpected error getting variant matrix schema: %v", err)
	}

	if retrievedSchema.Name != "color" {
		t.Fatalf("expected schema name 'color', got %q", retrievedSchema.Name)
	}

	if !retrievedSchema.Required {
		t.Fatal("expected schema Required to be true")
	}

	if len(retrievedSchema.Options) != 3 {
		t.Fatalf("expected 3 options, got %d", len(retrievedSchema.Options))
	}
}

func TestProductSetVariantMatrixSchemaEmpty(t *testing.T) {
	product := &Product{}

	err := product.SetVariantMatrixSchema(VariantMatrixSchema{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty schema with no name should still count as having dimensions
	// because the JSON is "{}" not "null" or ""
	if !product.HasVariantMatrixSchema() {
		t.Fatal("expected HasVariantMatrixSchema to return true for empty struct")
	}
}

func TestProductGetVariantMatrixSchemaNull(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_VARIANT_MATRIX_SCHEMA: "null",
	})

	schema, err := product.GetVariantMatrixSchema()
	if err != nil {
		t.Fatalf("unexpected error with null JSON: %v", err)
	}

	if schema.Name != "" {
		t.Fatalf("expected empty schema name for null, got %q", schema.Name)
	}

	if product.HasVariantMatrixSchema() {
		t.Fatal("expected HasVariantMatrixSchema to return false for null")
	}
}

func TestProductGetVariantMatrixSchemaInvalidJSON(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_VARIANT_MATRIX_SCHEMA: "{invalid",
	})

	_, err := product.GetVariantMatrixSchema()
	if err == nil {
		t.Fatal("expected error when parsing invalid JSON")
	}
}

// == VARIANT MATRIX VALUES TESTS =============================================

func TestProductSetVariantMatrixValues(t *testing.T) {
	product := &Product{}

	values := map[string]string{
		"color": "red",
		"size":  "9",
	}

	err := product.SetVariantMatrixValues(values)
	if err != nil {
		t.Fatalf("unexpected error setting variant matrix values: %v", err)
	}

	retrievedValues, err := product.GetVariantMatrixValues()
	if err != nil {
		t.Fatalf("unexpected error getting variant matrix values: %v", err)
	}

	if retrievedValues["color"] != "red" {
		t.Fatalf("expected color='red', got %q", retrievedValues["color"])
	}

	if retrievedValues["size"] != "9" {
		t.Fatalf("expected size='9', got %q", retrievedValues["size"])
	}
}

func TestProductGetVariantMatrixValuesEmpty(t *testing.T) {
	product := &Product{}

	values, err := product.GetVariantMatrixValues()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(values) != 0 {
		t.Fatalf("expected empty values, got %d", len(values))
	}
}

func TestProductGetVariantMatrixValuesNull(t *testing.T) {
	product := NewProductFromExistingData(map[string]string{
		COLUMN_VARIANT_MATRIX_VALUES: "null",
	})

	values, err := product.GetVariantMatrixValues()
	if err != nil {
		t.Fatalf("unexpected error with null JSON: %v", err)
	}

	if len(values) != 0 {
		t.Fatalf("expected empty values for null, got %d", len(values))
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
		t.Fatal("expected IsParent to return false for product without schema")
	}

	err := product.SetVariantMatrixSchema(VariantMatrixSchema{Name: "color"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !product.IsParent() {
		t.Fatal("expected IsParent to return true after setting schema")
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
