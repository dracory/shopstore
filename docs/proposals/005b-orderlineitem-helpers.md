# Proposal: OrderLineItem Business Logic Helpers

## Status: Draft

## Overview

Add business logic helpers to the OrderLineItem entity for subtotal calculations and relationship checks.

## Proposed Helpers

```go
// Calculate subtotal for line item
func (item *OrderLineItem) Subtotal() float64 {
    return item.PriceFloat() * float64(item.QuantityInt())
}

// Check if linked to a product
func (item *OrderLineItem) HasProduct() bool {
    return item.ProductID() != ""
}

// Check if linked to an order
func (item *OrderLineItem) HasOrder() bool {
    return item.OrderID() != ""
}
```

## Interface Updates

Add to `OrderLineItemInterface`:

```go
type OrderLineItemInterface interface {
    // ... existing methods ...
    
    Subtotal() float64
    HasProduct() bool
    HasOrder() bool
}
```

## Implementation Notes

- `Subtotal()` uses `float64` - document precision limitations
- These are simple O(1) lookups with no I/O

## Testing Requirements

Cover:
- Subtotal calculation with various price/quantity combinations
- Zero and negative values
- Empty string checks for `HasProduct()` and `HasOrder()`

## Acceptance Criteria

- [ ] `Subtotal()` implemented and tested
- [ ] `HasProduct()` implemented and tested
- [ ] `HasOrder()` implemented and tested
- [ ] Interface updated
