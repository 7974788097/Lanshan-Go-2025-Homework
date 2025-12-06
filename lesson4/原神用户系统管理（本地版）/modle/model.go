package manage

import "fmt"

var Users = make(map[int]User)

type User struct {
	name     string
	password string
}

type Manage interface {
	adduser(int)
	deleteuser(int)
	updatepassword(int)
}

func Adduser(m Manage, uid int) {
	m.adduser(uid)
}
func Deleteuser(m Manage, uid int) {
	m.deleteuser(uid)
}
func Updatepassword(m Manage, uid int) {
	m.updatepassword(uid)
}

type A struct{}

func (u A) adduser(uid int) {
	var name string
	var password string
	_, _ = fmt.Scan(&name)
	_, _ = fmt.Scan(&password)
	nu := User{name: name, password: password}
	Users[uid] = nu
}
func (u A) deleteuser(uid int) {
	delete(Users, uid)
}
func (u A) updatepassword(uid int) {
	_, exists := Users[uid]
	if !exists {
		fmt.Println("该用户不存在")
		return
	}
	var newpassword string
	_, _ = fmt.Scan(&newpassword)
	un := Users[uid]
	un.password = newpassword
	Users[uid] = un
}

type B struct{}

func (u B) adduser(uid int) {
	fmt.Println("已成功调用")
}
func (u B) deleteuser(uid int)     {}
func (u B) updatepassword(uid int) {}

func Test() {
	fmt.Println(Users)
}
