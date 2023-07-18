package model

import (
	"github.com/jinzhu/gorm"
)

//优惠券模型
type Coupon struct {
	gorm.Model
	Name          string `gorm:"size:255;index"` //优惠券名称
	Coupontype    uint   `gorm:"not null"`       //优惠券类型
	Title         string //优惠券标题
	Info          string `gorm:"size:1000"` //优惠券描述
	FinishTime    string //优惠券结束时间
	CreateTime    string //优惠券创建时间
	DiscountPrice string //优惠券折扣金额
	Effect        bool   `gorm:"default:false"` //是否生效
	Num           int    //优惠券数量
	BossID        int    //优惠券所属商家ID
	BossName      string //优惠券所属商家名称
}
