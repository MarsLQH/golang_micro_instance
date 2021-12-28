package core

import (
	"MarsLuo/config"
	"errors"
	"reflect"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DbService struct {
}

func (*DbService) ConnectDB() (dbConn *gorm.DB, err error) {
	conf := config.C().Mysql

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: conf.User + ":" + conf.Pwd + "@tcp(" + conf.Host + ":" + conf.Port + ")/" +
			conf.DB + "?charset=utf8mb4&parseTime=True&loc=Local",
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "", SingularTable: true}})
	if err != nil {
		log.Err(err).Msg("connect db error")
		return nil, err
	}
	sqlDB, err := db.DB()
	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	sqlDB.SetMaxIdleConns(10)

	//sqlDB.Close()

	return db, nil
}

func CreateOrUpdate(db *gorm.DB, model interface{}, where interface{}, update interface{}) (interface{}, error) {
	var result interface{}
	err := db.Model(model).Where(where).First(result).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			//insert
			if err = db.Model(model).Create(update).Error; err != nil {
				return nil, err
			}
		}
	}
	//not update some field
	reflect.ValueOf(update).Elem().FieldByName("id").SetInt(0)
	if err = db.Model(model).Where(where).Updates(update).Error; err != nil {
		return nil, err
	}
	return update, nil
}
