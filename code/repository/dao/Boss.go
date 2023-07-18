package dao

import (
	"context"
	"mall/model"

	"gorm.io/gorm"
)

type BossDao struct {
	*gorm.DB
}

func NewBossDao(ctx context.Context) *BossDao {
	return &BossDao{NewDBClient(ctx)}
}

func NewBossDaoByDB(db *gorm.DB) *BossDao {
	return &BossDao{db}
}

// GetUserById 根据 id 获取用户
func (dao *BossDao) GetBossById(bId uint) (boss *model.Boss, err error) {
	err = dao.DB.Model(&model.Boss{}).Where("id=?", bId).
		First(&boss).Error
	return
}

// UpdateUserById 根据 id 更新用户信息
func (dao *BossDao) UpdateBossById(bId uint, boss *model.Boss) error {
	return dao.DB.Model(&model.Boss{}).Where("id=?", bId).
		Updates(&boss).Error
}

// ExistOrNotByUserName 根据username判断是否存在该名字
func (dao *BossDao) ExistOrNotByBossName(bossName string) (boss *model.Boss, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.Boss{}).Where("boss_name=?", bossName).
		Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	err = dao.DB.Model(&model.Boss{}).Where("boss_name=?", bossName).
		First(&boss).Error
	if err != nil {
		return nil, false, err
	}
	return boss, true, nil
}

// CreateUser 创建用户
func (dao *BossDao) CreateBoss(boss *model.Boss) (err error) {
	err = dao.DB.Model(&model.Boss{}).Create(&boss).Error
	return
}
