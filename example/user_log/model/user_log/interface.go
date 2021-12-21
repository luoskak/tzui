package user_log

import (
	"context"
	"time"

	"gitlab.com/tz/tzui/pkg/devui"
)

// UserLog 用户日志
// @tzui:"form:filter:LogType"
type UserLog struct {
	ID         int       `tzui:"header:ID"`
	UserID     int       `tzui:"header:用户ID"`
	LogContent string    `tzui:"header:日志"`
	CreateAt   time.Time `tzui:"header:日期"`
	LogType    int       `json:"-"`
}

func GetAll(ctx context.Context, req *devui.DataTableSourceRequest) (*devui.DataTableSourceResponse, error) {
	time1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2021-12-10T14:45:00Z08:00")
	return req.Slice([]*UserLog{
		{
			ID:         0,
			UserID:     1,
			LogContent: "登录",
			CreateAt:   time1.Add(time.Minute * 20),
			LogType:    0,
		},
		{
			ID:         0,
			UserID:     1,
			LogContent: "查询",
			CreateAt:   time1.Add(time.Hour),
			LogType:    1,
		},
	}), nil
}
