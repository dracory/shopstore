package shopstore

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/dromara/carbon/v2"
)

// CategoryInterface defines the contract for category entities in the shop store.
// Categories support hierarchical structures (parent-child relationships),
// soft deletion, metadata storage, and status management.
type CategoryInterface interface {
	// DataObject methods

	// Data returns a map of all field values for serialization.
	Data() map[string]string
	// DataChanged returns a map of only the fields that have been modified since load.
	DataChanged() map[string]string
	// MarkAsNotDirty resets the dirty state, clearing all change tracking.
	MarkAsNotDirty()

	// Setters and Getters

	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) CategoryInterface

	// GetDescription returns the category description.
	GetDescription() string
	// SetDescription sets the category description.
	SetDescription(description string) CategoryInterface

	// GetID returns the unique identifier.
	GetID() string
	// SetID sets the unique identifier.
	SetID(id string) CategoryInterface

	// GetMemo returns the internal memo.
	GetMemo() string
	// SetMemo sets the internal memo.
	SetMemo(memo string) CategoryInterface

	// GetMetas returns all metadata as a map.
	GetMetas() (map[string]string, error)
	// GetMeta returns a specific metadata value by name.
	GetMeta(name string) string
	// SetMeta sets a single metadata value.
	SetMeta(name string, value string) error
	// SetMetas replaces all metadata with the provided map.
	SetMetas(metas map[string]string) error
	// MetasUpsert merges the provided metadata with existing values.
	MetasUpsert(metas map[string]string) error
	// MetaRemove removes a single metadata entry.
	MetaRemove(name string) error
	// MetasRemove removes multiple metadata entries.
	MetasRemove(names []string) error

	// GetParentID returns the parent category ID (empty if root).
	GetParentID() string
	// SetParentID sets the parent category ID.
	SetParentID(parentID string) CategoryInterface

	// GetStatus returns the current status.
	GetStatus() string
	// SetStatus sets the current status.
	SetStatus(status string) CategoryInterface

	// GetTitle returns the category title.
	GetTitle() string
	// SetTitle sets the category title.
	SetTitle(title string) CategoryInterface

	// GetSoftDeletedAt returns the soft deletion timestamp.
	GetSoftDeletedAt() string
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(deletedAt string) CategoryInterface

	// GetUpdatedAt returns the last update timestamp.
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) CategoryInterface

	// Status predicates

	// IsActive returns true if status is active.
	IsActive() bool
	// IsDraft returns true if status is draft.
	IsDraft() bool
	// IsInactive returns true if status is inactive.
	IsInactive() bool
	// IsSoftDeleted returns true if the category is soft deleted.
	IsSoftDeleted() bool

	// Hierarchy predicates

	// IsRoot returns true if this is a root category (no parent).
	IsRoot() bool
	// IsChild returns true if this category has a parent.
	IsChild() bool
}

