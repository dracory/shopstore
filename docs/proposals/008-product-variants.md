# Proposal: Product Variants System

## Status: Draft

## Overview

Support multiple versions of a product (variants) - different sizes, colors, materials, etc. Each variant has its own SKU, price, stock level, while sharing common attributes with the parent product.

## Problem Statement

Current system only supports single-version products:
- Cannot track inventory per size/color
- Single SKU per product
- No variant-specific pricing (sale price for red version only)
- Cannot display "Small / Red / Cotton" as one purchasable option

## Options Analysis

### Option 1: Parent-Child Products (Recommended)

Each variant is a full Product with `ParentID` reference.

```go
type Product struct {
    // ... existing fields ...
    ParentID string  // Empty = parent product, Set = variant
}

// Examples:
// Parent: "Nike Air Max" (not buyable, display only)
//   - Child 1: "Nike Air Max - Size 9 - Red" (buyable)
//   - Child 2: "Nike Air Max - Size 10 - Blue" (buyable)
```

**Pros:**
- Simple - uses existing Product entity
- Each variant has full Product capabilities (price, stock, SKU, images)
- Easy querying: `store.ProductList(ctx, NewProductQuery().SetParentID(parentID))`
- Variants can have their own metas, media, tags
- Minimal schema changes (add one column)

**Cons:**
- Data duplication (title, description repeated)
- No enforced consistency across variants
- Must ensure only children appear in cart/checkout

**API Example:**
```go
// Create parent (display-only)
parent := NewProduct().
    SetTitle("Nike Air Max").
    SetStatus(PRODUCT_STATUS_PARENT) // Not purchasable

// Create variants
variant1 := NewProduct().
    SetTitle("Nike Air Max").
    SetParentID(parent.GetID()).
    SetSKU("NAM-RED-9").
    SetPriceFloat(129.99).
    SetQuantityInt(15).
    SetStatus(PRODUCT_STATUS_ACTIVE).
    SetMeta("size", "9").
    SetMeta("color", "red")

// Get all variants
variants, _ := store.ProductList(ctx, NewProductQuery().SetParentID(parentID))
```

---

### Option 2: Separate ProductVariant Entity

New table dedicated to variants, Product stores common data.

```go
type ProductVariant struct {
    ID          string
    ProductID   string
    SKU         string
    Price       float64
    Quantity    int64
    OptionValues map[string]string  // {"size": "9", "color": "red"}
    Status      string
    CreatedAt   string
    UpdatedAt   string
}
```

**Pros:**
- No data duplication (common data in Product)
- Enforced structure for variant options
- Clean separation of concerns
- Can enforce variant option consistency

**Cons:**
- Complex implementation (new entity, query builder, store methods)
- Variants lack full Product features (no direct metas, media, tags)
- More joins required for reads
- Breaking change to existing product workflows

**New Files:**
- `type_product_variant.go`
- `query_product_variant.go`
- `store_product_variant.go`
- Plus tests for all

---

### Option 3: Options Matrix (Attribute Combinations)

Product defines options, system generates variants.

```go
type Product struct {
    // ... existing ...
    Options []ProductOption  // [{"name": "Size", "values": ["S", "M", "L"]}, {"name": "Color", "values": ["Red", "Blue"]}]
    Variants []ProductVariant  // Auto-generated: S-Red, S-Blue, M-Red, M-Blue, etc.
}
```

**Pros:**
- Most user-friendly for merchants
- Auto-generates all combinations
- Enforces no gaps in matrix

**Cons:**
- Most complex implementation
- Risk of combinatorial explosion (5 sizes × 5 colors × 3 materials = 75 variants)
- Requires UI for managing options matrix
- Breaking change to Product structure

---

### Option 4: Metas-Based (Not Recommended)

Store variants in product metas as JSON.

```go
product.SetMeta("variants", `[{"sku": "RED-9", "size": "9", "color": "red", "qty": 5}]`)
```

**Pros:**
- Zero schema changes
- Simple to implement

**Cons:**
- Cannot query variants in SQL (e.g., "find all red products in stock")
- No database constraints (duplicates, invalid data)
- Manual JSON management
- No individual variant IDs for order line items
- Does not solve the actual problem

**Verdict:** ❌ Not viable for production

## Recommendation

**Option 1: Parent-Child Products**

