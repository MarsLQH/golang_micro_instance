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
	var models []models.Table
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


		response.Lists = append(response.Lists, tmp)
	}
	return nil
}