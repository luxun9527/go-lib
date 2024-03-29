// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"go-lib/sdk/gorm/gen/dao/model"
)

func newAsset(db *gorm.DB, opts ...gen.DOOption) asset {
	_asset := asset{}

	_asset.assetDo.UseDB(db, opts...)
	_asset.assetDo.UseModel(&model.Asset{})

	tableName := _asset.assetDo.TableName()
	_asset.ALL = field.NewAsterisk(tableName)
	_asset.ID = field.NewInt32(tableName, "id")
	_asset.UserID = field.NewInt32(tableName, "user_id")
	_asset.Username = field.NewString(tableName, "username")
	_asset.CoinID = field.NewInt32(tableName, "coin_id")
	_asset.CoinName = field.NewString(tableName, "coin_name")
	_asset.AvailableQty = field.NewString(tableName, "available_qty")
	_asset.FrozenQty = field.NewString(tableName, "frozen_qty")
	_asset.CreatedAt = field.NewTime(tableName, "created_at")
	_asset.UpdatedAt = field.NewTime(tableName, "updated_at")

	_asset.fillFieldMap()

	return _asset
}

type asset struct {
	assetDo assetDo

	ALL          field.Asterisk
	ID           field.Int32
	UserID       field.Int32  // 用户ID
	Username     field.String // 用户名
	CoinID       field.Int32  // 数字货币ID
	CoinName     field.String // 数字货币名称
	AvailableQty field.String // 可用余额
	FrozenQty    field.String // 冻结金额
	CreatedAt    field.Time   // 创建时间
	UpdatedAt    field.Time   // 修改时间

	fieldMap map[string]field.Expr
}

func (a asset) Table(newTableName string) *asset {
	a.assetDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a asset) As(alias string) *asset {
	a.assetDo.DO = *(a.assetDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *asset) updateTableName(table string) *asset {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewInt32(table, "id")
	a.UserID = field.NewInt32(table, "user_id")
	a.Username = field.NewString(table, "username")
	a.CoinID = field.NewInt32(table, "coin_id")
	a.CoinName = field.NewString(table, "coin_name")
	a.AvailableQty = field.NewString(table, "available_qty")
	a.FrozenQty = field.NewString(table, "frozen_qty")
	a.CreatedAt = field.NewTime(table, "created_at")
	a.UpdatedAt = field.NewTime(table, "updated_at")

	a.fillFieldMap()

	return a
}

func (a *asset) WithContext(ctx context.Context) *assetDo { return a.assetDo.WithContext(ctx) }

func (a asset) TableName() string { return a.assetDo.TableName() }

func (a asset) Alias() string { return a.assetDo.Alias() }

func (a asset) Columns(cols ...field.Expr) gen.Columns { return a.assetDo.Columns(cols...) }

func (a *asset) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *asset) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 9)
	a.fieldMap["id"] = a.ID
	a.fieldMap["user_id"] = a.UserID
	a.fieldMap["username"] = a.Username
	a.fieldMap["coin_id"] = a.CoinID
	a.fieldMap["coin_name"] = a.CoinName
	a.fieldMap["available_qty"] = a.AvailableQty
	a.fieldMap["frozen_qty"] = a.FrozenQty
	a.fieldMap["created_at"] = a.CreatedAt
	a.fieldMap["updated_at"] = a.UpdatedAt
}

func (a asset) clone(db *gorm.DB) asset {
	a.assetDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a asset) replaceDB(db *gorm.DB) asset {
	a.assetDo.ReplaceDB(db)
	return a
}

type assetDo struct{ gen.DO }

func (a assetDo) Debug() *assetDo {
	return a.withDO(a.DO.Debug())
}

func (a assetDo) WithContext(ctx context.Context) *assetDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a assetDo) ReadDB() *assetDo {
	return a.Clauses(dbresolver.Read)
}

func (a assetDo) WriteDB() *assetDo {
	return a.Clauses(dbresolver.Write)
}

func (a assetDo) Session(config *gorm.Session) *assetDo {
	return a.withDO(a.DO.Session(config))
}

func (a assetDo) Clauses(conds ...clause.Expression) *assetDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a assetDo) Returning(value interface{}, columns ...string) *assetDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a assetDo) Not(conds ...gen.Condition) *assetDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a assetDo) Or(conds ...gen.Condition) *assetDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a assetDo) Select(conds ...field.Expr) *assetDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a assetDo) Where(conds ...gen.Condition) *assetDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a assetDo) Order(conds ...field.Expr) *assetDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a assetDo) Distinct(cols ...field.Expr) *assetDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a assetDo) Omit(cols ...field.Expr) *assetDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a assetDo) Join(table schema.Tabler, on ...field.Expr) *assetDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a assetDo) LeftJoin(table schema.Tabler, on ...field.Expr) *assetDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a assetDo) RightJoin(table schema.Tabler, on ...field.Expr) *assetDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a assetDo) Group(cols ...field.Expr) *assetDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a assetDo) Having(conds ...gen.Condition) *assetDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a assetDo) Limit(limit int) *assetDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a assetDo) Offset(offset int) *assetDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a assetDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *assetDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a assetDo) Unscoped() *assetDo {
	return a.withDO(a.DO.Unscoped())
}

func (a assetDo) Create(values ...*model.Asset) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a assetDo) CreateInBatches(values []*model.Asset, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a assetDo) Save(values ...*model.Asset) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a assetDo) First() (*model.Asset, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Asset), nil
	}
}

func (a assetDo) Take() (*model.Asset, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Asset), nil
	}
}

func (a assetDo) Last() (*model.Asset, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Asset), nil
	}
}

func (a assetDo) Find() ([]*model.Asset, error) {
	result, err := a.DO.Find()
	return result.([]*model.Asset), err
}

func (a assetDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Asset, err error) {
	buf := make([]*model.Asset, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a assetDo) FindInBatches(result *[]*model.Asset, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a assetDo) Attrs(attrs ...field.AssignExpr) *assetDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a assetDo) Assign(attrs ...field.AssignExpr) *assetDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a assetDo) Joins(fields ...field.RelationField) *assetDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a assetDo) Preload(fields ...field.RelationField) *assetDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a assetDo) FirstOrInit() (*model.Asset, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Asset), nil
	}
}

func (a assetDo) FirstOrCreate() (*model.Asset, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Asset), nil
	}
}

func (a assetDo) FindByPage(offset int, limit int) (result []*model.Asset, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a assetDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a assetDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a assetDo) Delete(models ...*model.Asset) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *assetDo) withDO(do gen.Dao) *assetDo {
	a.DO = *do.(*gen.DO)
	return a
}
