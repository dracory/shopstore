# Proposal: Entity Consistency Checks for Delete Operations

## Status: Draft

## Overview

Deleting or soft-deleting a parent entity currently leaves child/related entities in an "active" state. This can leave the database in an inconsistent state where child rows still reference a deleted parent. This proposal adds guard checks to `Store.Delete` and `Store.SoftDelete` methods so that the database is not left in an incomplete state.

## Affected Relationships

| Parent | Child/Related | Foreign Key | Current Behavior |
|--------|---------------|-------------|------------------|
| `Order` | `OrderLineItem` | `order_id` | `OrderLineItem` rows remain active when an order is soft-deleted. `OrderLineItemList` still returns them. |
| `Product` | `Product` (variants) | `parent_id` | Deleting a parent product leaves child variants still referencing a deleted parent. |
| `Category` | `Category` (children) | `parent_id` | Deleting a parent category leaves child categories still referencing a deleted parent. |
| `Product` | `OrderLineItem` | `product_id` | Deleting a product leaves `OrderLineItem` rows still referencing a deleted product. |
| `Product`/`Category`/`Order` | `Media` | `entity_id` | `Media` files remain associated with a deleted entity. |
| `Discount` | N/A | N/A | `Discount` has no child or related entities; no consistency check is needed. |

## Proposed Changes

`Store.Delete` and `Store.SoftDelete` methods should check for active children before completing the deletion. If active children or related entities are found, the method should return an error and **not** delete the parent.

A later enhancement can add a `Cascade` option or separate `DeleteWithCascade` / `SoftDeleteWithCascade` methods, but the immediate goal is to avoid leaving the database in an incomplete state.

### Example: `OrderDelete` (same logic for `OrderSoftDelete`)

```go
func (store *Store) OrderDelete(ctx context.Context, order OrderInterface) error {
    if order == nil {
        return errors.New("order is nil")
    }

    return store.OrderDeleteByID(ctx, order.GetID())
}

func (store *Store) OrderDeleteByID(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("order id is empty")
    }

    if err := store.assertOrderDeletable(ctx, id); err != nil {
        return err
    }

    _, err := store.db.Query().Table(store.orderTableName).Where(COLUMN_ID+" = ?", id).Delete()
    return err
}

// assertOrderDeletable returns an error if the order has active line items or media.
func (store *Store) assertOrderDeletable(ctx context.Context, orderID string) error {
    lineItemCount, err := store.OrderLineItemCount(ctx, NewOrderLineItemQuery().SetOrderID(orderID))
    if err != nil {
        return err
    }
    if lineItemCount > 0 {
        return ErrOrderHasActiveLineItems
    }

    mediaCount, err := store.MediaCount(ctx, NewMediaQuery().SetEntityID(orderID))
    if err != nil {
        return err
    }
    if mediaCount > 0 {
        return ErrOrderHasActiveMedia
    }

    return nil
}
```

### Example: `OrderSoftDelete` and `OrderSoftDeleteByID`

```go
func (store *Store) OrderSoftDelete(ctx context.Context, order OrderInterface) error {
    if order == nil {
        return errors.New("order is nil")
    }

    if err := store.assertOrderDeletable(ctx, order.GetID()); err != nil {
        return err
    }

    order.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
    return store.OrderUpdate(ctx, order)
}

func (store *Store) OrderSoftDeleteByID(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("order id is empty")
    }

    order, err := store.OrderFindByID(ctx, id)
    if err != nil {
        return err
    }
    if order == nil {
        return nil
    }

    return store.OrderSoftDelete(ctx, order)
}
```

### Example: `ProductDelete` (same logic for `ProductSoftDelete`)

