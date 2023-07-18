package dao

import (
	"context"
	"mall/model"

	"gorm.io/gorm"
)

type UserCouponDao struct {
	*gorm.DB
}

func NewUserCouponDao(ctx context.Context) *UserCouponDao {
	return &UserCouponDao{NewDBClient(ctx)}
}

func NewUserCouponDaoByDB(db *gorm.DB) *UserCouponDao {
	return &UserCouponDao{db}
}

// // GetUserById 根据 id 获取用户
// func (dao *UserCouponDao) GetUserById(uId uint) (usercoupon *model.UserCoupon, err error) {
// 	err = dao.DB.Model(&model.UserCoupon{}).Where("id=?", uId).
// 		First(&usercoupon).Error
// 	return
// }

// // UpdateUserById 根据 id 更新用户信息
// func (dao *UserDao) UpdateUserById(uId uint, user *model.User) error {
// 	return dao.DB.Model(&model.User{}).Where("id=?", uId).
// 		Updates(&user).Error
// }

// // ExistOrNotByUserName 根据username判断是否存在该名字
// func (dao *UserDao) ExistOrNotByUserName(userName string) (user *model.User, exist bool, err error) {
// 	var count int64
// 	err = dao.DB.Model(&model.User{}).Where("user_name=?", userName).
// 		Count(&count).Error
// 	if count == 0 {
// 		return nil, false, err
// 	}
// 	err = dao.DB.Model(&model.User{}).Where("user_name=?", userName).
// 		First(&user).Error
// 	if err != nil {
// 		return nil, false, err
// 	}
// 	return user, true, nil
// }

// 插入一条领取记录
func (dao *UserCouponDao) CreateUserCoupon(usercoupon *model.UserCoupon) (err error) {
	err = dao.DB.Model(&model.UserCoupon{}).Create(&usercoupon).Error
	return
}
