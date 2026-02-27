package commentmanage

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
)

var dbContext *model.DBcontext

func ConnectCommentDB(info *model.DBInfo) {
	dbContext = new(model.DBcontext)
	connectCommentMysql(info)
}

type CommentInfo struct {
	model.CommentInfo
	GetWay string
	Info   []model.CommentInfo
}

func Make(info *map[string]interface{}) *CommentInfo {
	newInfo := &CommentInfo{}
	if (*info)["comment_id"] != nil {
		newInfo.CommentID = int64((*info)["comment_id"].(float64))
	}
	if (*info)["comment_content"] != nil {
		newInfo.Content = (*info)["comment_content"].(string)
	}
	if (*info)["username"] != nil {
		newInfo.Username = (*info)["username"].(string)
	}
	if (*info)["parent_node_id"] != nil {
		newInfo.ParentNodeID = int64((*info)["parent_node_id"].(float64))
	}
	if (*info)["get_way"] != nil {
		newInfo.GetWay = (*info)["get_way"].(string)
	}
	if (*info)["comment_path"] != nil {
		newInfo.Path = (*info)["comment_path"].(string)
	}
	return newInfo
}
func (Comment *CommentInfo) AddInfo() bool {
	var dbError model.DbError

	dbError.MysqlError = SetCommentOfMysql(&(Comment.CommentInfo))

	return other.ErrorOutput(dbError)
} //添加评论，返回是否发生错误

func (Comment *CommentInfo) ChangeInfo() (bool, bool) {
	var dbError model.DbError
	var exist bool
	exist, dbError.MysqlError = UpdateCommentOfMysql(&(Comment.CommentInfo))

	return exist, other.ErrorOutput(dbError)
} //更新评论，返回是否发生错误

func (Comment *CommentInfo) GetInfo() (bool, bool) {
	var dbError model.DbError
	var bo bool
	if Comment.GetWay == "" {
		Comment.GetWay = "comment_id"
	}
	if Comment.GetWay == "path" {
		bo, dbError.MysqlError = GetCommentOfMysqlByPath(Comment)
	}
	if Comment.GetWay == "comment_id" {
		bo, dbError.MysqlError = GetCommentOfMysqlByID(Comment)
		if len(Comment.Info) != 0 {
			Comment.CommentInfo = Comment.Info[0]
		}
	}
	if Comment.GetWay == "username" {
		bo, dbError.MysqlError = GetCommentOfMysqlByUsername(Comment)
	}

	return bo, other.ErrorOutput(dbError)
} //获取comment信息(id查询会更改)，返回是否存在和是否发生错误

func (Comment *CommentInfo) DelInfo() bool {
	var dbError model.DbError

	dbError.MysqlError = DeleteCommentOfMysqlByPath(&(Comment.CommentInfo))

	return other.ErrorOutput(dbError)
}
