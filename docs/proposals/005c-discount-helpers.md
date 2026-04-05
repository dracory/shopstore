# Proposal: Discount Business Logic Helpers

## Status: Draft

## Overview

Add business logic helpers to the Discount entity for calculation and code matching.

## Proposed Helpers

```go
// Calculate discount amount for a given price
func (d *Discount) CalculateDiscount(originalPrice float64) float64 {
    if !d.IsValidNow() {
        return 0
    }
    
    if d.Type() == DISCOUNT_TYPE_PERCENT {
        return originalPrice * (d.Amount() / 100)
    }
    
    // Amount discount
    if d.Amount() > originalPrice {
        return originalPrice
    }
    return d.Amount()
}

// Calculate final price after discount
func (d *Discount) ApplyDiscount(originalPrice float64) float64 {
    return originalPrice - d.CalculateDiscount(originalPrice)
}

// Check if discount code matches
func (d *Discount) MatchesCode(code string) bool {
    return strings.EqualFold(d.Code(), code)
}
```

## Interface Updates

Add to `DiscountInterface`:

```go
type DiscountInterface interface {
    // ... existing methods ...
    
    CalculateDiscount(originalPrice float64) float64
    ApplyDiscount(originalPrice float64) float64
    MatchesCode(code string) bool
}
```

## Implementation Notes

### Edge Cases
- Negative `originalPrice` should return 0
- Discount amount capped at price for amount-type discounts
- Invalid discounts (not active, not started, expired) return 0

### Considerations
- `MatchesCode()` uses case-insensitive comparison
- Consider normalized code comparison (strip whitespace)

## Testing Requirements

Cover:
- Percent discount calculations
- Amount discount calculations
- Discount larger than price (capped at price)
- Invalid discount states (draft, expired, not started)
- Case-insensitive code matching
- Unicode in discount codes

## Acceptance Criteria

- [ ] `CalculateDiscount()` implemented and tested
- [ ] `ApplyDiscount()` implemented and tested
- [ ] `MatchesCode()` implemented and tested
- [ ] Interface updated
