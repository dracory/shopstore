package shopstore

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
	_ "modernc.org/sqlite"
)

func initDB(filepath string) (*sql.DB, error) {
	if filepath != ":memory:" {
		err := os.Remove(filepath) // remove database

		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			return nil, err
		}
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initStore(filepath string) (StoreInterface, error) {
	db, err := initDB(filepath)

	if err != nil {
		return nil, err
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		CategoryTableName:      "shop_category",
		DiscountTableName:      "shop_discount",
		MediaTableName:         "shop_media",
		OrderTableName:         "shop_order",
		OrderLineItemTableName: "shop_order_line_item",
		ProductTableName:       "shop_product",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	return store, nil
}

func TestStoreCategoryCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("CATEGORY_TITLE")

	err = store.CategoryCreate(database.Context(context.Background(), store.DB()), category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreCategoryDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("CATEGORY_TITLE")

	ctx := database.Context(context.Background(), store.DB())

	err = store.CategoryCreate(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategoryDelete(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.CategoryFindByID(ctx, category.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound != nil {
		t.Fatal("unexpected category found")
	}
}

func TestStoreCategoryDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("CATEGORY_TITLE")

	ctx := database.Context(context.Background(), store.DB())

	err = store.CategoryCreate(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategoryDeleteByID(ctx, category.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.CategoryFindByID(ctx, category.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound != nil {
		t.Fatal("unexpected category found")
	}
}

func TestStoreCategoryFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("CATEGORY_TITLE")

	ctx := database.Context(context.Background(), store.DB())

	err = store.CategoryCreate(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.CategoryFindByID(ctx, category.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound == nil {
		t.Fatal("unexpected nil category")
	}

	if categoryFound.ID() != category.ID() {
		t.Fatal("unexpected category id")
	}

	if categoryFound.Title() != category.Title() {
		t.Fatal("unexpected category title")
	}

	if categoryFound.Status() != category.Status() {
		t.Fatal("unexpected category status")
	}

	if categoryFound.ParentID() != category.ParentID() {
		t.Fatal("unexpected category parent id")
	}

	if !strings.Contains(categoryFound.SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", categoryFound.SoftDeletedAt())
		return
	}
}

func TestStoreCategorySoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("CATEGORY_TITLE")

	ctx := database.Context(context.Background(), store.DB())

	err = store.CategoryCreate(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategorySoftDelete(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.CategoryFindByID(ctx, category.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound != nil {
		t.Fatal("category must be nil as it was soft deleted")
	}

	list, err := store.CategoryList(ctx, NewCategoryQuery().SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(list) < 1 {
		t.Fatal("unexpected empty list")
	}

	if list[0].ID() != category.ID() {
		t.Fatal("unexpected category id")
	}

	if strings.Contains(list[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Category MUST be soft deleted, but found: ", list[0].SoftDeletedAt())
		return
	}
}

func TestStoreCategorySoftDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("CATEGORY_TITLE")

	ctx := database.Context(context.Background(), store.DB())

	err = store.CategoryCreate(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.CategorySoftDeleteByID(ctx, category.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.CategoryFindByID(ctx, category.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound != nil {
		t.Fatal("category must be nil as it was soft deleted")
	}

	list, err := store.CategoryList(ctx, NewCategoryQuery().SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(list) < 1 {
		t.Fatal("unexpected empty list")
	}

	if list[0].ID() != category.ID() {
		t.Fatal("unexpected category id")
	}

	if strings.Contains(list[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Category MUST be soft deleted, but found: ", list[0].SoftDeletedAt())
		return
	}
}

func TestStoreCategoryUpdate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	category := NewCategory().
		SetStatus(CATEGORY_STATUS_DRAFT).
		SetTitle("CATEGORY_TITLE")

	ctx := database.Context(context.Background(), store.DB())

	err = store.CategoryCreate(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	category.SetTitle("CATEGORY_TITLE_UPDATED")

	err = store.CategoryUpdate(ctx, category)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.CategoryFindByID(ctx, category.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound.Title() != "CATEGORY_TITLE_UPDATED" {
		t.Fatal("unexpected category title: ", categoryFound.Title())
	}
}

func TestStoreDiscountCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE")

	ctx := context.Background()
	err = store.DiscountCreate(ctx, discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreDiscountDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE")

	ctx := context.Background()
	err = store.DiscountCreate(ctx, discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.DiscountDelete(ctx, discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	discountFound, errFind := store.DiscountFindByID(ctx, discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound != nil {
		t.Fatal("Exam MUST be nil")
		return
	}
}

func TestStoreDiscountDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE")

	ctx := context.Background()
	err = store.DiscountCreate(ctx, discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.DiscountDeleteByID(ctx, discount.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	discountFound, errFind := store.DiscountFindByID(ctx, discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound != nil {
		t.Fatal("Exam MUST be nil")
		return
	}
}

func TestStoreDiscountFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE").
		SetDescription("DISCOUNT_DESCRIPTION").
		SetType(DISCOUNT_TYPE_AMOUNT).
		SetAmount(19.99).
		SetStartsAt(`2022-01-01 00:00:00`).
		SetEndsAt(`2022-01-01 23:59:59`)

	ctx := context.Background()
	err = store.DiscountCreate(ctx, discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	discountFound, errFind := store.DiscountFindByID(ctx, discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound == nil {
		t.Fatal("Discount MUST NOT be nil")
		return
	}

	if discountFound.Title() != "DISCOUNT_TITLE" {
		t.Fatal("Exam title MUST BE 'DISCOUNT_TITLE', found: ", discountFound.Title())
		return
	}

	if discountFound.Description() != "DISCOUNT_DESCRIPTION" {
		t.Fatal("Exam description MUST BE 'DISCOUNT_DESCRIPTION', found: ", discountFound.Description())
	}

	if discountFound.Status() != DISCOUNT_STATUS_DRAFT {
		t.Fatal("Exam status MUST BE 'draft', found: ", discountFound.Status())
		return
	}

	if discountFound.Type() != DISCOUNT_TYPE_AMOUNT {
		t.Fatal("Exam type MUST BE 'amount', found: ", discountFound.Type())
	}

	if discountFound.Type() != DISCOUNT_TYPE_AMOUNT {
		t.Fatal("Exam type MUST BE 'amount', found: ", discountFound.Type())
	}

	if discountFound.Amount() != 19.9900 {
		t.Fatal("Exam price MUST BE '19.9900', found: ", discountFound.Amount())
		return
	}

	if discountFound.StartsAt() != "2022-01-01 00:00:00 +0000 UTC" {
		t.Fatal("Exam start date MUST BE '2022-01-01 00:00:00', found: ", discountFound.StartsAt())
	}

	if discountFound.EndsAt() != "2022-01-01 23:59:59 +0000 UTC" {
		t.Fatal("Exam end date MUST BE '2022-01-01 23:59:59', found: ", discountFound.EndsAt())
	}

	if !strings.Contains(discountFound.SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", discountFound.SoftDeletedAt())
		return
	}
}

func TestStoreDiscountSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE")

	ctx := context.Background()
	err = store.DiscountCreate(ctx, discount)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.DiscountSoftDelete(ctx, discount)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	discountFound, errFind := store.DiscountFindByID(ctx, discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound != nil {
		t.Fatal("Discount MUST be nil")
		return
	}

	discountList, errList := store.DiscountList(ctx, NewDiscountQuery().
		SetID(discount.ID()).
		SetSoftDeletedIncluded(true))

	if errList != nil {
		t.Fatal("unexpected error:", errList)
		return
	}

	if len(discountList) != 1 {
		t.Fatal("Discount list MUST be 1")
		return
	}
}

func TestStoreDiscountUpdate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE").
		SetDescription("DISCOUNT_DESCRIPTION").
		SetType(DISCOUNT_TYPE_AMOUNT).
		SetAmount(19.99).
		SetStartsAt(`2022-01-01 00:00:00`).
		SetEndsAt(`2022-01-01 23:59:59`)

	ctx := context.Background()
	err = store.DiscountCreate(ctx, discount)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	discount.SetTitle("DISCOUNT_TITLE_UPDATED")

	err = store.DiscountUpdate(ctx, discount)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	discountFound, errFind := store.DiscountFindByID(ctx, discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if discountFound == nil {
		t.Fatal("Discount MUST NOT be nil")
	}

	if discountFound.Title() != "DISCOUNT_TITLE_UPDATED" {
		t.Fatal("Discount title MUST BE 'DISCOUNT_TITLE_UPDATED', found: ", discountFound.Title())
	}
}

func TestStoreMediaCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID("ENTITY_O1").
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/image.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreMediaDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID("ENTITY_O1").
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/image.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	ctx := context.Background()

	err = store.MediaCreate(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MediaDelete(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.MediaFindByID(ctx, media.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound != nil {
		t.Fatal("unexpected media found")
	}
}

func TestStoreMediaDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID("ENTITY_O1").
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/image.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	ctx := database.Context(context.Background(), store.DB())

	err = store.MediaCreate(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MediaDeleteByID(ctx, media.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	mediaFound, errFind := store.MediaFindByID(ctx, media.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if mediaFound != nil {
		t.Fatal("unexpected media found")
	}
}

func TestStoreMediaFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID("ENTITY_O1").
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/image.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	ctx := database.Context(context.Background(), store.DB())

	err = store.MediaCreate(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	mediaFound, errFind := store.MediaFindByID(ctx, media.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if mediaFound == nil {
		t.Fatal("unexpected nil media")
	}

	if mediaFound.ID() != media.ID() {
		t.Fatal("unexpected media id")
	}

	if mediaFound.Title() != media.Title() {
		t.Fatal("unexpected media title")
	}

	if mediaFound.Status() != media.Status() {
		t.Fatal("unexpected category status")
	}

	if mediaFound.EntityID() != media.EntityID() {
		t.Fatal("unexpected category parent id")
	}

	if !strings.Contains(mediaFound.SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", mediaFound.SoftDeletedAt())
		return
	}
}

func TestStoreMediaSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID("ENTITY_O1").
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/image.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	ctx := database.Context(context.Background(), store.DB())

	err = store.MediaCreate(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MediaSoftDelete(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	mediaFound, errFind := store.MediaFindByID(ctx, media.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if mediaFound != nil {
		t.Fatal("media must be nil as it was soft deleted")
	}

	list, err := store.MediaList(ctx, NewMediaQuery().SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(list) < 1 {
		t.Fatal("unexpected empty list")
	}

	if list[0].ID() != media.ID() {
		t.Fatal("unexpected media id")
	}

	if strings.Contains(list[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Media MUST be soft deleted, but found: ", list[0].SoftDeletedAt())
		return
	}
}

func TestStoreMediaSoftDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID("ENTITY_O1").
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/image.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	ctx := database.Context(context.Background(), store.DB())

	err = store.MediaCreate(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MediaSoftDeleteByID(ctx, media.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	mediaFound, errFind := store.MediaFindByID(ctx, media.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if mediaFound != nil {
		t.Fatal("category must be nil as it was soft deleted")
	}

	list, err := store.MediaList(ctx, NewMediaQuery().SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(list) < 1 {
		t.Fatal("unexpected empty list")
	}

	if list[0].ID() != media.ID() {
		t.Fatal("unexpected media id")
	}

	if strings.Contains(list[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Media MUST be soft deleted, but found: ", list[0].SoftDeletedAt())
		return
	}
}

func TestStoreMediaUpdate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	media := NewMedia().
		SetStatus(MEDIA_STATUS_DRAFT).
		SetEntityID("ENTITY_O1").
		SetTitle("MEDIA_TITLE").
		SetURL("https://example.com/image.jpg").
		SetType(MEDIA_TYPE_IMAGE_JPG).
		SetSequence(1)

	ctx := database.Context(context.Background(), store.DB())

	err = store.MediaCreate(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	media.SetTitle("MEDIA_TITLE_UPDATED")

	err = store.MediaUpdate(ctx, media)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	mediaFound, errFind := store.MediaFindByID(ctx, media.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if mediaFound.Title() != "MEDIA_TITLE_UPDATED" {
		t.Fatal("unexpected media title: ", mediaFound.Title())
	}
}
