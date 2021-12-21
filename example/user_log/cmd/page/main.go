package main

import (
	"context"
	"log"

	"gitlab.com/tz/tzui"
	"gitlab.com/tz/tzui/example/user_log/model/user_log"
	"gitlab.com/tz/tzui/pkg/devui"
	"gitlab.com/tz/tzui/pkg/utils"
)

func main() {
	pb := tzui.NewTzuiPageBuilder("", "user_log", "用户操作日志")
	pb.AddTzComponent(
		devui.NewDataTableBuilder(
			user_log.UserLog{},
			user_log.GetAll,
			devui.HeaderFixedDataTable,
		),
	)
	res, err := pb.Handle(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	byt, err := utils.JsonAPI.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(byt))
}