Best balance of simplicity and functionality:
- Uses existing Product infrastructure
- Variants are first-class products (full feature set)
- Single column addition (`parent_id`)
- Backwards compatible (existing products have empty ParentID)
- Can be extended later (add variant-specific fields)

## Implementation Details

### Database Changes

```sql
-- Add to shop_product table
ALTER TABLE shop_product ADD COLUMN parent_id VARCHAR(255) DEFAULT '';
ALTER TABLE shop_product ADD COLUMN variant_dimensions JSON DEFAULT NULL;

CREATE INDEX idx_product_parent ON shop_product(parent_id);

-- Optional: enforce that only leaf products can be purchased
-- (parent products should not appear in cart)
```

### Entity Changes

```go
// type_product.go
func (p *Product) GetParentID() string {
    return p.Get(COLUMN_PARENT_ID)
}

func (p *Product) SetParentID(parentID string) ProductInterface {
    p.Set(COLUMN_PARENT_ID, parentID)
    return p
}

// IsVariant returns true if this product has a parent
func (p *Product) IsVariant() bool {
    return p.GetParentID() != ""
}

// IsParent returns true if this product is a parent (has children)
func (p *Product) IsParent() bool {
    // Requires store lookup, see below
    return false // placeholder
}

// GetVariants returns child products (requires store access)
// func (p *Product) GetVariants(store StoreInterface) ([]ProductInterface, error)
```

### Interface Updates

```go
type ProductInterface interface {
    // ... existing methods ...
    
    GetParentID() string
    SetParentID(parentID string) ProductInterface
    IsVariant() bool
}
```

### Store Methods

```go
type StoreInterface interface {
    // ... existing methods ...
    
    // Variant operations
    ProductVariantList(ctx context.Context, parentID string) ([]ProductInterface, error)
    ProductIsParent(ctx context.Context, productID string) (bool, error)
    ProductGetParent(ctx context.Context, productID string) (ProductInterface, error)
    
    // Query by parent
    // (Add SetParentID to ProductQuery)
}
```

### Query Builder Update

```go
type ProductQueryInterface interface {
    // ... existing ...
    
    SetParentID(parentID string) ProductQueryInterface
    SetParentIDNull() ProductQueryInterface // Only parent products
}
```

### Constants

```go
const COLUMN_PARENT_ID = "parent_id"

// Optional status for parent products
const PRODUCT_STATUS_PARENT = "parent" // Display only, not purchasable
```

### Variant Dimensions via JSON Column

To enforce consistency across variants and support the options matrix, store dimension definitions in a dedicated JSON column:

**Database Schema:**
```sql
-- Add variant_dimensions column (stores dimension schema for parents)
ALTER TABLE shop_product ADD COLUMN variant_dimensions JSON DEFAULT NULL;
```

**Entity Changes:**

```go
// type_product.go

// SetVariantDimensions defines which attributes variants must have
// Example: ["color", "size"] or [{"name":"color","required":true}]
func (p *Product) SetVariantDimensions(dims interface{}) error {
    if p.GetParentID() != "" {
        return errors.New("cannot set dimensions on a variant")
    }
    jsonBytes, err := json.Marshal(dims)
    if err != nil {
        return err
    }
    p.Set(COLUMN_VARIANT_DIMENSIONS, string(jsonBytes))
    return nil
}

// GetVariantDimensions returns the dimension schema
func (p *Product) GetVariantDimensions() ([]VariantDimension, error) {
    dimJSON := p.Get(COLUMN_VARIANT_DIMENSIONS)
    if dimJSON == "" || dimJSON == "null" {
        return []VariantDimension{}, nil
    }
    var dims []VariantDimension
    err := json.Unmarshal([]byte(dimJSON), &dims)
    return dims, err
}

// GetVariantDimensionNames returns just the dimension names
func (p *Product) GetVariantDimensionNames() ([]string, error) {
    dimJSON := p.Get(COLUMN_VARIANT_DIMENSIONS)
    if dimJSON == "" || dimJSON == "null" {
        return []string{}, nil
    }
    
    // Try simple string array first
    var simple []string
    if err := json.Unmarshal([]byte(dimJSON), &simple); err == nil {
        return simple, nil
    }
    
    // Try structured format
    var structured []VariantDimension
    if err := json.Unmarshal([]byte(dimJSON), &structured); err != nil {
        return nil, err
    }
    
    names := make([]string, len(structured))
    for i, d := range structured {
        names[i] = d.Name
    }
    return names, nil
}

// HasVariantDimensions returns true if dimensions are defined
func (p *Product) HasVariantDimensions() bool {
    dimJSON := p.Get(COLUMN_VARIANT_DIMENSIONS)
    return dimJSON != "" && dimJSON != "null"
}

// VariantDimension represents a single dimension configuration
type VariantDimension struct {
    Name     string   `json:"name"`      // "color", "size"
    Required bool     `json:"required"`  // must variant have this?
    Options  []string `json:"options,omitempty"` // allowed values (optional)
}
```

