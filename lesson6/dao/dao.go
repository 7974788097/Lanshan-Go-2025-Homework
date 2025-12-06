package dao

import (
	"encoding/json"
	"fmt"
	"os"
)

// 模拟数据库
var Database = map[string]string{}

const userfile = "users.json"

func AddUser(username string, password string) error {
	Database[username] = password
	err := UpdateFile()
	if err != nil {
		return err
	}
	return nil
}

func FindUser(username string, password string) bool {
	if pwd, ok := Database[username]; ok {
		if pwd == password {
			return true
		}
	}
	return false
}
func SelectPasswordFromUsername(username string) string {
	return Database[username]
}
func UpdatePassword(username string, password string) error {
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
