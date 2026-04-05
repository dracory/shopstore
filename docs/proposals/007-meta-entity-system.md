# Proposal: Meta Entity System

## Status: Draft

## Overview

Introduce a dedicated Meta entity for query-critical metadata while preserving JSON metas for arbitrary display data. This hybrid approach enables efficient filtering/searching by meta values without sacrificing flexibility.

## Problem with Current Approach

Current JSON-based metas have limitations:
- Cannot efficiently query/filter by meta in SQL
- No database indexing on meta keys/values
- All reads require JSON parsing
- Large JSON blobs slow down queries

## Proposed Solution

Two-tier meta system:

| Storage | Use Case | API |
|---------|----------|-----|
| **Meta Table** | Queryable/filterable data (color, size, brand) | `SetSearchableMeta()`, `FindByMeta()` |
| **JSON Column** | Arbitrary display data (custom CSS, notes) | `SetMeta()`, `GetMeta()` |

## Entity Design

### Meta Entity

```go
type Meta struct {
    dataobject.DataObject
}

type MetaInterface interface {
    GetID() string
    SetID(id string) MetaInterface

    GetEntityType() string  // 'product', 'order', 'category'
    SetEntityType(entityType string) MetaInterface

    GetEntityID() string
    SetEntityID(entityID string) MetaInterface

    GetKey() string
    SetKey(key string) MetaInterface

    GetValue() string
    SetValue(value string) MetaInterface

    GetValueInt() int64
    SetValueInt(value int64) MetaInterface

    GetValueFloat() float64
    SetValueFloat(value float64) MetaInterface

    GetValueBool() bool
    SetValueBool(value bool) MetaInterface

    GetCreatedAt() string
    SetCreatedAt(createdAt string) MetaInterface

    GetUpdatedAt() string
    SetUpdatedAt(updatedAt string) MetaInterface
}
```

## New Files Required

| File | Purpose |
|------|---------|
| `type_meta.go` | Meta entity implementation |
| `type_meta_test.go` | Meta entity tests |
| `query_meta.go` | Meta query builder |
| `store_meta.go` | Meta store methods |

## Database Schema

```sql
CREATE TABLE shop_meta (
    id VARCHAR(255) PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    meta_key VARCHAR(255) NOT NULL,
    meta_value TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    
    UNIQUE KEY unique_entity_meta (entity_type, entity_id, meta_key),
    INDEX idx_entity_lookup (entity_type, entity_id),
    INDEX idx_meta_key_value (meta_key, meta_value(100))
);
```

## Store Interface Additions

```go
type StoreInterface interface {
    // ... existing methods ...

    // Meta CRUD
    MetaCount(ctx context.Context, options MetaQueryInterface) (int64, error)
    MetaCreate(ctx context.Context, meta MetaInterface) error
    MetaDelete(ctx context.Context, meta MetaInterface) error
    MetaDeleteByID(ctx context.Context, metaID string) error
    MetaFindByID(ctx context.Context, metaID string) (MetaInterface, error)
    MetaList(ctx context.Context, options MetaQueryInterface) ([]MetaInterface, error)
    MetaUpdate(ctx context.Context, meta MetaInterface) error

    // Convenience methods
    EntityMetaSet(ctx context.Context, entityType, entityID, key, value string) error
    EntityMetaGet(ctx context.Context, entityType, entityID, key string) (string, error)
    EntityMetaDelete(ctx context.Context, entityType, entityID, key string) error
    EntityMetasList(ctx context.Context, entityType, entityID string) (map[string]string, error)
    
    // Query by meta (the key feature)
    ProductsByMeta(ctx context.Context, key, value string) ([]ProductInterface, error)
    ProductsByMetaQuery(ctx context.Context, query MetaQueryInterface) ([]ProductInterface, error)
    OrdersByMeta(ctx context.Context, key, value string) ([]OrderInterface, error)
}
```

## Entity Helper Methods

Add to existing entities:

```go
// Product example
func (p *Product) SetSearchableMeta(key, value string) error
func (p *Product) GetSearchableMeta(key string) (string, error)
func (p *Product) DeleteSearchableMeta(key string) error
func (p *Product) GetAllSearchableMetas() (map[string]string, error)

// Type conversion helpers
func (p *Product) SetSearchableMetaInt(key string, value int64) error
func (p *Product) GetSearchableMetaInt(key string) (int64, error)
func (p *Product) SetSearchableMetaFloat(key string, value float64) error
func (p *Product) GetSearchableMetaFloat(key string) (float64, error)
func (p *Product) SetSearchableMetaBool(key string, value bool) error
func (p *Product) GetSearchableMetaBool(key string) (bool, error)
```

## Query Builder Features

```go
type MetaQueryInterface interface {
    SetEntityType(entityType string) MetaQueryInterface
    SetEntityID(entityID string) MetaQueryInterface
    SetKey(key string) MetaQueryInterface
    SetValue(value string) MetaQueryInterface
    SetValueLike(pattern string) MetaQueryInterface  // partial match
    SetValueGte(value string) MetaQueryInterface    // range queries
    SetValueLte(value string) MetaQueryInterface
    
    // Standard pagination
    SetLimit(limit int) MetaQueryInterface
    SetOffset(offset int) MetaQueryInterface
    SetOrderBy(orderBy string) MetaQueryInterface
    SetSortDirection(direction string) MetaQueryInterface
    
    // Count support
    SetCountOnly(countOnly bool) MetaQueryInterface
}
```

