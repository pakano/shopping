package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName  string `gorm:"type:varchar(30)"`
	Password  string `gorm:"type:varchar(100)"`
	Password2 string `gorm:"-"`
	Salt      string `gorm:"type:varchar(100)"`
	Token     string `gorm:"type:varchar(500)"`
	IsDeleted bool
	IsAdmin   bool
}

func NewUser(username, password, password2 string) *User {
	return &User{
		UserName:  username,
		Password:  password,
		Password2: password2,
	}
}
