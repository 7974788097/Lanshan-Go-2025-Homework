package dao

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// 模拟数据库
var Database = map[string]string{}

var lock = sync.RWMutex{}

const userfile = "users.json"

func AddUser(username string, password string) error {
	lock.Lock()
	defer lock.Unlock()
	Database[username] = password
	err := UpdateFile()
	if err != nil {
		return err
	}
	return nil
}

func FindUser(username string, password string) uint8 {
	lock.RLock()
	defer lock.RUnlock()
	if pwd, ok := Database[username]; ok {
		if pwd == password {
			return 0
		}
		return 2
	}
	return 1
}
func SelectPasswordFromUsername(username string) string {
	lock.RLock()
	defer lock.RUnlock()
	return Database[username]
}
func UpdatePassword(username string, password string) error {
	lock.Lock()
	defer lock.Unlock()
	Database[username] = password
	err := UpdateFile()
	if err != nil {
		return err
	}
	return nil
}
func UpdateFile() error {
	jsondata, err := json.MarshalIndent(Database, "", "\t")
	if err != nil {
		fmt.Printf("json marshal err: %v\n", err)
		return err
	}
	err = os.WriteFile(userfile, jsondata, 0644)
	if err != nil {
		fmt.Printf("write file err: %v\n", err)
		return err
	}
	return nil
}
func ReadFile() error {
	if _, err := os.Stat(userfile); os.IsNotExist(err) {
		file, err := os.Create(userfile)
		if err != nil {
			fmt.Printf("create file err: %v\n", err)
			return err
		}
		_, err = file.WriteString("{}")
		if err != nil {
			fmt.Printf("write empty file err: %v\n", err)
			return err
		}
		_ = file.Close()
		fmt.Println("create file success")
		return nil
	}
	file, err := os.ReadFile(userfile)
	if err != nil {
		fmt.Printf("open file err: %v\n", err)
		return err
	}
	err = json.Unmarshal(file, &Database)
	if err != nil {
		fmt.Printf("json unmarshal err: %v\n", err)
		return err
	}
	return nil
}
