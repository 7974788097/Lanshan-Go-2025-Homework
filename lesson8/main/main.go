package main

import (
	"moddle/api"
	"moddle/dao"
)

func main() {
	dao.ConnectDatabase()
	api.InitRouter_gin()
}

//代办：
//1.对各部分数据库操作添加日志记录
//2.将token的密码更换为更安全的方式获取，而不是存储在代码中
