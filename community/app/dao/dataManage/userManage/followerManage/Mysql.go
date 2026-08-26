package followermanage

import (
	"gorm.io/gorm"
)

func (U *FollowerStruct) addInMysql(dbContext *gorm.DB) error {
	result := dbContext.Create(U)
	return result.Error
}