## Usage Examples

### Setting Searchable Metas

```go
// On product
product.SetSearchableMeta("color", "red")
product.SetSearchableMetaInt("size_cm", 42)
product.SetSearchableMetaBool("in_stock", true)
product.SetSearchableMetaFloat("weight_kg", 1.5)

// Or via store
store.EntityMetaSet(ctx, "product", product.GetID(), "brand", "Nike")
```

### Querying by Meta

```go
// Find all red products
redProducts, err := store.ProductsByMeta(ctx, "color", "red")

// Complex query
query := shopstore.NewMetaQuery().
    SetKey("size_cm").
    SetValueGte("40").
    SetValueLte("44")
products, err := store.ProductsByMetaQuery(ctx, query)

// Cross-entity meta search
allRedItems, err := store.MetaList(ctx, shopstore.NewMetaQuery().
    SetKey("color").
    SetValue("red"))
// Returns metas for products, orders, categories that are "red"
```

### When to Use Which

| Scenario | Use | Example |
|----------|-----|---------|
| Filter products by color | **Meta Table** | `SetSearchableMeta("color", "blue")` |
| Product display attributes | **JSON Metas** | `SetMeta("hover_animation", "fade")` |
| Price range filtering | **Meta Table** | `SetSearchableMetaFloat("sale_price", 19.99)` |
| Internal notes | **JSON Metas** | `SetMeta("staff_notes", "Check supplier")` |
| Brand filtering | **Meta Table** | `SetSearchableMeta("brand", "Nike")` |
| Custom CSS | **JSON Metas** | `SetMeta("custom_css", ".big{font-size:2em}")` |

## Constants to Add

```go
// Meta entity types
const META_ENTITY_TYPE_PRODUCT = "product"
const META_ENTITY_TYPE_ORDER = "order"
const META_ENTITY_TYPE_CATEGORY = "category"
const META_ENTITY_TYPE_DISCOUNT = "discount"
const META_ENTITY_TYPE_MEDIA = "media"

// Columns
const COLUMN_META_KEY = "meta_key"
const COLUMN_META_VALUE = "meta_value"
```

## Implementation Notes

### Type Conversion
Meta values stored as strings but typed getters/setters handle conversion:
- Int: `strconv.FormatInt()` / `strconv.ParseInt()`
- Float: `strconv.FormatFloat()` / `strconv.ParseFloat()`
- Bool: `strconv.FormatBool()` / `strconv.ParseBool()`

### Data Consistency
- Updating a searchable meta updates `shop_meta` table row
- No automatic sync between JSON and searchable metas - explicit API choice
- Deleting entity should cascade delete its searchable metas

### Indexing Strategy
- `UNIQUE` constraint prevents duplicate keys per entity
- Composite index on `(entity_type, entity_id)` for entity lookups
- Composite index on `(meta_key, meta_value)` for filtering

### Backwards Compatibility
- Existing JSON metas continue to work unchanged
- `GetMeta()`/`SetMeta()` behavior unchanged
- New methods `GetSearchableMeta()`/`SetSearchableMeta()` for table-based

## Testing Requirements

1. Meta entity defaults and getters/setters
2. Type conversion (int/float/bool/string)
3. Unique constraint enforcement (same entity can't have duplicate keys)
4. Query filtering by key/value
5. Range queries (Gte/Lte)
6. Entity integration (product.SetSearchableMeta)
7. Store convenience methods
8. Cascade delete behavior
9. Cross-entity meta queries

## Acceptance Criteria

- [ ] Meta entity implemented (`type_meta.go`)
- [ ] Meta query builder (`query_meta.go`)
- [ ] Store methods for meta CRUD
- [ ] Entity helper methods (`SetSearchableMeta`, etc.)
- [ ] Products/Orders by meta query methods
- [ ] Database migration support
- [ ] Constants defined
- [ ] Comprehensive tests
- [ ] Interface updates

## Estimated Effort

**Large** - 10-14 hours including tests. Complexity factors:
- New entity with type conversion
- Query builder additions
- Store method additions
- Entity integration across 5+ types
- Index optimization considerations

## Migration Path

For existing projects wanting to migrate JSON metas to searchable:

```go
// Migration helper
func (s *Store) MigrateJSONMetaToSearchable(ctx context.Context, entityType, key string) error {
    // Find all entities with this key in JSON metas
    // Create corresponding rows in shop_meta
    // Optional: remove from JSON after migration
}
```

## Future Extensions

- **Multi-value metas**: Store array of values (comma-separated or separate rows)
- **Meta groups**: Group related metas (dimensions, specifications)
- **Required metas**: Enforce certain metas must exist per entity type
- **Meta validation**: Schema/constraints on meta keys/values
- **Meta search indexing**: Full-text search on meta values
