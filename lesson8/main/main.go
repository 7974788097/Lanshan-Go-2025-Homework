package main

import (
	"moddle/api"
	"moddle/dao"
)

func main() {
	dao.ConnectDatabase()
	api.InitRouter_gin()
}
