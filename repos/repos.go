package repos

import (
	"errors"

	"github.com/ajdwfnhaps/easy-gorm/schema"

	"github.com/jinzhu/gorm"
)

//Repository 仓储基类
type Repository struct {
	db  *gorm.DB
	fn  func() interface{}
	fns func() interface{}
}

// NewRepository 创建实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

//GetDb 获取DB
func (rep *Repository) GetDb() *gorm.DB {
	return rep.db
}

//Init 初始化
func (rep *Repository) Init(fns ...func() interface{}) {
	for i, f := range fns {
		switch i {
		case 0:
			rep.fn = f
		case 1:
			rep.fns = f
		}

	}
}

// //Create 创建
// func (rep *Repository) Create(in interface{}) {
// 	dbSet := rep.db
// 	dbSet.Create(in)
// }

// //Update 更新
// func (rep *Repository) Update(attrs ...interface{}) error {
// 	if rep.fn == nil {
// 		return errors.New("Repository未初始化")
// 	}

// 	dbSet := rep.db
// 	m := rep.fn()
// 	dbSet.Model(m).Update(attrs)
// 	return nil
// }

//GetByKey 通过主键获取
func (rep *Repository) GetByKey(out interface{}, selectFields interface{}, key int) {
	query := rep.db
	query = query.Model(out)
	if selectFields != nil {
		query = query.Select(selectFields)
	}

	query.First(out, key)
}

//Get 获取
func (rep *Repository) Get(out interface{}, selectFields interface{}, funcs ...func(*gorm.DB) *gorm.DB) {
	query := rep.db
	for _, fun := range funcs {
		query = query.Scopes(fun)
	}

	query = query.Model(out)

	if selectFields != nil {
		query = query.Select(selectFields)
	}

	query.Limit(1).Find(out)

}

//GetList 获取列表
func (rep *Repository) GetList(selectFields interface{}, funcs ...func(*gorm.DB) *gorm.DB) (interface{}, error) {

	if rep.fns == nil {
		return nil, errors.New("Repository未初始化")
	}

	query := rep.db
	for _, fun := range funcs {
		query = query.Scopes(fun)
	}

	if selectFields != nil {
		query = query.Select(selectFields)
	}

	out := rep.fns()
	query.Find(out)

	return out, nil
}

//GetPageList 分页查询
func (rep *Repository) GetPageList(pageIndex int, pageSize int, selectFields interface{}, funcs ...func(*gorm.DB) *gorm.DB) (*schema.PagingList, error) {

	if rep.fns == nil {
		return nil, errors.New("Repository未初始化")
	}

	query := rep.db
	for _, fun := range funcs {
		query = query.Scopes(fun)
	}

	count := 0
	m := rep.fn()
	query.Model(m).Count(&count)

	// 如果分页大小小于0或者分页索引小于0，则不查询数据
	if pageSize < 0 || pageIndex < 0 {
		return nil, errors.New("请设置正确的pageIndex，pageSize")
	}

	if pageIndex > 0 && pageSize > 0 {
		query = query.Offset((pageIndex - 1) * pageSize)
	}
	if pageSize > 0 {
		query = query.Limit(pageSize)
	}

	if selectFields != nil {
		query = query.Select(selectFields)
	}

	out := rep.fns()
	query.Find(out)

	return &schema.PagingList{
		PageIndex:      pageIndex,
		PageSize:       pageSize,
		TotalCount:     count,
		PageTotalCount: getPageCount(count, pageSize),
		Items:          out,
	}, nil
}

func getPageCount(totalCount, pageSize int) int {
	if totalCount <= 0 {
		return 0
	}

	if mod := totalCount % pageSize; mod > 0 {
		return (totalCount / pageSize) + 1
	}
	return totalCount / pageSize
}
