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

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) CategoryInterface

	Description() string
	SetDescription(description string) CategoryInterface

	ID() string
	SetID(id string) CategoryInterface

	Memo() string
	SetMemo(memo string) CategoryInterface

	Metas() (map[string]string, error)
	Meta(name string) string
	SetMeta(name string, value string) error
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	ParentID() string
	SetParentID(parentID string) CategoryInterface

	Status() string
	SetStatus(status string) CategoryInterface

	Title() string
	SetTitle(title string) CategoryInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) CategoryInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) CategoryInterface
}

type DiscountInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Setters and Getters

	Amount() float64
	SetAmount(amount float64) DiscountInterface

	Code() string
	SetCode(code string) DiscountInterface

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) DiscountInterface

	Description() string
	SetDescription(description string) DiscountInterface

	EndsAt() string
	EndsAtCarbon() *carbon.Carbon
	SetEndsAt(endsAt string) DiscountInterface

	ID() string
	SetID(id string) DiscountInterface

	Memo() string
	SetMemo(memo string) DiscountInterface

	Meta(name string) string
	MetaRemove(name string) error
	SetMeta(name string, value string) error

	Metas() (map[string]string, error)
	MetasRemove(names []string) error
	MetasUpsert(metas map[string]string) error
	SetMetas(metas map[string]string) error

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) DiscountInterface

	StartsAt() string
	StartsAtCarbon() *carbon.Carbon
	SetStartsAt(startsAt string) DiscountInterface

	Status() string
	SetStatus(status string) DiscountInterface

	Title() string
	SetTitle(title string) DiscountInterface

	Type() string
	SetType(type_ string) DiscountInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) DiscountInterface
}

type MediaInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Setters and Getters

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) MediaInterface

	Description() string
	SetDescription(description string) MediaInterface

	EntityID() string
	SetEntityID(entityID string) MediaInterface

	ID() string
	SetID(id string) MediaInterface

	Memo() string
	SetMemo(memo string) MediaInterface

	Metas() (map[string]string, error)
	Meta(name string) string
	SetMeta(name string, value string) error
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	Sequence() int
	SetSequence(sequence int) MediaInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(softDeletedAt string) MediaInterface

	Status() string
	SetStatus(status string) MediaInterface

	Title() string
	SetTitle(title string) MediaInterface

	Type() string
	SetType(type_ string) MediaInterface

	URL() string
	SetURL(url string) MediaInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) MediaInterface
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
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) OrderInterface

	CustomerID() string
	SetCustomerID(customerID string) OrderInterface

	ID() string
	SetID(id string) OrderInterface

	Memo() string
	SetMemo(memo string) OrderInterface

	Meta(name string) string
	SetMeta(name string, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	Price() string
	SetPrice(price string) OrderInterface
	PriceFloat() float64
	SetPriceFloat(price float64) OrderInterface

	Quantity() string
	SetQuantity(quantity string) OrderInterface
	QuantityInt() int64
	SetQuantityInt(quantity int64) OrderInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) OrderInterface

	Status() string
	SetStatus(status string) OrderInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) OrderInterface
}

type OrderLineItemInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) OrderLineItemInterface

	ID() string
	SetID(id string) OrderLineItemInterface

	Memo() string
	SetMemo(memo string) OrderLineItemInterface

	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	Meta(name string) string
	SetMeta(name string, value string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	OrderID() string
	SetOrderID(orderID string) OrderLineItemInterface

	Price() string
	SetPrice(price string) OrderLineItemInterface

	PriceFloat() float64
	SetPriceFloat(price float64) OrderLineItemInterface

	ProductID() string
	SetProductID(productID string) OrderLineItemInterface

	Quantity() string
	SetQuantity(quantity string) OrderLineItemInterface

	QuantityInt() int64
	SetQuantityInt(quantity int64) OrderLineItemInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) OrderLineItemInterface

	Status() string
	SetStatus(status string) OrderLineItemInterface

	Title() string
	SetTitle(title string) OrderLineItemInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) OrderLineItemInterface
}

type ProductInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Methods

	IsActive() bool
	IsDisabled() bool
	IsDraft() bool
	IsSoftDeleted() bool
	IsFree() bool
	Slug() string

	// Setters and Getters

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) ProductInterface

	Description() string
	SetDescription(description string) ProductInterface

	ID() string
	SetID(id string) ProductInterface

	Memo() string
	SetMemo(memo string) ProductInterface

	Meta(name string) string
	SetMeta(name string, value string) error

	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	MetasUpsert(metas map[string]string) error
	MetaRemove(name string) error
	MetasRemove(names []string) error

	Price() string
	SetPrice(price string) ProductInterface
	PriceFloat() float64
	SetPriceFloat(price float64) ProductInterface

	Quantity() string
	SetQuantity(quantity string) ProductInterface
	QuantityInt() int64
	SetQuantityInt(quantity int64) ProductInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) ProductInterface

	ShortDescription() string
	SetShortDescription(shortDescription string) ProductInterface

	Status() string
	SetStatus(status string) ProductInterface

	Title() string
	SetTitle(title string) ProductInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) ProductInterface
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
}