// DiscountInterface defines the contract for discount/promotion entities.
// Discounts support temporal validity (start/end dates), amount-based discounts,
// soft deletion, metadata storage, and status management.
type DiscountInterface interface {
	// DataObject methods

	// Data returns a map of all field values for serialization.
	Data() map[string]string
	// DataChanged returns a map of only the fields that have been modified since load.
	DataChanged() map[string]string
	// MarkAsNotDirty resets the dirty state, clearing all change tracking.
	MarkAsNotDirty()

	// Setters and Getters

	// GetAmount returns the discount amount.
	GetAmount() float64
	// SetAmount sets the discount amount.
	SetAmount(amount float64) DiscountInterface

	// GetCode returns the unique discount code.
	GetCode() string
	// SetCode sets the unique discount code.
	SetCode(code string) DiscountInterface

	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) DiscountInterface

	// GetDescription returns the discount description.
	GetDescription() string
	// SetDescription sets the discount description.
	SetDescription(description string) DiscountInterface

	// GetEndsAt returns the end date/time as a string.
	GetEndsAt() string
	// GetEndsAtCarbon returns the end date/time as a Carbon instance.
	GetEndsAtCarbon() *carbon.Carbon
	// SetEndsAt sets the end date/time.
	SetEndsAt(endsAt string) DiscountInterface

	// GetID returns the unique identifier.
	GetID() string
	// SetID sets the unique identifier.
	SetID(id string) DiscountInterface

	// GetMemo returns the internal memo.
	GetMemo() string
	// SetMemo sets the internal memo.
	SetMemo(memo string) DiscountInterface

	// GetMeta returns a specific metadata value by name.
	GetMeta(name string) string
	// GetMetas returns all metadata as a map.
	GetMetas() (map[string]string, error)
	// SetMeta sets a single metadata value.
	SetMeta(name string, value string) error
	// SetMetas replaces all metadata with the provided map.
	SetMetas(metas map[string]string) error
	// MetasUpsert merges the provided metadata with existing values.
	MetasUpsert(metas map[string]string) error
	// MetaRemove removes a single metadata entry.
	MetaRemove(name string) error
	// MetasRemove removes multiple metadata entries.
	MetasRemove(names []string) error

	// GetSoftDeletedAt returns the soft deletion timestamp.
	GetSoftDeletedAt() string
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(deletedAt string) DiscountInterface

	// GetStartsAt returns the start date/time as a string.
	GetStartsAt() string
	// GetStartsAtCarbon returns the start date/time as a Carbon instance.
	GetStartsAtCarbon() *carbon.Carbon
	// SetStartsAt sets the start date/time.
	SetStartsAt(startsAt string) DiscountInterface

	// GetStatus returns the current status.
	GetStatus() string
	// SetStatus sets the current status.
	SetStatus(status string) DiscountInterface

	// GetTitle returns the discount title.
	GetTitle() string
	// SetTitle sets the discount title.
	SetTitle(title string) DiscountInterface

	// GetType returns the discount type.
	GetType() string
	// SetType sets the discount type.
	SetType(type_ string) DiscountInterface

	// GetUpdatedAt returns the last update timestamp.
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) DiscountInterface

	// Status predicates

	// IsActive returns true if status is active.
	IsActive() bool
	// IsDraft returns true if status is draft.
	IsDraft() bool
	// IsInactive returns true if status is inactive.
	IsInactive() bool

	// Temporal predicates

	// IsStarted returns true if the discount period has started.
	IsStarted() bool
	// IsEnded returns true if the discount period has ended.
	IsEnded() bool
	// IsExpired returns true if the discount is no longer valid.
	IsExpired() bool
	// IsValidNow returns true if the discount is currently valid (started and not ended).
	IsValidNow() bool
}

