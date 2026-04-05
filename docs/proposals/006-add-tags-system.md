# Proposal: Add Tags System

## Status: Draft

## Overview

Add a flexible tagging system that allows categorizing products, orders, and potentially other entities with user-defined labels. Tags enable better organization, filtering, and search capabilities.

## Entity Design

### Tag Entity

```go
type Tag struct {
    dataobject.DataObject
}

type TagInterface interface {
    Data() map[string]string
    DataChanged() map[string]string
    MarkAsNotDirty()

    // Core fields
    GetID() string
    SetID(id string) TagInterface

    GetName() string
    SetName(name string) TagInterface

    GetSlug() string
    SetSlug(slug string) TagInterface

    GetDescription() string
    SetDescription(description string) TagInterface

    GetStatus() string
    SetStatus(status string) TagInterface

    // Status predicates
    IsActive() bool
    IsDraft() bool
    IsInactive() bool
    IsSoftDeleted() bool

    // Metadata
    GetMeta(name string) string
    SetMeta(name string, value string) error
    GetMetas() (map[string]string, error)
    SetMetas(metas map[string]string) error

    // Timestamps
    GetCreatedAt() string
    GetCreatedAtCarbon() *carbon.Carbon
    SetCreatedAt(createdAt string) TagInterface

    GetUpdatedAt() string
    GetUpdatedAtCarbon() *carbon.Carbon
    SetUpdatedAt(updatedAt string) TagInterface

    GetSoftDeletedAt() string
    GetSoftDeletedAtCarbon() *carbon.Carbon
    SetSoftDeletedAt(deletedAt string) TagInterface
}
```

### Tag Relationship (Polymorphic)

```go
type TagRelation struct {
    dataobject.DataObject
}

type TagRelationInterface interface {
    GetID() string
    SetID(id string) TagRelationInterface

    GetTagID() string
    SetTagID(tagID string) TagRelationInterface

    GetEntityType() string  // "product", "order", "category", etc.
    SetEntityType(entityType string) TagRelationInterface

    GetEntityID() string
    SetEntityID(entityID string) TagRelationInterface

    GetCreatedAt() string
    SetCreatedAt(createdAt string) TagRelationInterface
}
```

## New Files Required

| File | Purpose |
|------|---------|
| `type_tag.go` | Tag entity implementation |
| `type_tag_relation.go` | Tag relationship implementation |
| `type_tag_test.go` | Tag entity tests |
| `type_tag_relation_test.go` | Tag relation tests |
| `query_tag.go` | Tag query builder |
| `query_tag_relation.go` | Tag relation query builder |
| `store_tag.go` | Tag store methods |
| `store_tag_relation.go` | Tag relation store methods |

## Constants to Add

```go
// Tag statuses
const TAG_STATUS_DRAFT = "draft"
const TAG_STATUS_ACTIVE = "active"
const TAG_STATUS_INACTIVE = "inactive"

// Entity types for tag relations
const TAG_ENTITY_TYPE_PRODUCT = "product"
const TAG_ENTITY_TYPE_ORDER = "order"
const TAG_ENTITY_TYPE_CATEGORY = "category"
const TAG_ENTITY_TYPE_DISCOUNT = "discount"
const TAG_ENTITY_TYPE_MEDIA = "media"

// Columns
const COLUMN_TAG_ID = "tag_id"
const COLUMN_ENTITY_TYPE = "entity_type"
```

## Store Interface Additions

```go
type StoreInterface interface {
    // ... existing methods ...

    // Tag operations
    TagCount(ctx context.Context, options TagQueryInterface) (int64, error)
    TagCreate(ctx context.Context, tag TagInterface) error
    TagDelete(ctx context.Context, tag TagInterface) error
    TagDeleteByID(ctx context.Context, tagID string) error
    TagFindByID(ctx context.Context, tagID string) (TagInterface, error)
    TagFindBySlug(ctx context.Context, slug string) (TagInterface, error)
    TagList(ctx context.Context, options TagQueryInterface) ([]TagInterface, error)
    TagSoftDelete(ctx context.Context, tag TagInterface) error
    TagSoftDeleteByID(ctx context.Context, tagID string) error
    TagUpdate(ctx context.Context, tag TagInterface) error

    // Tag relation operations
    TagRelationCount(ctx context.Context, options TagRelationQueryInterface) (int64, error)
    TagRelationCreate(ctx context.Context, relation TagRelationInterface) error
    TagRelationDelete(ctx context.Context, relation TagRelationInterface) error
    TagRelationDeleteByID(ctx context.Context, relationID string) error
    TagRelationList(ctx context.Context, options TagRelationQueryInterface) ([]TagRelationInterface, error)
    
    // Convenience methods
    EntityTagAdd(ctx context.Context, entityType, entityID string, tag TagInterface) error
    EntityTagRemove(ctx context.Context, entityType, entityID string, tagID string) error
    EntityTagsList(ctx context.Context, entityType, entityID string) ([]TagInterface, error)
    TagEntitiesList(ctx context.Context, tagID string, entityType string) ([]string, error) // returns entity IDs
}
```

