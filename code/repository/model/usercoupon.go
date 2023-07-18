package model

import (
	"github.com/jinzhu/gorm"
)

//用户优惠券表
type UserCoupon struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	CouponID   uint   `gorm:"not null"`
	CouponCode string `gorm:"not null"`
	Type       uint   // 1 生效  2 未生效  3 已使用
}
