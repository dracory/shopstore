package shopstore

import (
	"context"
	"errors"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	"github.com/samber/lo"
)

func (store *Store) MediaCount(ctx context.Context, options MediaQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.mediaQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	store.logSql("select", sqlStr, params...)

	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (store *Store) MediaCreate(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	media.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := media.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.mediaTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("insert", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	return nil
}

func (store *Store) MediaDelete(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	return store.MediaDeleteByID(ctx, media.ID())
}

func (store *Store) MediaDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.mediaTableName).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	return nil
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

	q, columns, err := store.mediaQuery(options)

	if err != nil {
		return nil, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Select(columns...).
		ToSQL()

	if errSql != nil {
		return nil, errSql
	}

	store.logSql("select", sqlStr, params...)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return nil, err
	}

	list := []MediaInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewMediaFromExistingData(modelMap)
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

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.mediaTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(media.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, params...)

	_, err = database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	media.MarkAsNotDirty()

	return nil
}

func (store *Store) mediaQuery(options MediaQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("category options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.mediaTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasEntityID() {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).Eq(options.EntityID()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasTitleLike() {
		q = q.Where(goqu.C(COLUMN_TITLE).ILike(options.TitleLike()))
	}

	if options.HasType() {
		q = q.Where(goqu.C(COLUMN_TYPE).Eq(options.Type()))
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted media requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
