# Proposal: Product Business Logic Helpers

## Status: Draft

## Overview

Add business logic helpers to the Product entity for inventory value and availability checks.

## Proposed Helpers

```go
// Calculate total inventory value
func (product *Product) InventoryValue() float64 {
    return product.PriceFloat() * float64(product.QuantityInt())
}

// Check if product can be purchased
func (product *Product) IsAvailable() bool {
    return product.IsActive() && product.HasStock() && !product.IsSoftDeleted()
}

// Get product summary for display
func (product *Product) Summary() string {
    if product.ShortDescription() != "" {
        return product.ShortDescription()
    }
    desc := product.Description()
    if len(desc) > 100 {
        return desc[:100] + "..."
    }
    return desc
}
```

## Interface Updates

Add to `ProductInterface`:

```go
type ProductInterface interface {
    // ... existing methods ...
    
    InventoryValue() float64
    IsAvailable() bool
    Summary() string
}
```

## Implementation Notes

### Naming Consideration
Consider `IsPurchasable()` instead of `IsAvailable()` for clearer e-commerce intent.

### Edge Cases
- `InventoryValue()` with negative price/quantity
- `Summary()` handles empty descriptions gracefully

### Future Extension
Consider `IsLowStock()` helper:
```go
func (product *Product) IsLowStock() bool {
    return product.GetQuantityInt() > 0 && product.GetQuantityInt() < 5
}
```

## Testing Requirements

Cover:
- Inventory value calculations
- `IsAvailable()` with all status/stock/deletion combinations
- `Summary()` with short description, long description, empty description
- Edge cases (zero values, boundary conditions)

## Acceptance Criteria

- [ ] `InventoryValue()` implemented and tested
- [ ] `IsAvailable()` implemented and tested
- [ ] `Summary()` implemented and tested
- [ ] Interface updated
