# Proposal: Order Business Logic Helpers

## Status: Draft

## Overview

Add business logic helpers to the Order entity for common e-commerce calculations and state checks.

## Proposed Helpers

```go
// Calculate total order value
func (order *Order) Total() float64 {
    return order.PriceFloat() * float64(order.QuantityInt())
}

// Check if order has associated customer
func (order *Order) HasCustomer() bool {
    return order.CustomerID() != ""
}

// Check if order can be cancelled
func (order *Order) CanCancel() bool {
    return order.IsPending() || order.IsAwaitingPayment()
}

// Check if order can be refunded
func (order *Order) CanRefund() bool {
    return order.IsCompleted() || order.IsAwaitingFulfillment()
}

// Check if order requires shipping
func (order *Order) RequiresShipping() bool {
    return !order.IsCancelled() && !order.IsDeclined() && !order.IsRefunded()
}
```

## Interface Updates

Add to `OrderInterface`:

```go
type OrderInterface interface {
    // ... existing methods ...
    
    Total() float64
    HasCustomer() bool
    CanCancel() bool
    CanRefund() bool
    RequiresShipping() bool
}
```

## Implementation Notes

### Naming Consideration
Consider `OrderTotal()` instead of `Total()` for consistency with `OrderLineItem.Subtotal()` pattern.

### Edge Cases
- `CanCancel()` and `CanRefund()` currently don't include time-based constraints
- `Total()` uses `float64` - document precision limitations for financial calculations

### Future Extension
Consider `TotalWithItems()` helper when line items are available:
```go
func (order *Order) TotalWithItems(items []OrderLineItemInterface) float64
```

## Testing Requirements

Cover:
- Normal operation for each helper
- Edge cases (zero values, empty customer ID)
- All status combinations for `CanCancel()`, `CanRefund()`, `RequiresShipping()`
- Negative price/quantity handling

## Acceptance Criteria

- [ ] `Total()` implemented and tested
- [ ] `HasCustomer()` implemented and tested
- [ ] `CanCancel()` implemented and tested
- [ ] `CanRefund()` implemented and tested
- [ ] `RequiresShipping()` implemented and tested
- [ ] Interface updated
