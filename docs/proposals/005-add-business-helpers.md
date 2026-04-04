# Proposal: Add Business Logic Helpers

## Status: Draft

## Overview

Add common business logic helpers to entities that simplify calculations and provide semantic operations frequently needed in e-commerce applications.

## Proposed Helpers

### Order Entity

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

### OrderLineItem Entity

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

### Discount Entity

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

### Product Entity

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

### Category Entity

```go
// Check if this is a root-level category
func (category *Category) IsRoot() bool {
    return category.ParentID() == ""
}

// Check if this is a child category
func (category *Category) IsChild() bool {
    return category.ParentID() != ""
}

// Get breadcrumb path suggestion (stubs for future implementation)
func (category *Category) BreadcrumbPath() []string {
    // Returns array of category IDs from root to this category
    // Would require store access or pre-populated data
    return []string{category.ID()}
}
```

### Media Entity

```go
// Check if media is attached to an entity
func (media *Media) IsAttached() bool {
    return media.EntityID() != ""
}

// Check if media has a valid URL
func (media *Media) HasValidURL() bool {
    url := media.URL()
    return url != "" && (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"))
}

// Get file extension suggestion from URL
func (media *Media) SuggestedExtension() string {
    url := media.URL()
    parts := strings.Split(url, ".")
    if len(parts) > 1 {
        return parts[len(parts)-1]
    }
    return ""
}
```

## Interface Updates

Add helper methods to respective interfaces:

```go
type OrderInterface interface {
    // ... existing methods ...
    
    Total() float64
    HasCustomer() bool
    CanCancel() bool
    CanRefund() bool
    RequiresShipping() bool
}

type DiscountInterface interface {
    // ... existing methods ...
    
    CalculateDiscount(originalPrice float64) float64
    ApplyDiscount(originalPrice float64) float64
    MatchesCode(code string) bool
}
```

## Testing Requirements

Each helper needs tests covering:
1. Normal operation
2. Edge cases (zero values, empty strings)
3. Boundary conditions

Example test:
```go
func TestDiscountCalculateDiscount(t *testing.T) {
    // Percent discount
    d := NewDiscount().
        SetStatus(DISCOUNT_STATUS_ACTIVE).
        SetType(DISCOUNT_TYPE_PERCENT).
        SetAmount(10).
        SetStartsAt(carbon.Now().SubDay().ToDateTimeString())
    
    if got := d.CalculateDiscount(100); got != 10 {
        t.Fatalf("expected 10%% discount on 100 to be 10, got %f", got)
    }
    
    // Amount discount
    d.SetType(DISCOUNT_TYPE_AMOUNT).SetAmount(15)
    if got := d.CalculateDiscount(100); got != 15 {
        t.Fatalf("expected amount discount of 15, got %f", got)
    }
    
    // Discount larger than price
    d.SetAmount(150)
    if got := d.CalculateDiscount(100); got != 100 {
        t.Fatalf("expected discount capped at price, got %f", got)
    }
    
    // Invalid discount
    d.SetStatus(DISCOUNT_STATUS_DRAFT)
    if got := d.CalculateDiscount(100); got != 0 {
        t.Fatalf("expected invalid discount to return 0, got %f", got)
    }
}
```

## Benefits

- **Reduced Duplication:** Common calculations in one place
- **Business Logic Encapsulation:** Rules defined with the entity
- **Testability:** Business rules can be unit tested
- **Discoverability:** Helpers available via IDE autocomplete
- **Consistency:** Same calculations across the application

## Acceptance Criteria

- [ ] Order has `Total()`, `HasCustomer()`, `CanCancel()`, `CanRefund()`, `RequiresShipping()`
- [ ] OrderLineItem has `Subtotal()`, `HasProduct()`, `HasOrder()`
- [ ] Discount has `CalculateDiscount()`, `ApplyDiscount()`, `MatchesCode()`
- [ ] Product has `InventoryValue()`, `IsAvailable()`, `Summary()`
- [ ] Category has `IsRoot()`, `IsChild()`
- [ ] Media has `IsAttached()`, `HasValidURL()`, `SuggestedExtension()`
- [ ] All interfaces updated
- [ ] Comprehensive tests for all helpers

## Estimated Effort

**Medium** - 6-8 hours including tests and documentation.
