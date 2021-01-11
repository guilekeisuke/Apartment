package model

import (
	"fmt"
	"time"
)

type User struct {
	ID          int       `json:"id" gorm:"praimaly_key"`
	UserName    string    `json:"username"`
	Password    string    `json:"password"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	RoleId      int       `json:"role_id" gorm:"default:1"`
	CreatedAt   time.Time `json:"create_at"`
	LastLoginAt time.Time `json:"last_login_at" gorm:"default:null"`
}

func CreateUser(user *User) {
	db.Create(user)
}

func FindUser(u *User) User {
	var user User
	db.Where(u).First(&user)
	fmt.Println(user)
	return user
}
