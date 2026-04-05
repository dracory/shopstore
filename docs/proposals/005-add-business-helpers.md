# Proposal: Add Business Logic Helpers (Overview)

## Status: Draft

## Overview

Add common business logic helpers to entities that simplify calculations and provide semantic operations frequently needed in e-commerce applications.

## Split Proposals

This proposal has been split into individual entity-specific proposals for easier review:

| Proposal | Entity | Helpers | Status |
|----------|--------|---------|--------|
| [005a](005a-order-helpers.md) | Order | `Total()`, `HasCustomer()`, `CanCancel()`, `CanRefund()`, `RequiresShipping()` | Draft |
| [005b](005b-orderlineitem-helpers.md) | OrderLineItem | `Subtotal()`, `HasProduct()`, `HasOrder()` | Draft |
| [005c](005c-discount-helpers.md) | Discount | `CalculateDiscount()`, `ApplyDiscount()`, `MatchesCode()` | Draft |
| [005d](005d-product-helpers.md) | Product | `InventoryValue()`, `IsAvailable()`, `Summary()` | Draft |
| [005e](005e-media-helpers.md) | Media | `IsAttached()` + `IsValidURL()`, `GetFileExtension()` helpers | Draft |

## Cross-Cutting Concerns

### Floating-Point Precision
All monetary calculations use `float64`, which can introduce precision errors. These helpers are for **display and simple logic**, not financial calculations requiring exact precision.

### Interface Breaking Changes
Adding methods to interfaces is a **breaking change** for any external implementations. This affects:
- `OrderInterface`
- `OrderLineItemInterface`
- `DiscountInterface`
- `ProductInterface`
- `MediaInterface`

### Testing Pattern
Each helper should have comprehensive tests covering:
1. Normal operation
2. Edge cases (zero values, empty strings)
3. Boundary conditions

## Already Implemented

Category helpers (`IsRoot()`, `IsChild()`) are already implemented and do not need to be added.

## Benefits

- **Reduced Duplication:** Common calculations in one place
- **Business Logic Encapsulation:** Rules defined with the entity
- **Testability:** Business rules can be unit tested
- **Discoverability:** Helpers available via IDE autocomplete
- **Consistency:** Same calculations across the application

## Estimated Effort

**Medium** - 6-8 hours total across all entities including tests.
