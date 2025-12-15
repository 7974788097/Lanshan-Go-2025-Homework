package dao

import (
	"log"
	"moddle/model"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 模拟数据库
//var Database = map[string]string{}

var lock = sync.RWMutex{}
var db *gorm.DB

const Userdatabase = "root:123456@tcp(127.0.0.1:3306)/Users?charset=utf8mb4&parseTime=True&loc=Local"

func ConnectDatabase() {
	var err error
	db, err = gorm.Open(mysql.Open(Userdatabase), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("自动迁移成功")
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}

func AddUser(username string, password string) {
	lock.Lock()
	add := &model.User{Username: username, Password: password}
	defer lock.Unlock()
	db.Create(add)
	//err := UpdateFile()
	//if err != nil {
	//	return err
	//}
}

func FindUser(username string) bool {
	lock.RLock()
	defer lock.RUnlock()
	var user model.User
	db.Where("username = ?", username).First(&user)
	if user.Username == "" {
		return false
	}
	return true
	//if pwd, ok := Database[username]; ok {
	//	if pwd == password {
	//		return 0
	//	}
	//	return 2
	//}
	//return 1
}
func JudgePassword(username string, password string) bool {
	var user model.User
	db.Where("username = ?", username).First(&user)
	if user.Password == password {
		return true
	}
	return false
}
func SelectPasswordFromUsername(username string) string {
	lock.RLock()
	defer lock.RUnlock()
	var user model.User
	db.Where("username = ?", username).First(&user)
	return user.Password
	//return Database[username]
}
func UpdatePassword(username string, password string) {
	lock.Lock()
	defer lock.Unlock()
	var user model.User
	db.Where("username = ?", username).First(&user)
	db.Model(&user).Update("password", password)
	//Database[username] = password
	//err := UpdateFile()
	//if err != nil {
	//	return err
	//}
	//return nil
}

//func UpdateFile() error {
//	jsondata, err := json.MarshalIndent(Database, "", "\t")
//	if err != nil {
//		fmt.Printf("json marshal err: %v\n", err)
//		return err
//	}
//	err = os.WriteFile(userfile, jsondata, 0644)
//	if err != nil {
//		fmt.Printf("write file err: %v\n", err)
//		return err
//	}
//	return nil
//}
//func ReadDatabase() error {
//	if _, err := os.Stat(userfile); os.IsNotExist(err) {
//		file, err := os.Create(userfile)
//		if err != nil {
//			fmt.Printf("create file err: %v\n", err)
//			return err
//		}
//		_, err = file.WriteString("{}")
//		if err != nil {
//			fmt.Printf("write empty file err: %v\n", err)
//			return err
//		}
//		_ = file.Close()
//		fmt.Println("create file success")
//		return nil
//	}
//	file, err := os.ReadFile(userfile)
//	if err != nil {
//		fmt.Printf("open file err: %v\n", err)
//		return err
//	}
//	err = json.Unmarshal(file, &Database)
//	if err != nil {
//		fmt.Printf("json unmarshal err: %v\n", err)
//		return err
//	}
//	return nil
//}
