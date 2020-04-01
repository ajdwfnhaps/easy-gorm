package easygorm

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	// gorm存储注入
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// NewDB 创建DB实例
func NewDB(c *Option) (*gorm.DB, error) {
	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return db, nil
}

// UseGorm gorm with mysql
func UseGorm(opt *Option) (*gorm.DB, func(), error) {

	var dsn string
	switch opt.DBType {
	case "mysql":
		dsn = opt.MySQL.DSN()
	default:
		return nil, nil, errors.New("unknown db")
	}

	if len(opt.DSN) == 0 {
		opt.DSN = dsn
	}

	db, err := NewDB(opt)

	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		db.Close()
	}, err

}
