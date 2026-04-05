package shopstore

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/dromara/carbon/v2"
)

type CategoryInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Setters and Getters

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) CategoryInterface

	GetDescription() string
	SetDescription(description string) CategoryInterface

	GetID() string
	SetID(id string) CategoryInterface

	GetMemo() string
	SetMemo(memo string) CategoryInterface

	GetMetas() (map[string]string, error)
	GetMeta(name string) string
	SetMeta(name string, value string) error
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	GetParentID() string
	SetParentID(parentID string) CategoryInterface

	GetStatus() string
	SetStatus(status string) CategoryInterface

	GetTitle() string
	SetTitle(title string) CategoryInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) CategoryInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) CategoryInterface

	// Status predicates
	IsActive() bool
	IsDraft() bool
	IsInactive() bool
	IsSoftDeleted() bool

	// Hierarchy predicates
	IsRoot() bool
	IsChild() bool
}

type DiscountInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Setters and Getters

	GetAmount() float64
	SetAmount(amount float64) DiscountInterface

	GetCode() string
	SetCode(code string) DiscountInterface

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) DiscountInterface

	GetDescription() string
	SetDescription(description string) DiscountInterface

	GetEndsAt() string
	GetEndsAtCarbon() *carbon.Carbon
	SetEndsAt(endsAt string) DiscountInterface

	GetID() string
	SetID(id string) DiscountInterface

	GetMemo() string
	SetMemo(memo string) DiscountInterface

	GetMeta(name string) string
	GetMetas() (map[string]string, error)
	SetMeta(name string, value string) error
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) DiscountInterface

	GetStartsAt() string
	GetStartsAtCarbon() *carbon.Carbon
	SetStartsAt(startsAt string) DiscountInterface

	GetStatus() string
	SetStatus(status string) DiscountInterface

	GetTitle() string
	SetTitle(title string) DiscountInterface

	GetType() string
	SetType(type_ string) DiscountInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) DiscountInterface

	// Status predicates
	IsActive() bool
	IsDraft() bool
	IsInactive() bool

	// Temporal predicates
	IsStarted() bool
	IsEnded() bool
	IsExpired() bool
	IsValidNow() bool
}

type MediaInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Setters and Getters

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) MediaInterface

	GetDescription() string
	SetDescription(description string) MediaInterface

	GetEntityID() string
	SetEntityID(entityID string) MediaInterface

	GetID() string
	SetID(id string) MediaInterface

	GetMemo() string
	SetMemo(memo string) MediaInterface

	GetMetas() (map[string]string, error)
	GetMeta(name string) string
	SetMeta(name string, value string) error
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	GetSequence() int
	SetSequence(sequence int) MediaInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(softDeletedAt string) MediaInterface

	GetStatus() string
	SetStatus(status string) MediaInterface

	GetTitle() string
	SetTitle(title string) MediaInterface

	GetType() string
	SetType(mediaType string) MediaInterface

	GetURL() string
	SetURL(url string) MediaInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) MediaInterface

	// Status predicates
	IsActive() bool
	IsDraft() bool
	IsInactive() bool
	IsSoftDeleted() bool

	// Type predicates
	IsImage() bool
	IsVideo() bool
}

type OrderInterface interface {
	// Inherited from DataObject
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Methods
	IsAwaitingFulfillment() bool
	IsAwaitingPayment() bool
	IsAwaitingPickup() bool
	IsAwaitingShipment() bool
	IsCancelled() bool
	IsCompleted() bool
	IsDeclined() bool
	IsDisputed() bool
	IsManualVerificationRequired() bool
	IsPending() bool
	IsRefunded() bool
	IsShipped() bool

	// Setters and Getters
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) OrderInterface

	GetCustomerID() string
	SetCustomerID(customerID string) OrderInterface

	GetID() string
	SetID(id string) OrderInterface

	GetMemo() string
	SetMemo(memo string) OrderInterface

	GetMeta(name string) string
	SetMeta(name string, value string) error
	GetMetas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	GetPrice() string
	SetPrice(price string) OrderInterface
	GetPriceFloat() float64
	SetPriceFloat(price float64) OrderInterface

	GetQuantity() string
	SetQuantity(quantity string) OrderInterface
	GetQuantityInt() int64
	SetQuantityInt(quantity int64) OrderInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) OrderInterface

	GetStatus() string
	SetStatus(status string) OrderInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) OrderInterface
}

type OrderLineItemInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) OrderLineItemInterface

	GetID() string
	SetID(id string) OrderLineItemInterface

	GetMemo() string
	SetMemo(memo string) OrderLineItemInterface

	GetMetas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	GetMeta(name string) string
	SetMeta(name string, value string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	GetOrderID() string
	SetOrderID(orderID string) OrderLineItemInterface

	GetPrice() string
	SetPrice(price string) OrderLineItemInterface
	GetPriceFloat() float64
	SetPriceFloat(price float64) OrderLineItemInterface

	GetProductID() string
	SetProductID(productID string) OrderLineItemInterface

	GetQuantity() string
	SetQuantity(quantity string) OrderLineItemInterface
	GetQuantityInt() int64
	SetQuantityInt(quantity int64) OrderLineItemInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) OrderLineItemInterface

	GetStatus() string
	SetStatus(status string) OrderLineItemInterface

	GetTitle() string
	SetTitle(title string) OrderLineItemInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) OrderLineItemInterface

	// Status predicates
	IsActive() bool
	IsCancelled() bool
	IsCompleted() bool
	IsDraft() bool

	// Business logic predicates
	HasQuantity() bool
	IsFree() bool
}

type VariantDimension struct {
	Name     string   `json:"name"`              // "color", "size"
	Required bool     `json:"required"`          // must variant have this?
	Options  []string `json:"options,omitempty"` // allowed values (optional)
}

type ProductInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Methods
	IsActive() bool
	IsDisabled() bool
	IsDraft() bool
	IsParent() bool
	IsSoftDeleted() bool
	IsFree() bool
	IsVariant() bool
	Slug() string

	// Stock predicates
	HasStock() bool
	IsOutOfStock() bool

	// Price predicates
	IsPaid() bool

	// Setters and Getters

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) ProductInterface

	GetDescription() string
	SetDescription(description string) ProductInterface

	GetID() string
	SetID(id string) ProductInterface

	GetMemo() string
	SetMemo(memo string) ProductInterface

	GetMeta(name string) string
	SetMeta(name string, value string) error

	GetMetas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	GetParentID() string
	SetParentID(parentID string) ProductInterface

	GetPrice() string
	SetPrice(price string) ProductInterface
	GetPriceFloat() float64
	SetPriceFloat(price float64) ProductInterface

	GetQuantity() string
	SetQuantity(quantity string) ProductInterface
	GetQuantityInt() int64
	SetQuantityInt(quantity int64) ProductInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) ProductInterface

	GetShortDescription() string
	SetShortDescription(shortDescription string) ProductInterface

	GetStatus() string
	SetStatus(status string) ProductInterface

	GetTitle() string
	SetTitle(title string) ProductInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) ProductInterface

	GetVariantDimensions() ([]VariantDimension, error)
	SetVariantDimensions(dims interface{}) error
	GetVariantDimensionNames() ([]string, error)
	HasVariantDimensions() bool
	SetVariantMatrixSchema(schema string) ProductInterface
	SetVariantMatrixValues(values string) ProductInterface
}