// MediaInterface defines the contract for media entities (images, videos, etc).
// Media files are associated with entities via EntityID, support sequencing for ordering,
// soft deletion, metadata storage, and type classification.
type MediaInterface interface {
	// DataObject methods

	// Data returns a map of all field values for serialization.
	Data() map[string]string
	// DataChanged returns a map of only the fields that have been modified since load.
	DataChanged() map[string]string
	// MarkAsNotDirty resets the dirty state, clearing all change tracking.
	MarkAsNotDirty()

	// Setters and Getters

	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) MediaInterface

	// GetDescription returns the media description.
	GetDescription() string
	// SetDescription sets the media description.
	SetDescription(description string) MediaInterface

	// GetEntityID returns the associated entity ID.
	GetEntityID() string
	// SetEntityID sets the associated entity ID.
	SetEntityID(entityID string) MediaInterface

	// GetID returns the unique identifier.
	GetID() string
	// SetID sets the unique identifier.
	SetID(id string) MediaInterface

	// GetMemo returns the internal memo.
	GetMemo() string
	// SetMemo sets the internal memo.
	SetMemo(memo string) MediaInterface

	// GetMetas returns all metadata as a map.
	GetMetas() (map[string]string, error)
	// GetMeta returns a specific metadata value by name.
	GetMeta(name string) string
	// SetMeta sets a single metadata value.
	SetMeta(name string, value string) error
	// SetMetas replaces all metadata with the provided map.
	SetMetas(metas map[string]string) error
	// MetasUpsert merges the provided metadata with existing values.
	MetasUpsert(metas map[string]string) error
	// MetaRemove removes a single metadata entry.
	MetaRemove(name string) error
	// MetasRemove removes multiple metadata entries.
	MetasRemove(names []string) error

	// GetSequence returns the display sequence/order.
	GetSequence() int
	// SetSequence sets the display sequence/order.
	SetSequence(sequence int) MediaInterface

	// GetSoftDeletedAt returns the soft deletion timestamp.
	GetSoftDeletedAt() string
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(softDeletedAt string) MediaInterface

	// GetStatus returns the current status.
	GetStatus() string
	// SetStatus sets the current status.
	SetStatus(status string) MediaInterface

	// GetTitle returns the media title.
	GetTitle() string
	// SetTitle sets the media title.
	SetTitle(title string) MediaInterface

	// GetType returns the media type (image, video, etc).
	GetType() string
	// SetType sets the media type.
	SetType(mediaType string) MediaInterface

	// GetURL returns the media file URL.
	GetURL() string
	// SetURL sets the media file URL.
	SetURL(url string) MediaInterface

	// GetUpdatedAt returns the last update timestamp.
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) MediaInterface

	// Status predicates

	// IsActive returns true if status is active.
	IsActive() bool
	// IsDraft returns true if status is draft.
	IsDraft() bool
	// IsInactive returns true if status is inactive.
	IsInactive() bool
	// IsSoftDeleted returns true if the media is soft deleted.
	IsSoftDeleted() bool

	// Type predicates

	// IsImage returns true if media type is image.
	IsImage() bool
	// IsVideo returns true if media type is video.
	IsVideo() bool
}

// OrderInterface defines the contract for order entities.
// Orders track customer purchases with status workflow, pricing, quantity management,
// soft deletion, and metadata storage. Supports various order states from pending to completed.
type OrderInterface interface {
	// DataObject methods

	// Data returns a map of all field values for serialization.
	Data() map[string]string
	// DataChanged returns a map of only the fields that have been modified since load.
	DataChanged() map[string]string
	// MarkAsNotDirty resets the dirty state, clearing all change tracking.
	MarkAsNotDirty()

	// Order status predicates

	// IsAwaitingFulfillment returns true if order is awaiting fulfillment.
	IsAwaitingFulfillment() bool
	// IsAwaitingPayment returns true if order is awaiting payment.
	IsAwaitingPayment() bool
	// IsAwaitingPickup returns true if order is awaiting pickup.
	IsAwaitingPickup() bool
	// IsAwaitingShipment returns true if order is awaiting shipment.
	IsAwaitingShipment() bool
	// IsCancelled returns true if order is cancelled.
	IsCancelled() bool
	// IsCompleted returns true if order is completed.
	IsCompleted() bool
	// IsDeclined returns true if order is declined.
	IsDeclined() bool
	// IsDisputed returns true if order is disputed.
	IsDisputed() bool
	// IsManualVerificationRequired returns true if order requires manual verification.
	IsManualVerificationRequired() bool
	// IsPending returns true if order is pending.
	IsPending() bool
	// IsRefunded returns true if order is refunded.
	IsRefunded() bool
	// IsShipped returns true if order is shipped.
	IsShipped() bool

	// Setters and Getters

	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) OrderInterface

	// GetCustomerID returns the customer ID.
	GetCustomerID() string
	// SetCustomerID sets the customer ID.
	SetCustomerID(customerID string) OrderInterface

	// GetID returns the unique identifier.
	GetID() string
	// SetID sets the unique identifier.
	SetID(id string) OrderInterface

	// GetMemo returns the internal memo.
	GetMemo() string
	// SetMemo sets the internal memo.
	SetMemo(memo string) OrderInterface

	// GetMeta returns a specific metadata value by name.
	GetMeta(name string) string
	// SetMeta sets a single metadata value.
	SetMeta(name string, value string) error
	// GetMetas returns all metadata as a map.
	GetMetas() (map[string]string, error)
	// SetMetas replaces all metadata with the provided map.
	SetMetas(metas map[string]string) error
	// MetasUpsert merges the provided metadata with existing values.
	MetasUpsert(metas map[string]string) error
	// MetaRemove removes a single metadata entry.
	MetaRemove(name string) error
	// MetasRemove removes multiple metadata entries.
	MetasRemove(names []string) error

	// GetPrice returns the price as a string.
	GetPrice() string
	// SetPrice sets the price from a string.
	SetPrice(price string) OrderInterface
	// GetPriceFloat returns the price as a float64.
	GetPriceFloat() float64
	// SetPriceFloat sets the price from a float64.
	SetPriceFloat(price float64) OrderInterface

	// GetQuantity returns the quantity as a string.
	GetQuantity() string
	// SetQuantity sets the quantity from a string.
	SetQuantity(quantity string) OrderInterface
	// GetQuantityInt returns the quantity as an int64.
	GetQuantityInt() int64
	// SetQuantityInt sets the quantity from an int64.
	SetQuantityInt(quantity int64) OrderInterface

	// GetSoftDeletedAt returns the soft deletion timestamp.
	GetSoftDeletedAt() string
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(deletedAt string) OrderInterface

	// GetStatus returns the current status.
	GetStatus() string
	// SetStatus sets the current status.
	SetStatus(status string) OrderInterface

	// GetUpdatedAt returns the last update timestamp.
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) OrderInterface
}

