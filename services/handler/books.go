package handler

import (
	"MarsLuo/dao"
	"MarsLuo/models"
	"context"
	"encoding/json"

	pb "MarsLuo/services/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"strconv"
)

type Videos struct{}



func (v *Videos) Lists(ctx context.Context, request *pb.ListsRequest, response *pb.ListsResponse) error {

	var dbOp dao.DbVideo
	var models []models.table
	var a int64
	err := dbOp.Lists(request, &models, &a)

	if err == nil {
		response.Code = 0
		response.Msg = ""
		response.Total = int32(a)
	} else {
		response.Code = -1
		response.Msg = err.Error()
		response.Total = 0
		return err
	}
	response.Records = []*structpb.Struct{}
	for _, rec := range models {

		tmp := new(pb.Book)

		pos, _ := rec.RecommendPos.MarshalJSON()
		result := new([]int32)
		json.Unmarshal(pos, result)


		response.Lists = append(response.VideoLists, tmp)
	}
	return nil
}

func (v *Videos) Recommend(ctx context.Context, request *pb.RecommendRequest, rsp *pb.RecommendResponse) error {
	var dbOp dao.DbVideo
	var models models.HicourtVideo

	if err := dbOp.UpdateRecommend(request, models); err != nil {
		return err
	}
	return nil
}

func (v *Videos) Detail(ctx context.Context, request *pb.DetailRequest, rsp *pb.DetailResponse) error {
	var dbOp dao.DbVideo
	var models models.HicourtVideo
	if err := dbOp.Detail(request.Id, &models); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.Code = 0
	rsp.Msg = ""
	tmp := new(pb.Video)
	tmp.VideoName = models.VideoName
	tmp.GoalType = models.GoalType
	tmp.GoalTime = models.GoalTime
	tmp.GoalScore = models.GoalScore
	tmp.VideoUrl = models.VideoUrl
	tmp.Id = models.Id
	tmp.Appid = models.Appid
	tmp.Cid = models.CID
	tmp.CourtId = strconv.Itoa(int(models.CourtId))
	tmp.CoverImage = models.CoverImage
	rsp.Video = tmp
	return nil
}

func (v *Videos) GetArena(ctx context.Context, request *pb.GetArenaRequest, rsp *pb.GetArenaResponse) error {
	var dbOp dao.DbVideo
	var models []models.HicourtArena
	var total int64
	if err := dbOp.GetArenaAppId(request, &models, &total); err != nil {
		rsp.Code = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.Code = 0
	rsp.Msg = ""
	var Appids []string
	if total > 0 && len(models) > 0 {
		for _, model := range models {
			Appids = append(Appids, model.Appid)
		}
	}
	rsp.Total = int32(total)
	rsp.AppIds = Appids
	return nil
}