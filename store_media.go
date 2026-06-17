package shopstore

import (
	"context"
	"errors"
	"strings"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (store *Store) MediaCount(ctx context.Context, options MediaQueryInterface) (int64, error) {
	q, err := store.mediaQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (store *Store) MediaCreate(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	media.SetSoftDeletedAt(MAX_DATETIME)

	data := media.Data()
	row := map[string]any{}
	for k, v := range data {
		row[k] = v
	}

	err := store.db.Query().Table(store.mediaTableName).Create(row)
	if err != nil {
		return err
	}

	return nil
}

func (store *Store) MediaDelete(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	return store.MediaDeleteByID(ctx, media.GetID())
}

func (store *Store) MediaDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	_, err := store.db.Query().Table(store.mediaTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}

func (store *Store) MediaFindByID(ctx context.Context, id string) (MediaInterface, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	q := NewMediaQuery().SetID(id).SetLimit(1)

	list, err := store.MediaList(ctx, q)

	if err != nil {
		return nil, err
	}

	if len(list) < 1 {
		return nil, nil
	}

	return list[0], nil
}

func (store *Store) MediaList(ctx context.Context, options MediaQueryInterface) ([]MediaInterface, error) {
	err := options.Validate()

	if err != nil {
		return nil, err
	}

	q, err := store.mediaQuery(options)

	if err != nil {
		return nil, err
	}

	var results []map[string]any
	if err := q.Get(&results); err != nil {
		return nil, err
	}

	list := []MediaInterface{}

	lo.ForEach(results, func(result map[string]any, index int) {
		model := NewMediaFromExistingData(mapAnyToString(result))
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) MediaSoftDelete(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.MediaUpdate(ctx, media)
}

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

func (store *Store) MediaUpdate(ctx context.Context, media MediaInterface) (err error) {
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := media.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	row := map[string]any{}
	for k, v := range dataChanged {
		row[k] = v
	}

	_, err = store.db.Query().Table(store.mediaTableName).Where(COLUMN_ID+" = ?", media.GetID()).Update(row)

	if err != nil {
		return err
	}

	media.MarkAsNotDirty()

	return nil
}

func (store *Store) mediaQuery(options MediaQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		return nil, errors.New("category options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.mediaTableName)

	if options.HasID() {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasIDIn() {
		ids := make([]any, len(options.IDIn()))
		for i, id := range options.IDIn() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.HasEntityID() {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID())
	}

	if options.HasStatus() {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.HasTitleLike() {
		searchTerm := strings.ReplaceAll(options.TitleLike(), "'", "''")
		searchTerm = strings.ReplaceAll(searchTerm, "%", "\\%")
		searchTerm = strings.ReplaceAll(searchTerm, "_", "\\_")
		q = q.Where(COLUMN_TITLE+" LIKE ?", "%"+searchTerm+"%")
	}

	if options.HasType() {
		q = q.Where(COLUMN_MEDIA_TYPE+" = ?", options.Type())
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(cast.ToInt(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(cast.ToInt(options.Offset()))
		}
	}

	sortOrder := lo.Ternary(options.HasSortDirection(), options.SortDirection(), "desc")

	if options.HasOrderBy() {
		q = q.OrderBy(options.OrderBy(), sortOrder)
	}

	if options.SoftDeletedIncluded() {
		return q, nil // soft deleted media requested specifically
	}

	q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", MAX_DATETIME)

	return q, nil
}
