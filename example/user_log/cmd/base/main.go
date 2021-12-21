package main

import (
	"context"
	"log"

	"gitlab.com/tz/tzui/example/user_log/model/user_log"
	"gitlab.com/tz/tzui/pkg/devui"
	"gitlab.com/tz/tzui/pkg/utils"
)

func main() {
	req := new(devui.DataTableSourceRequest)
	req.Page = 1
	req.PerPage = 5
	res, err := user_log.GetAll(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	byt, err := utils.JsonAPI.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(byt))
}