// OrderLineItemInterface defines the contract for order line item entities.
// Line items represent individual products within an order with their own
// pricing, quantity, and status tracking. Supports soft deletion and metadata storage.
type OrderLineItemInterface interface {
	// DataObject methods

	// Data returns a map of all field values for serialization.
	Data() map[string]string
	// DataChanged returns a map of only the fields that have been modified since load.
	DataChanged() map[string]string
	// MarkAsNotDirty resets the dirty state, clearing all change tracking.
	MarkAsNotDirty()

	// Setters and Getters

	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) OrderLineItemInterface

	// GetID returns the unique identifier.
	GetID() string
	// SetID sets the unique identifier.
	SetID(id string) OrderLineItemInterface

	// GetMemo returns the internal memo.
	GetMemo() string
	// SetMemo sets the internal memo.
	SetMemo(memo string) OrderLineItemInterface

	// GetMetas returns all metadata as a map.
	GetMetas() (map[string]string, error)
	// SetMetas replaces all metadata with the provided map.
	SetMetas(metas map[string]string) error
	// GetMeta returns a specific metadata value by name.
	GetMeta(name string) string
	// SetMeta sets a single metadata value.
	SetMeta(name string, value string) error
	// MetasUpsert merges the provided metadata with existing values.
	MetasUpsert(metas map[string]string) error
	// MetaRemove removes a single metadata entry.
	MetaRemove(name string) error
	// MetasRemove removes multiple metadata entries.
	MetasRemove(names []string) error

	// GetOrderID returns the associated order ID.
	GetOrderID() string
	// SetOrderID sets the associated order ID.
	SetOrderID(orderID string) OrderLineItemInterface

	// GetPrice returns the price as a string.
	GetPrice() string
	// SetPrice sets the price from a string.
	SetPrice(price string) OrderLineItemInterface
	// GetPriceFloat returns the price as a float64.
	GetPriceFloat() float64
	// SetPriceFloat sets the price from a float64.
	SetPriceFloat(price float64) OrderLineItemInterface

	// GetProductID returns the associated product ID.
	GetProductID() string
	// SetProductID sets the associated product ID.
	SetProductID(productID string) OrderLineItemInterface

	// GetQuantity returns the quantity as a string.
	GetQuantity() string
	// SetQuantity sets the quantity from a string.
	SetQuantity(quantity string) OrderLineItemInterface
	// GetQuantityInt returns the quantity as an int64.
	GetQuantityInt() int64
	// SetQuantityInt sets the quantity from an int64.
	SetQuantityInt(quantity int64) OrderLineItemInterface

	// GetSoftDeletedAt returns the soft deletion timestamp.
	GetSoftDeletedAt() string
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(deletedAt string) OrderLineItemInterface

	// GetStatus returns the current status.
	GetStatus() string
	// SetStatus sets the current status.
	SetStatus(status string) OrderLineItemInterface

	// GetTitle returns the line item title.
	GetTitle() string
	// SetTitle sets the line item title.
	SetTitle(title string) OrderLineItemInterface

	// GetUpdatedAt returns the last update timestamp.
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) OrderLineItemInterface

	// Status predicates

	// IsActive returns true if status is active.
	IsActive() bool
	// IsCancelled returns true if status is cancelled.
	IsCancelled() bool
	// IsCompleted returns true if status is completed.
	IsCompleted() bool
	// IsDraft returns true if status is draft.
	IsDraft() bool

	// Business logic predicates

	// HasQuantity returns true if quantity is greater than zero.
	HasQuantity() bool
	// IsFree returns true if price is zero.
	IsFree() bool
}

