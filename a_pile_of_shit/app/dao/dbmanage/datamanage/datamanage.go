package datamanage

import (
	"a_pile_of_shit/app/dao/dbmanage/datamanage/commentmanage"
	"a_pile_of_shit/app/dao/dbmanage/datamanage/messagemanage"
	"a_pile_of_shit/app/model"
)

func ConnectDataDB(info *model.DBInfo) {
	commentmanage.ConnectCommentDB(info)
	messagemanage.ConnectMessageDB(info)
}

type DataOperation interface {
	AddInfo() bool
	DelInfo() bool
	GetInfo() (bool, bool)
	ChangeInfo() (bool, bool)
}

func Add(Info DataOperation, bo *bool) DataOperation {
	*bo = Info.AddInfo()
	return Info
} //传入data包的struct和error判断的*bool
func Delete(Info DataOperation, bo *bool) DataOperation {
	*bo = Info.DelInfo()
	return Info
}
func Update(Info DataOperation, exist *bool, bo *bool) DataOperation {
	*exist, *bo = Info.ChangeInfo()
	return Info
}
func Get(Info DataOperation, exist *bool, bo *bool) DataOperation {
	*exist, *bo = Info.GetInfo()
	return Info
}
