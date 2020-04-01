package icore

import (
	"github.com/ajdwfnhaps/easy-gorm/schema"

	"github.com/jinzhu/gorm"
)

//IRepository 仓储接口
type IRepository interface {
	// //Create 创建
	// Create(in interface{})
	//GetDb 获取DB
	GetDb() *gorm.DB
	//GetByKey 通过主键获取
	GetByKey(out interface{}, selectFields interface{}, key int)
	//Get 获取
	Get(out interface{}, selectFields interface{}, funcs ...func(*gorm.DB) *gorm.DB)
	//GetList 获取列表
	GetList(selectFields interface{}, funcs ...func(*gorm.DB) *gorm.DB) (interface{}, error)
	//GetPageList 分页查询
	GetPageList(pageIndex int, pageSize int, selectFields interface{}, funcs ...func(*gorm.DB) *gorm.DB) (*schema.PagingList, error)
}