// VariantMatrixSchema defines a single dimension for product variants.
// Used to define variant attributes like color, size, material, etc.
// Supports optional predefined options and required/optional flags.
type VariantMatrixSchema struct {
	Name     string   `json:"name"`              // "color", "size"
	Required bool     `json:"required"`          // must variant have this?
	Options  []string `json:"options,omitempty"` // allowed values (optional)
}

// ProductInterface defines the contract for product entities.
// Products support parent-child relationships (variants), pricing, stock management,
// variant matrix configuration, soft deletion, metadata storage, and status management.
type ProductInterface interface {
	// DataObject methods

	// Data returns a map of all field values for serialization.
	Data() map[string]string
	// DataChanged returns a map of only the fields that have been modified since load.
	DataChanged() map[string]string
	// MarkAsNotDirty resets the dirty state, clearing all change tracking.
	MarkAsNotDirty()

	// Product status predicates

	// IsActive returns true if status is active.
	IsActive() bool
	// IsDisabled returns true if status is disabled.
	IsDisabled() bool
	// IsDraft returns true if status is draft.
	IsDraft() bool
	// IsParent returns true if this product has variants (has variant matrix schema).
	IsParent() bool
	// IsSoftDeleted returns true if the product is soft deleted.
	IsSoftDeleted() bool
	// IsFree returns true if price is zero.
	IsFree() bool
	// IsVariant returns true if this product is a variant of another product.
	IsVariant() bool
	// Slug returns the URL-friendly slug for the product.
	Slug() string

	// Stock predicates

	// HasStock returns true if quantity is greater than zero.
	HasStock() bool
	// IsOutOfStock returns true if quantity is zero.
	IsOutOfStock() bool

	// Price predicates

	// IsPaid returns true if price is greater than zero.
	IsPaid() bool

	// Setters and Getters

	// GetCreatedAt returns the creation timestamp as a string.
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance.
	GetCreatedAtCarbon() *carbon.Carbon
	// SetCreatedAt sets the creation timestamp.
	SetCreatedAt(createdAt string) ProductInterface

	// GetDescription returns the full product description.
	GetDescription() string
	// SetDescription sets the full product description.
	SetDescription(description string) ProductInterface

	// GetID returns the unique identifier.
	GetID() string
	// SetID sets the unique identifier.
	SetID(id string) ProductInterface

	// GetMemo returns the internal memo.
	GetMemo() string
	// SetMemo sets the internal memo.
	SetMemo(memo string) ProductInterface

	// GetMeta returns a specific metadata value by name.
	GetMeta(name string) string
	// SetMeta sets a single metadata value.
	SetMeta(name string, value string) error

	// GetMetas returns all metadata as a map.
	GetMetas() (map[string]string, error)
	// SetMetas replaces all metadata with the provided map.
	SetMetas(metas map[string]string) error
	// MetasUpsert merges the provided metadata with existing values.
	MetasUpsert(metas map[string]string) error
	// MetaRemove removes a single metadata entry.
	MetaRemove(name string) error
	// MetasRemove removes multiple metadata entries.
	MetasRemove(names []string) error

	// GetParentID returns the parent product ID (empty if not a variant).
	GetParentID() string
	// SetParentID sets the parent product ID.
	SetParentID(parentID string) ProductInterface

	// GetPrice returns the price as a string.
	GetPrice() string
	// SetPrice sets the price from a string.
	SetPrice(price string) ProductInterface
	// GetPriceFloat returns the price as a float64.
	GetPriceFloat() float64
	// SetPriceFloat sets the price from a float64.
	SetPriceFloat(price float64) ProductInterface

	// GetQuantity returns the stock quantity as a string.
	GetQuantity() string
	// SetQuantity sets the stock quantity from a string.
	SetQuantity(quantity string) ProductInterface
	// GetQuantityInt returns the stock quantity as an int64.
	GetQuantityInt() int64
	// SetQuantityInt sets the stock quantity from an int64.
	SetQuantityInt(quantity int64) ProductInterface

	// GetSoftDeletedAt returns the soft deletion timestamp.
	GetSoftDeletedAt() string
	// GetSoftDeletedAtCarbon returns the soft deletion timestamp as a Carbon instance.
	GetSoftDeletedAtCarbon() *carbon.Carbon
	// SetSoftDeletedAt sets the soft deletion timestamp.
	SetSoftDeletedAt(deletedAt string) ProductInterface

	// GetShortDescription returns the short/abbreviated description.
	GetShortDescription() string
	// SetShortDescription sets the short/abbreviated description.
	SetShortDescription(shortDescription string) ProductInterface

	// GetStatus returns the current status.
	GetStatus() string
	// SetStatus sets the current status.
	SetStatus(status string) ProductInterface

	// GetTitle returns the product title.
	GetTitle() string
	// SetTitle sets the product title.
	SetTitle(title string) ProductInterface

	// GetUpdatedAt returns the last update timestamp.
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance.
	GetUpdatedAtCarbon() *carbon.Carbon
	// SetUpdatedAt sets the last update timestamp.
	SetUpdatedAt(updatedAt string) ProductInterface

	// GetVariantMatrixSchema returns the variant matrix schema configuration.
	GetVariantMatrixSchema() (VariantMatrixSchema, error)
	// SetVariantMatrixSchema sets the variant matrix schema configuration.
	SetVariantMatrixSchema(schema VariantMatrixSchema) error
	// HasVariantMatrixSchema returns true if a variant matrix schema is defined.
	HasVariantMatrixSchema() bool

	// GetVariantMatrixValues returns the variant attribute values for this product.
	GetVariantMatrixValues() (map[string]string, error)
	// SetVariantMatrixValues sets the variant attribute values for this product.
	SetVariantMatrixValues(values map[string]string) error
}

