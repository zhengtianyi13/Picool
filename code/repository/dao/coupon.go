package dao

import (
	"context"
	"mall/model"

	"gorm.io/gorm"
)

type CouponDao struct {
	*gorm.DB
}

func NewCouponDao(ctx context.Context) *CouponDao {
	return &CouponDao{NewDBClient(ctx)}
}

func NewCouponDaoByDB(db *gorm.DB) *CouponDao {
	return &CouponDao{db}
}

// GetCouponById 通过 id 获取Coupon
func (dao *CouponDao) GetCouponById(id uint) (coupon *model.Coupon, err error) {
	err = dao.DB.Model(&model.Coupon{}).Where("id=?", id).
		First(&coupon).Error
	return
}

// ListCouponByCondition 获取优惠券列表
func (dao *CouponDao) ListCouponByCondition(condition map[string]interface{}, page model.BasePage) (coupons []*model.Coupon, err error) {
	err = dao.DB.Where(condition).
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).Find(&coupons).Error
	return
}

// CreateProduct 创建商品
func (dao *CouponDao) CreateCoupon(coupon *model.Coupon) (err error) {
	err = dao.DB.Model(&model.Coupon{}).Create(&coupon).Error
	return
}

// CountCouponByCondition 根据情况获取商品的数量
func (dao *CouponDao) CountCouponByCondition(condition map[string]interface{}) (total int64, err error) {
	err = dao.DB.Model(&model.Coupon{}).Where(condition).Count(&total).Error
	return
}

// DeleteCoupon 删除商品
func (dao *CouponDao) DeleteCoupon(pId uint) (err error) {
	err = dao.DB.Model(&model.Coupon{}).Delete(&model.Coupon{}).Error
	return
}

// UpdateCoupon 更新商品
func (dao *CouponDao) UpdateCoupon(pId uint, coupon *model.Coupon) (err error) {
	err = dao.DB.Model(&model.Coupon{}).Where("id=?", pId).
		Updates(&coupon).Error
	return
}

// SearchCoupon 搜索商品
func (dao *CouponDao) SearchCoupon(info string, page model.BasePage) (coupons []*model.Coupon, err error) {
	err = dao.DB.Model(&model.Coupon{}).
		Where("name LIKE ? OR info LIKE ?", "%"+info+"%", "%"+info+"%").
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).Find(&coupons).Error
	return
}