```go
func (store *Store) ProductDelete(ctx context.Context, product ProductInterface) error {
    if product == nil {
        return errors.New("product is nil")
    }

    return store.ProductDeleteByID(ctx, product.GetID())
}

func (store *Store) ProductDeleteByID(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("product id is empty")
    }

    if err := store.assertProductDeletable(ctx, id); err != nil {
        return err
    }

    _, err := store.db.Query().Table(store.productTableName).Where(COLUMN_ID+" = ?", id).Delete()
    return err
}

// assertProductDeletable returns an error if the product has active variants, line items, or media.
func (store *Store) assertProductDeletable(ctx context.Context, productID string) error {
    variantCount, err := store.ProductCount(ctx, NewProductQuery().SetParentID(productID))
    if err != nil {
        return err
    }
    if variantCount > 0 {
        return ErrProductHasActiveVariants
    }

    lineItemCount, err := store.OrderLineItemCount(ctx, NewOrderLineItemQuery().SetProductID(productID))
    if err != nil {
        return err
    }
    if lineItemCount > 0 {
        return ErrProductHasActiveLineItems
    }

    mediaCount, err := store.MediaCount(ctx, NewMediaQuery().SetEntityID(productID))
    if err != nil {
        return err
    }
    if mediaCount > 0 {
        return ErrProductHasActiveMedia
    }

    return nil
}
```

### Example: `ProductSoftDelete`

```go
func (store *Store) ProductSoftDelete(ctx context.Context, product ProductInterface) error {
    if product == nil {
        return errors.New("product is nil")
    }

    if err := store.assertProductDeletable(ctx, product.GetID()); err != nil {
        return err
    }

    product.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
    return store.ProductUpdate(ctx, product)
}
```

### Example: `CategoryDelete` (same logic for `CategorySoftDelete`)

```go
func (store *Store) CategoryDelete(ctx context.Context, category CategoryInterface) error {
    if category == nil {
        return errors.New("category is nil")
    }

    return store.CategoryDeleteByID(ctx, category.GetID())
}

func (store *Store) CategoryDeleteByID(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("id is empty")
    }

    if err := store.assertCategoryDeletable(ctx, id); err != nil {
        return err
    }

    _, err := store.db.Query().Table(store.categoryTableName).Where(COLUMN_ID+" = ?", id).Delete()
    return err
}

// assertCategoryDeletable returns an error if the category has active children or media.
func (store *Store) assertCategoryDeletable(ctx context.Context, categoryID string) error {
    childCount, err := store.CategoryCount(ctx, NewCategoryQuery().SetParentID(categoryID))
    if err != nil {
        return err
    }
    if childCount > 0 {
        return ErrCategoryHasActiveChildren
    }

    mediaCount, err := store.MediaCount(ctx, NewMediaQuery().SetEntityID(categoryID))
    if err != nil {
        return err
    }
    if mediaCount > 0 {
        return ErrCategoryHasActiveMedia
    }

    return nil
}
```

### Example: `CategorySoftDelete`

```go
func (store *Store) CategorySoftDelete(ctx context.Context, category CategoryInterface) error {
    if category == nil {
        return errors.New("category is nil")
    }

    if err := store.assertCategoryDeletable(ctx, category.GetID()); err != nil {
        return err
    }

    category.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
    return store.CategoryUpdate(ctx, category)
}
```

### Example: `MediaDelete` and `MediaSoftDelete`

`Media` is a leaf entity and has no child rows. Its `entity_id` points to an external parent, so the parent entity is responsible for blocking deletion while `Media` still references it. `MediaDelete` and `MediaSoftDelete` therefore only need to handle the "not found" case consistently.

```go
func (store *Store) MediaSoftDeleteByID(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("id is empty")
    }

    media, err := store.MediaFindByID(ctx, id)
    if err != nil {
        return err
    }
    if media == nil {
        return nil
    }

    return store.MediaSoftDelete(ctx, media)
}
```

`MediaDeleteByID` should follow the same pattern: return `nil` when `MediaFindByID` returns `nil`. `MediaDelete` and `MediaSoftDelete` methods that receive an object should return an error when the object is `nil`.

## Cross-Entity Reference Checks

In addition to the direct parent-child checks, each parent `Delete` and `SoftDelete` method must check cross-entity references to avoid broken foreign key-like relationships.

| Parent | Referenced By | Check |
|--------|---------------|-------|
| `Product` | `OrderLineItem` | `OrderLineItemCount` with `SetProductID` |
| `Product` | `Media` | `MediaCount` with `SetEntityID` |
| `Category` | `Media` | `MediaCount` with `SetEntityID` |
| `Order` | `Media` | `MediaCount` with `SetEntityID` |

## Interface Updates

