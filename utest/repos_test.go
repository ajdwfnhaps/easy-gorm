package utest

import (
	"fmt"
	"reflect"
	"testing"

	easygorm "github.com/ajdwfnhaps/easy-gorm"
	"github.com/ajdwfnhaps/easy-gorm/icore"
	"github.com/ajdwfnhaps/easy-gorm/repos"

	"github.com/jinzhu/gorm"
	"go.uber.org/dig"
)

//TestRepos 测试查询
func TestRepos(t *testing.T) {

	//创建依赖注入容器
	container := buildContainer()

	//使用时
	err := container.Invoke(func(
		db *gorm.DB,
		reps IProvinceRepository,
	) {

		fmt.Println("gogogo")

		//声明条件集合(可包括排序)
		var filters []func(*gorm.DB) *gorm.DB

		//追加查询条件
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("province LIKE  ?", "广%")
		})

		//查询单条记录 Get方法第二个参数为查询字段，nil为查所有字段
		one := Province{}
		reps.Get(&one, nil, filters[0])
		fmt.Println(one)

		//通过主键查询单条记录
		one = Province{}
		reps.GetByKey(&one, nil, 35)
		fmt.Println(one)

		//追加排序条件
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Order("id desc")
		})

		//分页查询
		if ret, err := reps.GetPageList(1, 10, "id,provinceID,province", filters[1:]...); err != nil {
			t.Error(err)
		} else {
			fmt.Println(reflect.TypeOf(ret.Items))

			fmt.Println(ret)

			if list, ok := ret.Items.(*[]Province); ok {
				for _, v := range *list {
					fmt.Println(v)
				}
			}

		}

		//列表查询 第一个参数为查询的字段，第二个参数可传入条件集合
		if list, err := reps.GetList(nil); err != nil {
			t.Error(err)
		} else {
			if items, ok := list.(*[]Province); ok {
				fmt.Println("--------")
				for _, v := range *items {
					fmt.Println(v)
				}
			}
		}

	})

	if err != nil {
		t.Error(err)
	}

}

//TestCreate 测试增删改
func TestCreate(t *testing.T) {
	//创建依赖注入容器
	container := buildContainer()

	//使用时
	err := container.Invoke(func(
		reps IProvinceRepository,
	) {

		db := reps.GetDb()
		//创建
		province := Province{ProvinceID: "100", Province: "外星"}
		db.Create(&province)
		//修改
		db.Model(&province).Update("province", "地球")

		//删除
		db.Delete(Province{}, "provinceID=?", "100")
		// or
		db.Where(Province{ProvinceID: "100"}).Delete(Province{})
	})

	if err != nil {
		t.Error(err)
	}
}

func NewDb() (*gorm.DB, error) {

	db, _, err := easygorm.UseGorm(&easygorm.Option{
		Debug:        true,
		DBType:       "mysql",
		MaxLifetime:  7200, //设置连接可以重用的最长时间(单位：秒)
		MaxOpenConns: 150,  //设置数据库的最大打开连接数
		MaxIdleConns: 50,   //设置空闲连接池中的最大连接数
		MySQL: easygorm.MySQLConf{
			Host:     "127.0.0.1",
			Port:     13306,
			User:     "root",
			Password: "123456",
			DBName:   "test_db",
		},
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}

type IProvinceRepository interface {
	icore.IRepository
}

type ProvinceRepository struct {
	*repos.Repository
}

func (p *ProvinceRepository) Init() {
	p.Repository.Init(
		func() interface{} {
			return &Province{}
		},
		func() interface{} {
			return &[]Province{}
		})
}

func NewProvinceRepository(repBase *repos.Repository) *ProvinceRepository {
	p := &ProvinceRepository{repBase}
	p.Init()
	return p
}

func buildContainer() *dig.Container {
	// 创建依赖注入容器
	container := dig.New()

	container.Provide(NewDb)
	container.Provide(repos.NewRepository)
	container.Provide(func(m *repos.Repository) icore.IRepository { return m })

	container.Provide(NewProvinceRepository)
	container.Provide(func(m *ProvinceRepository) IProvinceRepository { return m })

	return container

}
