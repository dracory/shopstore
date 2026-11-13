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

func (store *Store) SettingCount(ctx context.Context, options SettingQueryInterface) (int64, error) {
	q, err := store.settingQuery(options)
	if err != nil {
		return -1, err
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (store *Store) SettingCreate(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	if setting.GetKey() == "" {
		return errors.New("setting key is empty")
	}

	setting.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	setting.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	setting.SetSoftDeletedAt(MAX_DATETIME)

	data := setting.Data()
	row := map[string]any{}
	for k, v := range data {
		row[k] = v
	}

	err := store.db.Query().Table(store.settingTableName).Create(row)
	if err != nil {
		return err
	}

	setting.MarkAsNotDirty()

	return nil
}

func (store *Store) SettingDelete(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	return store.SettingDeleteByID(ctx, setting.GetID())
}

func (store *Store) SettingDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("setting id is empty")
	}

	_, err := store.db.Query().Table(store.settingTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}

func (store *Store) SettingDeleteByKey(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("setting key is empty")
	}

	_, err := store.db.Query().Table(store.settingTableName).Where(COLUMN_SETTING_KEY+" = ?", key).Delete()
	return err
}

func (store *Store) SettingFindByID(ctx context.Context, id string) (SettingInterface, error) {
	if id == "" {
		return nil, errors.New("setting id is empty")
	}

	list, err := store.SettingList(ctx, NewSettingQuery().
		SetID(id).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) SettingFindByKey(ctx context.Context, key string) (SettingInterface, error) {
	if key == "" {
		return nil, errors.New("setting key is empty")
	}

	list, err := store.SettingList(ctx, NewSettingQuery().
		SetKey(key).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// SettingGet is a shortcut method to get a setting value by key, or a
// default if not found.
func (store *Store) SettingGet(ctx context.Context, key string, valueDefault string) (string, error) {
	setting, err := store.SettingFindByKey(ctx, key)

	if err != nil {
		return valueDefault, err
	}

	if setting == nil {
		return valueDefault, nil
	}

	return setting.GetValue(), nil
}

// SettingSet is a shortcut method to save a setting value by key.
// Creates the setting if it does not exist, updates it otherwise.
func (store *Store) SettingSet(ctx context.Context, key string, value string) error {
	setting, err := store.SettingFindByKey(ctx, key)

	if err != nil {
		return err
	}

	if setting == nil {
		setting = NewSetting().SetKey(key).SetValue(value)
		return store.SettingCreate(ctx, setting)
	}

	setting.SetValue(value)

	return store.SettingUpdate(ctx, setting)
}

func (store *Store) SettingList(ctx context.Context, options SettingQueryInterface) ([]SettingInterface, error) {
	q, err := store.settingQuery(options)
	if err != nil {
		return []SettingInterface{}, err
	}

	var results []map[string]any
	if err := q.Get(&results); err != nil {
		return []SettingInterface{}, err
	}

	list := []SettingInterface{}

	lo.ForEach(results, func(result map[string]any, index int) {
		model := NewSettingFromExistingData(mapAnyToString(result))
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) SettingSoftDelete(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	setting.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SettingUpdate(ctx, setting)
}

func (store *Store) SettingSoftDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("setting id is empty")
	}

	setting, err := store.SettingFindByID(ctx, id)
	if err != nil {
		return err
	}
	if setting == nil {
		return nil
	}

	return store.SettingSoftDelete(ctx, setting)
}

func (store *Store) SettingUpdate(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	setting.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := setting.DataChanged()

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

	_, err := store.db.Query().Table(store.settingTableName).Where(COLUMN_ID+" = ?", setting.GetID()).Update(row)

	setting.MarkAsNotDirty()

	return err
}

func (store *Store) settingQuery(options SettingQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		options = NewSettingQuery()
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.settingTableName)

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

	if options.HasKey() {
		q = q.Where(COLUMN_SETTING_KEY+" = ?", options.Key())
	}

	if options.HasKeyLike() {
		searchTerm := strings.ReplaceAll(options.KeyLike(), "'", "''")
		searchTerm = strings.ReplaceAll(searchTerm, "%", "\\%")
		searchTerm = strings.ReplaceAll(searchTerm, "_", "\\_")
		q = q.Where(COLUMN_SETTING_KEY+" LIKE ?", "%"+searchTerm+"%")
	}

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" BETWEEN ? AND ?", options.CreatedAtGte(), options.CreatedAtLte())
	} else if options.HasCreatedAtGte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	} else if options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
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

	if !options.SoftDeletedIncluded() {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", MAX_DATETIME)
	}

	return q, nil
}
