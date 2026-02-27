package messagemanage

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
)

var dbContext *model.DBcontext

func ConnectMessageDB(info *model.DBInfo) {
	dbContext = &model.DBcontext{}
	connectMessageMysql(info)
}

type MessageInfo struct {
	model.MessageInfo
	GetWay string
	Info   []model.MessageInfo
}

func Make(info *map[string]any) *MessageInfo {
	newInfo := &MessageInfo{}
	if (*info)["message_id"] != nil {
		newInfo.MessageID = int64((*info)["message_id"].(float64))
	}
	if (*info)["message_name"] != nil {
		newInfo.MessageName = (*info)["message_name"].(string)
	}
	if (*info)["message_content"] != nil {
		newInfo.Content = (*info)["message_content"].(string)
	}
	if (*info)["username"] != nil {
		newInfo.Username = (*info)["username"].(string)
	}
	if (*info)["message_permission"] != nil {
		newInfo.Permission = uint8((*info)["message_permission"].(float64))
	} else {
		newInfo.Permission = 1
	}
	return newInfo
}
func (Message *MessageInfo) AddInfo() bool {
	var dbError model.DbError

	dbError.MysqlError = SetMassageOfMysql(&(Message.MessageInfo))

	return other.ErrorOutput(&dbError)
} //添加文章，返回是否发生错误

func (Message *MessageInfo) ChangeInfo() (bool, bool) {
	var dbError model.DbError
	var exist bool
	exist, dbError.MysqlError = UpdateMessageOfMysql(&(Message.MessageInfo))

	return exist, other.ErrorOutput(&dbError)
} //更新文章，返回是否发生错误

func (Message *MessageInfo) GetInfo() (bool, bool) {
	var dbError model.DbError
	var bo bool

	if Message.GetWay == "username" {
		bo, dbError.MysqlError = GetMessageOfMysqlByUsername(Message)
	}
	if Message.GetWay == "message_id" {
		bo, dbError.MysqlError = GetMessageOfMysqlByID(Message)
		if len(Message.Info) != 0 {
			Message.MessageInfo = Message.Info[0]
		}
	}
	if Message.GetWay == "message_name" {
		bo, dbError.MysqlError = GetMessageOfMysqlByMessageName(Message)
	}

	return bo, other.ErrorOutput(&dbError)
} //获取兵更改message信息，返回是否存在和是否发生错误
func (Message *MessageInfo) DelInfo() bool {
	var dbError model.DbError

	dbError.MysqlError = DeleteMessageOfMysql(&(Message.MessageInfo))

	return other.ErrorOutput(&dbError)
} //删除文章
