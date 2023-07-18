package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//User 用户模型
type Boss struct {
	gorm.Model
	BossName       string `gorm:"unique"`
	Email          string
	PasswordDigest string
	NickName       string
	Status         string
	Avatar         string `gorm:"size:1000"`
	Money          string
}

//SetPassword 设置密码
func (Boss *Boss) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	Boss.PasswordDigest = string(bytes)
	return nil
}

//CheckPassword 校验密码
func (Boss *Boss) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(Boss.PasswordDigest), []byte(password))
	return err == nil
}

//AvatarUrl 封面地址
func (Boss *Boss) AvatarURL() string {
	signedGetURL := Boss.Avatar
	return signedGetURL
}
