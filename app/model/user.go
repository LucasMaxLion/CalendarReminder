package model

import "fmt"

func GetUser(name string) *User {
	var ret User
	Conn.Table("user").Where("name =  ?", name).First(&ret)
	return &ret
}

func CreateUser(user *User) error {
	if err := Conn.Create(user).Error; err != nil {
		fmt.Printf("fatal form model createduser err :%s", err.Error())
		return err
	}
	return nil
}