No interface changes are required. The `Store` methods keep the same signatures and return an error when an active relationship exists.

## Implementation Notes

- All checks are performed at the application layer because the ORM does not enforce foreign-key constraints.
- The `*Count` methods count according to the query's `SoftDeletedIncluded` value. The checks should not call `SetSoftDeletedIncluded(true)`; the default `false` (or not setting the option) counts only active rows.
- `Media` uses `entity_id` without an `entity_type` column, so `MediaCount` by `SetEntityID` is the only available check and assumes the referenced entity ID is unique across entity types. The globally-unique `GenerateShortID` strategy makes this safe today, but this should be treated as a hard constraint until an `entity_type` column is added.
- Errors should be deterministic and easy to test. The `assert*Deletable` helpers return the sentinel errors defined in the Sentinel Errors section so callers can distinguish blocking conditions programmatically.
- `Delete`, `DeleteByID`, `SoftDelete`, and `SoftDeleteByID` methods should all delegate to the same `assert*Deletable` helper before performing the delete or update.
- `*SoftDeleteByID` methods (e.g., `ProductSoftDeleteByID`, `CategorySoftDeleteByID`) should follow the `OrderSoftDeleteByID` pattern: guard against empty `id`, call `*FindByID`, return `nil` when not found, otherwise call `*SoftDelete`. The empty-ID guard is a deliberate improvement over the existing `OrderSoftDeleteByID` and `ProductSoftDeleteByID` implementations.
- `DeleteByID` methods should also be idempotent; if the row does not exist, the `assert` helper will find no active children and the delete will affect zero rows, returning `nil`. One edge case: orphaned child rows that reference a non-existent parent (e.g., an `OrderLineItem` with an `order_id` that no longer exists) would block the delete; this should not happen in practice if the consistency checks are enforced on all deletes, but any existing orphan data should be cleaned up before enabling this guard.
- `SoftDelete` methods must not mutate the entity before the consistency check passes. Call `assert*Deletable` before calling `SetSoftDeletedAt`.

### Sentinel Errors

```go
var (
    ErrOrderHasActiveLineItems = errors.New("cannot delete order with active line items")
    ErrOrderHasActiveMedia     = errors.New("cannot delete order with active media")

    ErrProductHasActiveVariants  = errors.New("cannot delete product with active variants")
    ErrProductHasActiveLineItems = errors.New("cannot delete product referenced by active order line items")
    ErrProductHasActiveMedia     = errors.New("cannot delete product with active media")

    ErrCategoryHasActiveChildren = errors.New("cannot delete category with active children")
    ErrCategoryHasActiveMedia    = errors.New("cannot delete category with active media")
)
```

## Future Extensions

- `DeleteWithCascade` and `SoftDeleteWithCascade` methods that explicitly delete children.
- `Force` option in a new `DeleteOptions` struct.
- Add an `entity_type` column to `Media` to remove the ID uniqueness assumption.
- Database-level `ON DELETE` behaviors if the underlying storage supports it.

## Testing Requirements

For each affected `Delete` and `SoftDelete` method:

- Normal delete when no active children or related entities exist.
- Error returned when active children or related entities exist.
- No partial state change when the error is returned (parent remains not deleted).
- For `SoftDelete`, the entity's `SoftDeletedAt` value is unchanged when the check fails.
- `DeleteByID` and `SoftDeleteByID` behave the same as `Delete` and `SoftDelete`.
- `DeleteByID` and `SoftDeleteByID` return `nil` when the entity is not found.

## Acceptance Criteria

- [ ] `OrderDelete` and `OrderSoftDelete` return an error when active `OrderLineItem` rows or `Media` rows exist.
- [ ] `ProductDelete` and `ProductSoftDelete` return an error when active variants, `OrderLineItem` rows, or `Media` rows exist.
- [ ] `CategoryDelete` and `CategorySoftDelete` return an error when active child categories or `Media` rows exist.
- [ ] `MediaDelete` and `MediaSoftDelete` remain unchanged (no child check required; return `nil` when the entity is not found).
- [ ] `DiscountDelete` and `DiscountSoftDelete` remain unchanged (no child or related entities).
- [ ] Existing tests continue to pass.
- [ ] New tests cover the error paths for each entity and each delete type.
