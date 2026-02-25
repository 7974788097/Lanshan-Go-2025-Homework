package dao

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
	"fmt"
)

func getInformationOfMysql(Info *model.DBInfo) {
	password := other.GetenvDefault("MYSQL_ROOT_PASSWORD", "123456")
	username := other.GetenvDefault("MYSQL_USER", "root")
	userdbname := other.GetenvDefault("MYSQL_DBNAME_Users", "Users")
	datadbname := other.GetenvDefault("MYSQL_DBNAME_DATA", "Data")
	port := other.GetenvDefault("MYSQL_PORT", "3306")
	host := other.GetenvDefault("MYSQL_HOST", "localhost")
	Info.MysqlUserPath = fmt.Sprint(username, ":", password, "@tcp(", host, ":", port, ")/", userdbname, "?charset=utf8mb4&parseTime=True&loc=Local")
	Info.MysqlDataPath = fmt.Sprint(username, ":", password, "@tcp(", host, ":", port, ")/", datadbname, "?charset=utf8mb4&parseTime=True&loc=Local")
}
