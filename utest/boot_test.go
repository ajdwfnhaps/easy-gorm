package utest

import (
	"fmt"
	"testing"

	easygorm "github.com/ajdwfnhaps/easy-gorm"
)

func TestUseGorm(t *testing.T) {
	//Option实际项目中从配置文件设置
	db, relCall, err := easygorm.UseGorm(&easygorm.Option{
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

	defer func() {
		if relCall != nil {
			relCall()
			fmt.Println("已关闭数据库链接")
		}
	}()

	if err != nil {
		t.Error(err)
	}

	fmt.Println("数据库链接已打开")

	var provinces []Province
	db.Find(&provinces)

	for _, v := range provinces {
		fmt.Println(v)
	}

}

// Province Province实体
type Province struct {
	ID         int    `gorm:"column:id;primary_key;"`             // 记录内码
	ProvinceID string `gorm:"column:provinceID;type:varchar(6);"` // 编号
	Province   string `gorm:"column:province;type:varchar(40);"`  // 名称
}

// TableName  Set the table name to be `hat_province`
func (a Province) TableName() string {
	return "hat_province"
}