**Constants:**
```go
const COLUMN_VARIANT_DIMENSIONS = "variant_dimensions"
```

**Why separate column:**
- Clean separation from user-defined metas
- No risk of overwrite by `SetMetas()`
- Can extend matrix data structure later
- Queryable: `WHERE variant_dimensions IS NOT NULL` finds all parents with defined schemas

**Validation in Store:**

```go
func (s *Store) ProductCreate(ctx context.Context, product ProductInterface) error {
    // If this is a variant, validate against parent's dimensions
    if product.GetParentID() != "" {
        parent, err := s.ProductFindByID(ctx, product.GetParentID())
        if err != nil {
            return err
        }
        
        if parent.HasVariantDimensions() {
            dimNames, _ := parent.GetVariantDimensionNames()
            variantMetas, _ := product.GetMetas()
            
            // Ensure variant has all required dimension metas
            for _, dim := range dimNames {
                if _, ok := variantMetas[dim]; !ok {
                    return fmt.Errorf("variant missing required dimension meta: %s", dim)
                }
            }
        }
    }
    
    // ... continue with create
}
```

**Usage Flow:**

```go
// 1. Create parent with dimension schema (simple array)
parent := NewProduct().
    SetTitle("Nike Air Max").
    SetStatus(PRODUCT_STATUS_PARENT)
parent.SetVariantDimensions([]string{"color", "size"})
store.ProductCreate(ctx, parent)

// 2. Or with structured dimensions (validation rules)
parent.SetVariantDimensions([]shopstore.VariantDimension{
    {Name: "color", Required: true, Options: []string{"red", "blue", "black"}},
    {Name: "size", Required: true, Options: []string{"8", "9", "10", "11"}},
})

// 3. Create variant - validated against schema
variant := NewProduct().
    SetParentID(parent.GetID()).
    SetMeta("color", "red").  // Required
    SetMeta("size", "9")      // Required
    SetMeta("material", "leather")  // Optional, not in schema
store.ProductCreate(ctx, variant) // Succeeds

// 4. Invalid variant - missing required dimension
badVariant := NewProduct().
    SetParentID(parent.GetID()).
    SetMeta("color", "blue")  // Missing "size"!
store.ProductCreate(ctx, badVariant) // Fails: variant missing required dimension meta: size
```

**Matrix Generation Helper:**

```go
// Auto-generate all variant combinations from dimensions
func (p *Product) GenerateVariants(store StoreInterface, ctx context.Context) ([]ProductInterface, error) {
    if !p.HasVariantDimensions() {
        return nil, errors.New("parent has no dimensions defined")
    }
    
    dims, err := p.GetVariantDimensions()
    if err != nil {
        return nil, err
    }
    
    // Generate cartesian product of all option combinations
    combinations := cartesianProduct(dims)
    
    variants := make([]ProductInterface, len(combinations))
    for i, combo := range combinations {
        variant := NewProduct().
            SetParentID(p.GetID()).
            SetTitle(p.GetTitle())
        
        // Set metas for each dimension
        for dimName, value := range combo {
            variant.SetMeta(dimName, value)
        }
        
        // Generate SKU from combination
        sku := generateVariantSKU(p.GetID(), combo)
        variant.SetSKU(sku)
        
        variants[i] = variant
    }
    
    return variants, nil
}
```

## Helper Methods

