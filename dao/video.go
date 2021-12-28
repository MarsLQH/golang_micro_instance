package dao

import (
	"MarsLuo/dao/base"
	"MarsLuo/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"


	pb "MarsLuo/services/proto"
)

type DbVideo struct {
	Db   *gorm.DB
	Base base.Functions
}

var dbVideo DbVideo

func (v *DbVideo) Lists(request *pb.ListsRequest, videoModel *[]models.Table, a *int64) (err error) {

	v.Db, err = v.Base.GetDb()
	sqlDB, err := v.Db.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Error().Err(err).Msgf("error happen when get db select=%v", err)
		return err
	}

	var offSet int32
	var perPage int32 = 20
	var page int32 = 1
	var tmp *gorm.DB
	tmp = v.Db.Debug()

	if request.OrderBy != "" {
		tmp = tmp.Order(request.OrderBy + " desc,cols1 desc")
	} else {
		if request.cols2 != 0 {
			tmp = tmp.Order("cols2 desc")
		} else {
			tmp = tmp.Order("cols3 desc")
		}
	}
	if request.PerPage > 0 {
		perPage = request.PerPage
	}
	if request.Page > 1 {
		page = request.Page
	}
	offSet = (page - 1) * perPage

	if request.Pos != 0 {
		str := " JSON_CONTAINS(recommend_pos,JSON_Array(?)) "
		tmp = tmp.Where(str, request.Pos)

		//err =  v.Db.Debug().Find(&videoModel,datatypes.JSONQuery("recommend_pos").Equals(request.Pos,"pos")).Error	//错误示例 can not select correct records
		//err = v.Db.Debug().Where(" JSON_CONTAINS(recommend_pos,JSON_Array(?)) ", request.Pos).Find(&BookModel).Error	//debug可以查出想要的数据sql this can  select correct records
		//JSON_CONTAINS(recommend_pos,JSON_Array( 7))
		//  v.Db.Where(" JSON_CONTAINS(recommend_pos,JSON_Array(?)) ",request.Pos)
	}


	var model models.Table
	err = tmp.Model(&model).Count(a).Error
	if err != nil {
		log.Error().Err(err).Msgf("取总数时错误")
		return err
	}
	tmp = tmp.Offset(int(offSet)).Limit(int(perPage))
	err = tmp.Find(&videoModel).Error
	if err != nil {
		log.Error().Err(err).Msgf("error happen when exec db select=%v", request)
		return err
	}
	return nil
}