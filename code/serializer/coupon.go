package serializer

import (
	"mall/model"
)

type Coupon struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Coupontype    uint   `json:"category_id"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	FinishTime    string `json:"finishtime"`
	CreateTime    string `json:"createTime"`
	DiscountPrice string `json:"discount_price"`
	CreatedAt     int64  `json:"created_at"`
	Num           int    `json:"num"`
	BossID        int    `json:"boss_id"`
	BossName      string `json:"boss_name"`
}

// 序列化商品
func BuildCoupon(item *model.Coupon) Coupon { //由指针转化为实际的类
	return Coupon{
		ID:            item.ID,
		Name:          item.Name,
		Coupontype:    item.Coupontype,
		Title:         item.Title,
		Info:          item.Info,
		FinishTime:    item.FinishTime,
		CreateTime:    item.CreateTime,
		DiscountPrice: item.DiscountPrice,
		Num:           item.Num,
		CreatedAt:     item.CreatedAt.Unix(),
		BossID:        item.BossID,
		BossName:      item.BossName,
	}
}

//序列化商品列表
func BuildCoupons(items []*model.Coupon) (coupons []Coupon) {
	for _, item := range items {
		coupon := BuildCoupon(item)
		coupons = append(coupons, coupon)
	}
	return coupons
}