type StoreInterface interface {
	AutoMigrate() error
	DB() *sql.DB
	EnableDebug(debug bool, sqlLogger ...*slog.Logger)

	CategoryTableName() string
	DiscountTableName() string
	MediaTableName() string
	OrderTableName() string
	OrderLineItemTableName() string
	ProductTableName() string

	CategoryCount(ctx context.Context, options CategoryQueryInterface) (int64, error)
	CategoryCreate(context context.Context, category CategoryInterface) error
	CategoryDelete(context context.Context, category CategoryInterface) error
	CategoryDeleteByID(context context.Context, categoryID string) error
	CategoryFindByID(context context.Context, categoryID string) (CategoryInterface, error)
	CategoryList(context context.Context, options CategoryQueryInterface) ([]CategoryInterface, error)
	CategorySoftDelete(context context.Context, category CategoryInterface) error
	CategorySoftDeleteByID(context context.Context, categoryID string) error
	CategoryUpdate(contxt context.Context, category CategoryInterface) error

	DiscountCount(ctx context.Context, options DiscountQueryInterface) (int64, error)
	DiscountCreate(ctx context.Context, discount DiscountInterface) error
	DiscountDelete(ctx context.Context, discount DiscountInterface) error
	DiscountDeleteByID(ctx context.Context, discountID string) error
	DiscountFindByID(ctx context.Context, discountID string) (DiscountInterface, error)
	DiscountFindByCode(ctx context.Context, code string) (DiscountInterface, error)
	DiscountList(ctx context.Context, options DiscountQueryInterface) ([]DiscountInterface, error)
	DiscountSoftDelete(ctx context.Context, discount DiscountInterface) error
	DiscountSoftDeleteByID(ctx context.Context, discountID string) error
	DiscountUpdate(ctx context.Context, discount DiscountInterface) error

	MediaCount(ctx context.Context, options MediaQueryInterface) (int64, error)
	MediaCreate(ctx context.Context, media MediaInterface) error
	MediaDelete(ctx context.Context, media MediaInterface) error
	MediaDeleteByID(ctx context.Context, mediaID string) error
	MediaFindByID(ctx context.Context, mediaID string) (MediaInterface, error)
	MediaList(ctx context.Context, options MediaQueryInterface) ([]MediaInterface, error)
	MediaSoftDelete(ctx context.Context, media MediaInterface) error
	MediaSoftDeleteByID(ctx context.Context, mediaID string) error
	MediaUpdate(ctx context.Context, media MediaInterface) error

	OrderCount(ctx context.Context, options OrderQueryInterface) (int64, error)
	OrderCreate(ctx context.Context, order OrderInterface) error
	OrderDelete(ctx context.Context, order OrderInterface) error
	OrderDeleteByID(ctx context.Context, id string) error
	OrderFindByID(ctx context.Context, id string) (OrderInterface, error)
	OrderList(ctx context.Context, options OrderQueryInterface) ([]OrderInterface, error)
	OrderSoftDelete(ctx context.Context, order OrderInterface) error
	OrderSoftDeleteByID(ctx context.Context, id string) error
	OrderUpdate(ctx context.Context, order OrderInterface) error

	OrderLineItemCount(ctx context.Context, options OrderLineItemQueryInterface) (int64, error)
	OrderLineItemCreate(ctx context.Context, orderLineItem OrderLineItemInterface) error
	OrderLineItemDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error
	OrderLineItemDeleteByID(ctx context.Context, id string) error
	OrderLineItemFindByID(ctx context.Context, id string) (OrderLineItemInterface, error)
	OrderLineItemList(ctx context.Context, options OrderLineItemQueryInterface) ([]OrderLineItemInterface, error)
	OrderLineItemSoftDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error
	OrderLineItemSoftDeleteByID(ctx context.Context, id string) error
	OrderLineItemUpdate(ctx context.Context, orderLineItem OrderLineItemInterface) error

	ProductCount(ctx context.Context, options ProductQueryInterface) (int64, error)
	ProductCreate(ctx context.Context, product ProductInterface) error
	ProductDelete(ctx context.Context, product ProductInterface) error
	ProductDeleteByID(ctx context.Context, productID string) error
	ProductFindByID(ctx context.Context, productID string) (ProductInterface, error)
	ProductList(ctx context.Context, options ProductQueryInterface) ([]ProductInterface, error)
	ProductSoftDelete(ctx context.Context, product ProductInterface) error
	ProductSoftDeleteByID(ctx context.Context, productID string) error
	ProductUpdate(ctx context.Context, product ProductInterface) error

	// Variant operations
	ProductVariantList(ctx context.Context, parentID string) ([]ProductInterface, error)
	ProductIsParent(ctx context.Context, productID string) (bool, error)
	ProductGetParent(ctx context.Context, productID string) (ProductInterface, error)
}
