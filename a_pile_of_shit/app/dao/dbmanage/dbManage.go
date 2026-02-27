package dao

import (
	"a_pile_of_shit/app/dao/dbmanage/datamanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage"
	"a_pile_of_shit/app/model"
)

func getInformationOfDB() *model.DBInfo {
	Info := new(model.DBInfo)
	getInformationOfMysql(Info)
	getInformationOfRedis(Info)
	return Info
}
func ConnectDatabase() {
	Info := getInformationOfDB()
	usermanage.ConnectUserDB(Info)
	datamanage.ConnectDataDB(Info)
}
