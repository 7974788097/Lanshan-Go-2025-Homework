package main

import (
	"Project/modle"
	"fmt"
)

func main() {
	fmt.Println("输入账号，密码")
	A := manage.A{}
	B := manage.B{}
	manage.Adduser(B, 1)
	//B.Adduser(1)
	manage.Adduser(A, 114514)    //输入 一号 1
	manage.Adduser(A, 247915639) //输入 二号 2
	manage.Test()
	fmt.Println("...............................") //这只是一个可有可无的，毫无存在感的分割线
	manage.Deleteuser(A, 114514)
	manage.Test()
	fmt.Println("...............................") //这只是一个可有可无的，毫无存在感的分割线
	manage.Updatepassword(A, 247915639)            //输入22
	manage.Test()
	fmt.Println("...............................") //这只是一个可有可无的，毫无存在感的分割线
	manage.Updatepassword(A, 114514)               //输入11
	manage.Test()
}
