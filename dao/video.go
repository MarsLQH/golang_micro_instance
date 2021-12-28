package dao

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"git.myarena7.com/arena/hicourt/dao/base"
	"git.myarena7.com/arena/hicourt/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"

	pb "git.myarena7.com/arena/hicourt/services/proto"
)

type DbVideo struct {
	Db   *gorm.DB
	Base base.Functions
}

var dbVideo DbVideo

func (v *DbVideo) Lists(request *pb.ListsRequest, videoModel *[]models.HicourtVideo, a *int64) (err error) {

	log.Info().Msgf("请求参数%v", request)
	//err = v.GetDb().Model(&models.HicourtVideo{}).Count(a).Error ok
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
	//tmp = v.Db.Debug().Order("id desc")
	tmp = v.Db.Debug()

	if request.OrderBy != "" {
		tmp = tmp.Order(request.OrderBy + " desc,goal_time desc")
	} else {
		if request.Pos != 0 {
			tmp = tmp.Order("recommend_at desc")
		} else {
			tmp = tmp.Order("goal_time desc")
		}
	}
	if request.PerPage > 0 {
		perPage = request.PerPage
	}
	if request.Page > 1 {
		page = request.Page
	}
	offSet = (page - 1) * perPage
	//tmp = tmp.Offset(int(offSet)).Limit(int(perPage))

	if request.Pos != 0 {
		str := " JSON_CONTAINS(recommend_pos,JSON_Array(?)) "
		tmp = tmp.Where(str, request.Pos)
		//tmp2 = tmp2.Where(str, request.Pos)
		//err =  v.Db.Debug().Find(&videoModel,datatypes.JSONQuery("recommend_pos").Equals(request.Pos,"pos")).Error	//错误示例 can not select correct records
		//err = v.Db.Debug().Where(" JSON_CONTAINS(recommend_pos,JSON_Array(?)) ", request.Pos).Find(&videoModel).Error	//debug可以查出想要的数据sql this can  select correct records
		//JSON_CONTAINS(recommend_pos,JSON_Array( 7))
		//  v.Db.Where(" JSON_CONTAINS(recommend_pos,JSON_Array(?)) ",request.Pos)
	}
	if request.AppIds != nil {
		str := "appid in (?) "
		tmp = tmp.Where(str, request.AppIds)
		//tmp2 = tmp2.Where(str, request.AppIds)
	}
	if request.Cids != nil {
		str := "cid in (?) "
		tmp = tmp.Where(str, request.Cids)
		//tmp2 = tmp2.Where(str, request.Cids)
	}
	if request.CourtId != "" {
		tmp = tmp.Where("court_id", request.CourtId)
		//tmp2 = tmp2.Where("court_id", request.CourtId)
	}
	if request.GoalType != "" {
		tmp = tmp.Where("goal_type", request.GoalType)
		//tmp2 = tmp2.Where("goal_type", request.GoalType)
	}

	//if request.BeginAt.IsValid() {
	if request.BeginAt != "" {
		tmp = tmp.Where("goal_time >= ?", request.BeginAt)
		//tmp = tmp2.Where("goal_time >= ?", request.BeginAt)
	}
	if request.EndAt != "" {
		tmp = tmp.Where("goal_time <= ?", request.EndAt)
	}
	if request.UserId != "" {
		tmp = tmp.Joins("INNER JOIN `hicourt_user_video` `UserVideo` ON `UserVideo`.`hicourt_video_id` = `hicourt_video`.`id`").Where("`UserVideo`.`user_id`", request.UserId)
	}

	//======get total
	var model models.HicourtVideo
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
	//log.Info().Msgf("返回数据：%v", videoModel)
	return nil
}

func (v *DbVideo) UpdateRecommend(request *pb.RecommendRequest, videoModel models.HicourtVideo) (err error) {
	v.Db, err = v.Base.GetDb()
	sqlDB, err := v.Db.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Error().Err(err).Msgf("error getDb=%v", err.Error())
		return err
	}

	recommendAt := &sql.NullTime{Valid: false}
	var posJson []byte
	if request.Pos != nil {
		fmt.Println("===============pos=",request.Pos)
		recommendAt = &sql.NullTime{Time: time.Now(), Valid: true}
		posJson, err = json.Marshal(request.Pos)

		if err != nil {
			log.Error().Err(err).Msgf("Pos 参数错误 %v", request.Pos)
			return err
		}
	}else {
		posJson = nil
	}

	if request.Id <= 0 {
		log.Error().Err(err).Msgf("ID 参数错误 ID 为空 %v", request.Id)
		return errors.New("参数错误ID为空")
	}

	log.Info().Msgf("请求参数%v", request)
	log.Info().Msgf("posJson1 = %v", posJson)
	if err = v.Db.Debug().Model(&videoModel).Where("id=?", request.Id).Updates(models.HicourtVideo{RecommendPos: nil, RecommendAt:*recommendAt }).Error; err != nil {
		log.Error().Err(err).Msgf("db exec %v", request)
		return err
	}
	return nil
}

func (v *DbVideo) Detail(id int64, videoModel *models.HicourtVideo) (err error) {
	v.Db, err = v.Base.GetDb()
	sqlDB, err := v.Db.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Error().Err(err).Msgf("error getDb=%v", err.Error())
		return err
	}
	if err = v.Db.Debug().Find(&videoModel, id).Error; err != nil {
		log.Error().Err(err).Msgf("exec sql error =%v,id=%d", err.Error(), id)
		return err
	}
	log.Info().Msgf("exec sql result =%v,id=%d", videoModel, id)
	return nil
}

func (v *DbVideo) GetArenaAppId(request *pb.GetArenaRequest, videoModel *[]models.HicourtArena, a *int64) (err error) {
	log.Info().Msgf("请求参数%v", request)

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

	tmp = v.Db.Debug().Distinct("appid")
	offSet = (page - 1) * perPage

	if request.Cids != nil {
		str := "company_id in (?) "
		tmp = tmp.Where(str, request.Cids)
	}
	//======get total
	var model models.HicourtArena
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
	log.Info().Msgf("返回数据：%v", videoModel)

	return nil
}