```go
// On Product entity

// GetVariants loads child products from store
func (p *Product) GetVariants(store StoreInterface, ctx context.Context) ([]ProductInterface, error) {
    if p.GetParentID() != "" {
        return nil, errors.New("cannot get variants of a variant")
    }
    return store.ProductList(ctx, NewProductQuery().SetParentID(p.GetID()))
}

// GetParent loads parent product from store
func (p *Product) GetParent(store StoreInterface, ctx context.Context) (ProductInterface, error) {
    if p.GetParentID() == "" {
        return nil, errors.New("product has no parent")
    }
    return store.ProductFindByID(ctx, p.GetParentID())
}

// GetVariantOptions returns unique option values across variants
// e.g., {"size": ["8", "9", "10"], "color": ["red", "blue"]}
func (p *Product) GetVariantOptions(store StoreInterface, ctx context.Context) (map[string][]string, error)
```

## Use Cases

### 1. Simple Color Variants

```go
parent := NewProduct().
    SetTitle("T-Shirt").
    SetDescription("Premium cotton tee").
    SetStatus(PRODUCT_STATUS_PARENT)
store.ProductCreate(ctx, parent)

for _, color := range []string{"red", "blue", "green"} {
    variant := NewProduct().
        SetTitle("T-Shirt").
        SetParentID(parent.GetID()).
        SetSKU("TSH-" + strings.ToUpper(color)).
        SetPriceFloat(29.99).
        SetMeta("color", color).
        SetStatus(PRODUCT_STATUS_ACTIVE)
    store.ProductCreate(ctx, variant)
}
```

### 2. Size + Color Matrix

```go
// Parent: "Nike Air Max"
// Variants: "Nike Air Max - Size 8 - Red", "Nike Air Max - Size 8 - Blue", etc.

for _, size := range []string{"8", "9", "10", "11"} {
    for _, color := range []string{"red", "black", "white"} {
        variant := NewProduct().
            SetTitle("Nike Air Max").
            SetParentID(parent.GetID()).
            SetSKU(fmt.Sprintf("NAM-%s-%s", size, color)).
            SetPriceFloat(129.99 + sizePremium(size)). // Size 11 costs more
            SetQuantityInt(rand.Intn(50)).
            SetMeta("size", size).
            SetMeta("color", color)
        store.ProductCreate(ctx, variant)
    }
}
```

### 3. Variant-Specific Media

```go
// Red variant has red shoe images
// Blue variant has blue shoe images
redVariant.AddMedia(ctx, store, redShoeImage)
blueVariant.AddMedia(ctx, store, blueShoeImage)
```

### 4. Variant-Specific Pricing

```go
// Red is on sale
goldVariant.SetPriceFloat(199.99).SetMeta("sale_price", "149.99")
```

## UI Considerations

Product display page should:
1. Load parent product (title, description, common media)
2. Load variants via `ProductVariantList()`
3. Show option selectors (dropdowns for size/color)
4. On selection: update displayed price, stock, images
5. Add to cart: use variant product ID, not parent

## Testing Requirements

1. Create parent product
2. Create variant products with parent ID
3. Query variants by parent ID
4. Ensure variant has access to parent data
5. Test variant-specific fields (different price, stock)
6. Prevent circular references (variant as parent)
7. Test deletion: deleting parent should handle variants

## Acceptance Criteria

- [ ] `parent_id` column added to product table
- [ ] `GetParentID()` / `SetParentID()` methods on Product
- [ ] `IsVariant()` helper method
- [ ] `SetParentID()` added to ProductQuery
- [ ] `ProductVariantList()` store method
- [ ] Parent/variant relationship validation
- [ ] Comprehensive tests
- [ ] Documentation for variant workflow

## Migration Path

Existing products remain unchanged (ParentID = "").
To convert existing product to parent-variants:

1. Create parent product (copy of original)
2. Mark original as variant with parent ID
3. Create additional variants as needed
4. Update original product to be variant

Or use migration helper:
```go
func (s *Store) ConvertToParentProduct(ctx context.Context, productID string) (ProductInterface, error)
```

## Estimated Effort

**Medium** - 4-6 hours:
- Database migration: 30 min
- Entity methods: 1 hour
- Query builder update: 30 min
- Store methods: 1 hour
- Tests: 2 hours

## Future Extensions

- **Variant templates**: Pre-defined variant sets (S/M/L always created)
- **Variant generation**: Auto-create from option matrix
- **Variant inheritance**: Inherit price from parent unless overridden
- **Variant stock rollup**: Parent shows "low stock" if any variant low
- **Variant bundles**: "Family pack" containing one of each variant