// StoreInterface defines the contract for the shop store database operations.
// Provides CRUD operations, soft deletion, counting, listing with pagination,
// and variant management for all entity types (categories, discounts, media, orders, products).
type StoreInterface interface {
	// AutoMigrate creates or updates database tables to match the current schema.
	AutoMigrate() error
	// DB returns the underlying SQL database connection.
	DB() *sql.DB
	// EnableDebug enables or disables debug logging for SQL queries.
	EnableDebug(debug bool, sqlLogger ...*slog.Logger)

	// Table name methods

	// CategoryTableName returns the database table name for categories.
	CategoryTableName() string
	// DiscountTableName returns the database table name for discounts.
	DiscountTableName() string
	// MediaTableName returns the database table name for media.
	MediaTableName() string
	// OrderTableName returns the database table name for orders.
	OrderTableName() string
	// OrderLineItemTableName returns the database table name for order line items.
	OrderLineItemTableName() string
	// ProductTableName returns the database table name for products.
	ProductTableName() string

	// Category operations

	// CategoryCount returns the total count of categories matching the query options.
	CategoryCount(ctx context.Context, options CategoryQueryInterface) (int64, error)
	// CategoryCreate inserts a new category into the database.
	CategoryCreate(context context.Context, category CategoryInterface) error
	// CategoryDelete permanently deletes a category from the database.
	CategoryDelete(context context.Context, category CategoryInterface) error
	// CategoryDeleteByID permanently deletes a category by its ID.
	CategoryDeleteByID(context context.Context, categoryID string) error
	// CategoryFindByID retrieves a category by its unique ID.
	CategoryFindByID(context context.Context, categoryID string) (CategoryInterface, error)
	// CategoryList retrieves a list of categories matching the query options.
	CategoryList(context context.Context, options CategoryQueryInterface) ([]CategoryInterface, error)
	// CategorySoftDelete soft deletes a category by setting the deleted timestamp.
	CategorySoftDelete(context context.Context, category CategoryInterface) error
	// CategorySoftDeleteByID soft deletes a category by its ID.
	CategorySoftDeleteByID(context context.Context, categoryID string) error
	// CategoryUpdate updates an existing category in the database.
	CategoryUpdate(contxt context.Context, category CategoryInterface) error

	// Discount operations

	// DiscountCount returns the total count of discounts matching the query options.
	DiscountCount(ctx context.Context, options DiscountQueryInterface) (int64, error)
	// DiscountCreate inserts a new discount into the database.
	DiscountCreate(ctx context.Context, discount DiscountInterface) error
	// DiscountDelete permanently deletes a discount from the database.
	DiscountDelete(ctx context.Context, discount DiscountInterface) error
	// DiscountDeleteByID permanently deletes a discount by its ID.
	DiscountDeleteByID(ctx context.Context, discountID string) error
	// DiscountFindByID retrieves a discount by its unique ID.
	DiscountFindByID(ctx context.Context, discountID string) (DiscountInterface, error)
	// DiscountFindByCode retrieves a discount by its unique code.
	DiscountFindByCode(ctx context.Context, code string) (DiscountInterface, error)
	// DiscountList retrieves a list of discounts matching the query options.
	DiscountList(ctx context.Context, options DiscountQueryInterface) ([]DiscountInterface, error)
	// DiscountSoftDelete soft deletes a discount by setting the deleted timestamp.
	DiscountSoftDelete(ctx context.Context, discount DiscountInterface) error
	// DiscountSoftDeleteByID soft deletes a discount by its ID.
	DiscountSoftDeleteByID(ctx context.Context, discountID string) error
	// DiscountUpdate updates an existing discount in the database.
	DiscountUpdate(ctx context.Context, discount DiscountInterface) error

	// Media operations

	// MediaCount returns the total count of media matching the query options.
	MediaCount(ctx context.Context, options MediaQueryInterface) (int64, error)
	// MediaCreate inserts a new media into the database.
	MediaCreate(ctx context.Context, media MediaInterface) error
	// MediaDelete permanently deletes a media from the database.
	MediaDelete(ctx context.Context, media MediaInterface) error
	// MediaDeleteByID permanently deletes a media by its ID.
	MediaDeleteByID(ctx context.Context, mediaID string) error
	// MediaFindByID retrieves a media by its unique ID.
	MediaFindByID(ctx context.Context, mediaID string) (MediaInterface, error)
	// MediaList retrieves a list of media matching the query options.
	MediaList(ctx context.Context, options MediaQueryInterface) ([]MediaInterface, error)
	// MediaSoftDelete soft deletes a media by setting the deleted timestamp.
	MediaSoftDelete(ctx context.Context, media MediaInterface) error
	// MediaSoftDeleteByID soft deletes a media by its ID.
	MediaSoftDeleteByID(ctx context.Context, mediaID string) error
	// MediaUpdate updates an existing media in the database.
	MediaUpdate(ctx context.Context, media MediaInterface) error

	// Order operations

	// OrderCount returns the total count of orders matching the query options.
	OrderCount(ctx context.Context, options OrderQueryInterface) (int64, error)
	// OrderCreate inserts a new order into the database.
	OrderCreate(ctx context.Context, order OrderInterface) error
	// OrderDelete permanently deletes an order from the database.
	OrderDelete(ctx context.Context, order OrderInterface) error
	// OrderDeleteByID permanently deletes an order by its ID.
	OrderDeleteByID(ctx context.Context, id string) error
	// OrderFindByID retrieves an order by its unique ID.
	OrderFindByID(ctx context.Context, id string) (OrderInterface, error)
	// OrderList retrieves a list of orders matching the query options.
	OrderList(ctx context.Context, options OrderQueryInterface) ([]OrderInterface, error)
	// OrderSoftDelete soft deletes an order by setting the deleted timestamp.
	OrderSoftDelete(ctx context.Context, order OrderInterface) error
	// OrderSoftDeleteByID soft deletes an order by its ID.
	OrderSoftDeleteByID(ctx context.Context, id string) error
	// OrderUpdate updates an existing order in the database.
	OrderUpdate(ctx context.Context, order OrderInterface) error

	// OrderLineItem operations

	// OrderLineItemCount returns the total count of line items matching the query options.
	OrderLineItemCount(ctx context.Context, options OrderLineItemQueryInterface) (int64, error)
	// OrderLineItemCreate inserts a new line item into the database.
	OrderLineItemCreate(ctx context.Context, orderLineItem OrderLineItemInterface) error
	// OrderLineItemDelete permanently deletes a line item from the database.
	OrderLineItemDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error
	// OrderLineItemDeleteByID permanently deletes a line item by its ID.
	OrderLineItemDeleteByID(ctx context.Context, id string) error
	// OrderLineItemFindByID retrieves a line item by its unique ID.
	OrderLineItemFindByID(ctx context.Context, id string) (OrderLineItemInterface, error)
	// OrderLineItemList retrieves a list of line items matching the query options.
	OrderLineItemList(ctx context.Context, options OrderLineItemQueryInterface) ([]OrderLineItemInterface, error)
	// OrderLineItemSoftDelete soft deletes a line item by setting the deleted timestamp.
	OrderLineItemSoftDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error
	// OrderLineItemSoftDeleteByID soft deletes a line item by its ID.
	OrderLineItemSoftDeleteByID(ctx context.Context, id string) error
	// OrderLineItemUpdate updates an existing line item in the database.
	OrderLineItemUpdate(ctx context.Context, orderLineItem OrderLineItemInterface) error

	// Product operations

	// ProductCount returns the total count of products matching the query options.
	ProductCount(ctx context.Context, options ProductQueryInterface) (int64, error)
	// ProductCreate inserts a new product into the database.
	ProductCreate(ctx context.Context, product ProductInterface) error
	// ProductDelete permanently deletes a product from the database.
	ProductDelete(ctx context.Context, product ProductInterface) error
	// ProductDeleteByID permanently deletes a product by its ID.
	ProductDeleteByID(ctx context.Context, productID string) error
	// ProductFindByID retrieves a product by its unique ID.
	ProductFindByID(ctx context.Context, productID string) (ProductInterface, error)
	// ProductList retrieves a list of products matching the query options.
	ProductList(ctx context.Context, options ProductQueryInterface) ([]ProductInterface, error)
	// ProductSoftDelete soft deletes a product by setting the deleted timestamp.
	ProductSoftDelete(ctx context.Context, product ProductInterface) error
	// ProductSoftDeleteByID soft deletes a product by its ID.
	ProductSoftDeleteByID(ctx context.Context, productID string) error
	// ProductUpdate updates an existing product in the database.
	ProductUpdate(ctx context.Context, product ProductInterface) error

	// Variant operations

	// ProductVariantList retrieves all variants for a parent product.
	ProductVariantList(ctx context.Context, parentID string) ([]ProductInterface, error)
	// ProductIsParent checks if a product is a parent (has variants).
	ProductIsParent(ctx context.Context, productID string) (bool, error)
	// ProductGetParent retrieves the parent product for a variant.
	ProductGetParent(ctx context.Context, productID string) (ProductInterface, error)
}
