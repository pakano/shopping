package user

import (
	"shopping/utils/hash"

	"gorm.io/gorm"
)

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Salt == "" {
		// 为salt创建一个随机字符串
		salt := hash.CreateSalt()
		// 创建hash加密密码
		hashPassword, err := hash.HashPassword(u.Password + salt)
		if err != nil {
			return nil
		}
		u.Password = hashPassword
		u.Salt = salt
	}

	return
}
