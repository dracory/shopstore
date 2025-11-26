package shopstore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

func (store *Store) CategoryCount(ctx context.Context, options CategoryQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.categoryQuery(options)

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

func (store *Store) CategoryCreate(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	category.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	category.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := category.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.categoryTableName).
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

func (store *Store) CategoryDelete(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	return store.CategoryDeleteByID(ctx, category.ID())
}

func (store *Store) CategoryDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.categoryTableName).
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

func (store *Store) CategoryFindByID(ctx context.Context, id string) (CategoryInterface, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	q := NewCategoryQuery().SetID(id).SetLimit(1)

	list, err := store.CategoryList(ctx, q)

	if err != nil {
		return nil, err
	}

	if len(list) < 1 {
		return nil, nil
	}

	return list[0], nil
}

func (store *Store) CategoryList(ctx context.Context, options CategoryQueryInterface) ([]CategoryInterface, error) {
	err := options.Validate()

	if err != nil {
		return nil, err
	}

	q, columns, err := store.categoryQuery(options)

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

	list := []CategoryInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewCategoryFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) CategorySoftDelete(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.CategoryUpdate(ctx, category)
}

func (store *Store) CategorySoftDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	category, err := store.CategoryFindByID(ctx, id)

	if err != nil {
		return err
	}

	if category == nil {
		return nil
	}

	return store.CategorySoftDelete(ctx, category)
}

func (store *Store) CategoryUpdate(ctx context.Context, category CategoryInterface) (err error) {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := category.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.categoryTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(category.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err = database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	category.MarkAsNotDirty()

	return nil

}

func (store *Store) categoryQuery(options CategoryQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("category options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.categoryTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasParentID() {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasTitleLike() {
		q = q.Where(goqu.C(COLUMN_TITLE).ILike(options.TitleLike()))
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted blocks requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
