package services

import (
	"context"
	"encoding/json"
	"fmt"

	"pfserver/config"
)

type N struct{}

func Notifs() N {
	return N{}
}

type ResultMember struct {
	Msg  string `json:"msg"`
	Seen bool   `json:"seen"`
}
type Result struct {
	Score  int          `json:"Score"`
	Member ResultMember `json:"Member"`
}

func (N) Get(ctx context.Context, userId int) ([]Result, error) {
	redis, _ := config.Redis().Client()

	zcmd := redis.ZRangeWithScores(ctx, fmt.Sprintf("%s:%d", config.RDS_USER_NOTIFS, userId), 0, -1)
	var data []Result
	res, err := zcmd.Result()
	jsn, _ := json.Marshal(res)

	json.Unmarshal(jsn, &data)

	return data, err
}

// func (N) MarkSeen(ctx context.Context, userId int) ([]Result, error) {
// 	redis, _ := config.Redis().Client()

// 	zcmd := redis.ZRem(ctx, fmt.Sprintf("%s:%d", config.RDS_USER_NOTIFS, userId), 0, -1)
// 	var data []Result
// 	res, err := zcmd.Result()
// 	jsn, _ := json.Marshal(res)

// 	json.Unmarshal(jsn, &data)

// 	return data, err
// }
