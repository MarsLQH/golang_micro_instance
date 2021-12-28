package base

import (
	"MarsLuo/core"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"sync"
)

type Functions struct {
	core     core.DbService
	DbClient *gorm.DB
}

var funcs Functions

func (f *Functions) GetDb() (client *gorm.DB, err error) {
	if f.DbClient != nil {
		return f.DbClient, nil
	}

	once := &sync.Once{}
	once.Do(func() {
		log.Info().Msg("	Dbbook init ")
		DbClient, err := f.core.ConnectDB()
		if err != nil {
			log.Error().Err(err).Msg("链接数据库时发生错误")

		}
		f.DbClient = DbClient
	})
	if f.DbClient == nil {
		return nil, errors.New("链接数据库时发生错误")
	}
	return f.DbClient, nil
}