## Helper Methods on Entities

Consider adding to `Product`, `Order`, etc.:

```go
// On ProductInterface
func (p *Product) AddTag(tag TagInterface) error
func (p *Product) RemoveTag(tagID string) error
func (p *Product) GetTags() ([]TagInterface, error)
func (p *Product) HasTag(tagID string) bool
```

## Query Builder Features

### TagQuery
- `SetNameLike()` - search by name pattern
- `SetStatus()` - filter by status
- `SetSlug()` - find by slug
- Standard pagination/sorting

### TagRelationQuery
- `SetTagID()` - filter by specific tag
- `SetEntityType()` - filter by entity type
- `SetEntityID()` - filter by specific entity
- `SetEntityIDIn()` - filter by multiple entities

## Use Cases

1. **Product Categorization**: Tag products with "new-arrival", "bestseller", "sale"
2. **Order Organization**: Tag orders with "vip-customer", "rush-order", "gift-wrapped"
3. **Dynamic Collections**: Create product collections based on tags
4. **Filtering**: Filter products by multiple tags
5. **Reporting**: Generate reports by tag (e.g., sales of "bestseller" items)

## Database Schema

```sql
-- Tags table
CREATE TABLE shop_tag (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    metas JSON,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    soft_deleted_at DATETIME NOT NULL DEFAULT '9999-12-31 23:59:59'
);

-- Tag relations table (polymorphic)
CREATE TABLE shop_tag_relation (
    id VARCHAR(255) PRIMARY KEY,
    tag_id VARCHAR(255) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES shop_tag(id),
    UNIQUE KEY unique_tag_entity (tag_id, entity_type, entity_id)
);

-- Indexes
CREATE INDEX idx_tag_status ON shop_tag(status);
CREATE INDEX idx_tag_slug ON shop_tag(slug);
CREATE INDEX idx_tag_relation_tag ON shop_tag_relation(tag_id);
CREATE INDEX idx_tag_relation_entity ON shop_tag_relation(entity_type, entity_id);
```

## Implementation Notes

### Slug Generation
Auto-generate slug from name using existing `str.Slugify()` helper. Handle duplicates by appending number.

### Validation
- Tag name required, max length 255
- Slug unique across all tags
- Entity type restricted to known types (product, order, etc.)

### Soft Delete Behavior
- Soft-deleted tags should not appear in entity tag lists
- Hard delete of tag should cascade delete relations

### Performance
- Tag relations table will grow quickly - ensure proper indexing
- Consider caching entity tag lists for frequently accessed entities

## Testing Requirements

1. Tag CRUD operations
2. Tag relation CRUD operations
3. Slug uniqueness handling
4. Entity tag add/remove/list
5. Query filtering by tag
6. Soft delete behavior
7. Duplicate relation prevention (unique constraint)

## Acceptance Criteria

- [ ] Tag entity implemented (`type_tag.go`)
- [ ] Tag relation entity implemented (`type_tag_relation.go`)
- [ ] Tag query builder (`query_tag.go`)
- [ ] Tag relation query builder (`query_tag_relation.go`)
- [ ] Store methods for tags and relations
- [ ] Convenience methods on entities (optional)
- [ ] Database migration support
- [ ] Comprehensive tests
- [ ] Constants defined
- [ ] Interface updates

## Estimated Effort

**Large** - 12-16 hours including tests and documentation. More complex than business helpers due to:
- New entity with full CRUD
- Polymorphic relationship pattern
- Additional query builders
- Database schema changes
- More test coverage needed

## Future Extensions

- **Tag Groups**: Group tags into categories (e.g., "Product Types", "Seasons")
- **Tag Colors**: Visual tags with color coding
- **Tag Analytics**: Track which tags drive sales
- **Auto-tagging**: Rules-based automatic tagging
- **Tag Synonyms**: Alternative names for same tag (search improvement)
