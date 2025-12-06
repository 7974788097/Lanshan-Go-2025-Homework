package main

import (
	"moddle/api"
	"moddle/dao"
)

func main() {
	err := dao.ReadFile()
	if err != nil {
		panic(err)
	}
	api.InitRouter_gin()
}
