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

	categoryFound, errFind := store.CategoryFindByID(ctx, category.GetID())

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

	err = store.CategoryDeleteByID(ctx, category.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.CategoryFindByID(ctx, category.GetID())

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

	categoryFound, errFind := store.CategoryFindByID(ctx, category.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound == nil {
		t.Fatal("unexpected nil category")
	}

	if categoryFound.GetID() != category.GetID() {
		t.Fatal("unexpected category id")
	}

	if categoryFound.GetTitle() != category.GetTitle() {
		t.Fatal("unexpected category title")
	}

	if categoryFound.GetStatus() != category.GetStatus() {
		t.Fatal("unexpected category status")
	}

	if categoryFound.GetParentID() != category.GetParentID() {
		t.Fatal("unexpected category parent id")
	}

	if !strings.Contains(categoryFound.GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", categoryFound.GetSoftDeletedAt())
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

	categoryFound, errFind := store.CategoryFindByID(ctx, category.GetID())

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

	if list[0].GetID() != category.GetID() {
		t.Fatal("unexpected category id")
	}

	if strings.Contains(list[0].GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Category MUST be soft deleted, but found: ", list[0].GetSoftDeletedAt())
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

	err = store.CategorySoftDeleteByID(ctx, category.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	categoryFound, errFind := store.CategoryFindByID(ctx, category.GetID())

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

	if list[0].GetID() != category.GetID() {
		t.Fatal("unexpected category id")
	}

	if strings.Contains(list[0].GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Category MUST be soft deleted, but found: ", list[0].GetSoftDeletedAt())
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

	categoryFound, errFind := store.CategoryFindByID(ctx, category.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if categoryFound.GetTitle() != "CATEGORY_TITLE_UPDATED" {
		t.Fatal("unexpected category title: ", categoryFound.GetTitle())
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

	discountFound, errFind := store.DiscountFindByID(ctx, discount.GetID())

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

	err = store.DiscountDeleteByID(ctx, discount.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	discountFound, errFind := store.DiscountFindByID(ctx, discount.GetID())

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

	discountFound, errFind := store.DiscountFindByID(ctx, discount.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound == nil {
		t.Fatal("Discount MUST NOT be nil")
		return
	}

	if discountFound.GetTitle() != "DISCOUNT_TITLE" {
		t.Fatal("Exam title MUST BE 'DISCOUNT_TITLE', found: ", discountFound.GetTitle())
		return
	}

	if discountFound.GetDescription() != "DISCOUNT_DESCRIPTION" {
		t.Fatal("Exam description MUST BE 'DISCOUNT_DESCRIPTION', found: ", discountFound.GetDescription())
	}

	if discountFound.GetStatus() != DISCOUNT_STATUS_DRAFT {
		t.Fatal("Exam status MUST BE 'draft', found: ", discountFound.GetStatus())
		return
	}

	if discountFound.GetType() != DISCOUNT_TYPE_AMOUNT {
		t.Fatal("Exam type MUST BE 'amount', found: ", discountFound.GetType())
	}

	if discountFound.GetType() != DISCOUNT_TYPE_AMOUNT {
		t.Fatal("Exam type MUST BE 'amount', found: ", discountFound.GetType())
	}

	if discountFound.GetAmount() != 19.9900 {
		t.Fatal("Exam price MUST BE '19.9900', found: ", discountFound.GetAmount())
		return
	}

	if discountFound.GetStartsAt() != "2022-01-01 00:00:00 +0000 UTC" {
		t.Fatal("Exam start date MUST BE '2022-01-01 00:00:00', found: ", discountFound.GetStartsAt())
	}

	if discountFound.GetEndsAt() != "2022-01-01 23:59:59 +0000 UTC" {
		t.Fatal("Exam end date MUST BE '2022-01-01 23:59:59', found: ", discountFound.GetEndsAt())
	}

	if !strings.Contains(discountFound.GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", discountFound.GetSoftDeletedAt())
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

	discountFound, errFind := store.DiscountFindByID(ctx, discount.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound != nil {
		t.Fatal("Discount MUST be nil")
		return
	}

	discountList, errList := store.DiscountList(ctx, NewDiscountQuery().
		SetID(discount.GetID()).
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

	discountFound, errFind := store.DiscountFindByID(ctx, discount.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if discountFound == nil {
		t.Fatal("Discount MUST NOT be nil")
	}

	if discountFound.GetTitle() != "DISCOUNT_TITLE_UPDATED" {
		t.Fatal("Discount title MUST BE 'DISCOUNT_TITLE_UPDATED', found: ", discountFound.GetTitle())
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

	categoryFound, errFind := store.MediaFindByID(ctx, media.GetID())

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

	err = store.MediaDeleteByID(ctx, media.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	mediaFound, errFind := store.MediaFindByID(ctx, media.GetID())

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

	mediaFound, errFind := store.MediaFindByID(ctx, media.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if mediaFound == nil {
		t.Fatal("unexpected nil media")
	}

	if mediaFound.GetID() != media.GetID() {
		t.Fatal("unexpected media id")
	}

	if mediaFound.GetTitle() != media.GetTitle() {
		t.Fatal("unexpected media title")
	}

	if mediaFound.GetStatus() != media.GetStatus() {
		t.Fatal("unexpected category status")
	}

	if mediaFound.GetEntityID() != media.GetEntityID() {
		t.Fatal("unexpected category parent id")
	}

	if !strings.Contains(mediaFound.GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", mediaFound.GetSoftDeletedAt())
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

	mediaFound, errFind := store.MediaFindByID(ctx, media.GetID())

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

	if list[0].GetID() != media.GetID() {
		t.Fatal("unexpected media id")
	}

	if strings.Contains(list[0].GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Media MUST be soft deleted, but found: ", list[0].GetSoftDeletedAt())
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

	err = store.MediaSoftDeleteByID(ctx, media.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	mediaFound, errFind := store.MediaFindByID(ctx, media.GetID())

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

	if list[0].GetID() != media.GetID() {
		t.Fatal("unexpected media id")
	}

	if strings.Contains(list[0].GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Media MUST be soft deleted, but found: ", list[0].GetSoftDeletedAt())
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

	mediaFound, errFind := store.MediaFindByID(ctx, media.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if mediaFound.GetTitle() != "MEDIA_TITLE_UPDATED" {
		t.Fatal("unexpected media title: ", mediaFound.GetTitle())
	}
}
