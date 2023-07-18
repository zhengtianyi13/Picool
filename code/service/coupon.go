package service

import (
	"context"
	"database/sql"
	"errors"
	"mall/dao"
	"mall/model"
	"mall/pkg/e"
	util "mall/pkg/utils"
	"mall/serializer"
	"strconv"
	"sync"

	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//创建优惠券服务类，其属性是优惠券可以从前端传入的属性
type CouponService struct {
	ID            uint   `form:"id" json:"id"`
	Name          string `form:"name" json:"name"`
	Coupontype    int    `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" binding:"required,min=2,max=100"`
	Info          string `form:"info" json:"info" binding:"max=1000"`
	FinishTime    string `form:"finishtime" json:"finishtime"`
	CreateTime    string `form:"createTime" json:"createTime"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	Effect        bool   `form:"effect" json:"effect"`
	Num           int    `form:"num" json:"num"`
	model.BasePage
}

type ListCouponImgService struct {
}

// 商品
func (service *CouponService) Show(ctx context.Context, id string) serializer.Response {
	code := e.SUCCESS

	pId, _ := strconv.Atoi(id)

	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(uint(pId))
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data:   serializer.BuildProduct(product),
		Msg:    e.GetMsg(code),
	}
}

//创建优惠券
func (service *CouponService) Create(ctx context.Context, uId uint) serializer.Response {
	var boss *model.Boss //提交的商户类
	var err error
	code := e.SUCCESS

	bossDao := dao.NewBossDao(ctx)       //创建商户dao
	boss, err = bossDao.GetBossById(uId) //获取商户信息
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// path, err := UploadToQiNiu(tmp, files[0].Size)

	// if err != nil {
	// 	code = e.ErrorUploadFile
	// 	return serializer.Response{
	// 		Status: code,
	// 		Data:   e.GetMsg(code),
	// 		Error:  path,
	// 	}
	// }

	coupon := &model.Coupon{ //构建新的优惠券类
		Name:          service.Name, //优惠券名称
		Coupontype:    uint(service.Coupontype),
		Title:         service.Title,
		Info:          service.Info,
		FinishTime:    service.FinishTime, //优惠券结束时间
		CreateTime:    service.CreateTime, //优惠券创建时间
		DiscountPrice: service.DiscountPrice,
		Effect:        service.Effect,
		Num:           service.Num,
		BossID:        int(uId),
		BossName:      boss.BossName,
	}

	couponDao := dao.NewCouponDao(ctx)
	err = couponDao.CreateCoupon(coupon)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data:   serializer.BuildCoupon(coupon),
		Msg:    e.GetMsg(code),
	}
}

func (service *CouponService) List(ctx context.Context, uId uint) serializer.Response {
	var coupons []*model.Coupon
	var total int64
	code := e.SUCCESS

	bossDao := dao.NewBossDao(ctx)        //创建商户dao
	boss, err := bossDao.GetBossById(uId) //获取商户信息
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.PageSize == 0 { //如果没有传入每页的数量，则默认为15
		service.PageSize = 15
	}
	condition := make(map[string]interface{}) //创建一个map，用于存放查询条件
	condition["boss_id"] = int(boss.ID)       //查询条件为商户id
	// if service.Coupontype != 0 { //如果查询条件中有优惠券类型，则加入查询条件
	// 	condition["category_id"] = service.Coupontype
	// }
	couponDao := dao.NewCouponDao(ctx)
	total, errcoupon := couponDao.CountCouponByCondition(condition) //获取符合条件的优惠券总数
	if errcoupon != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	wg := new(sync.WaitGroup) //创建一个同步等待的组
	wg.Add(1)                 //添加一个等待
	go func() {               //开启一个协程
		couponDao = dao.NewCouponDaoByDB(couponDao.DB) //创建一个新的dao
		coupons, _ = couponDao.ListCouponByCondition(condition, service.BasePage)
		wg.Done()
	}()
	wg.Wait() //主线程阻塞，等待上面的协程查询完毕

	return serializer.BuildListResponse(serializer.BuildCoupons(coupons), uint(total))
}

//删除商品
func (service *CouponService) Delete(ctx context.Context, pId string) serializer.Response {
	code := e.SUCCESS

	couponDao := dao.NewCouponDao(ctx)
	couponId, _ := strconv.Atoi(pId)
	err := couponDao.DeleteCoupon(uint(couponId))
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// 更新优惠券
func (service *CouponService) Update(ctx context.Context, pId string) serializer.Response {
	code := e.SUCCESS
	couponDao := dao.NewCouponDao(ctx)

	couponId, _ := strconv.Atoi(pId)
	coupon := &model.Coupon{
		Name:          service.Name,
		Coupontype:    uint(service.Coupontype),
		Title:         service.Title,
		Info:          service.Info,
		FinishTime:    service.FinishTime,
		CreateTime:    service.CreateTime,
		DiscountPrice: service.DiscountPrice,
		Effect:        service.Effect,
	}
	err := couponDao.UpdateCoupon(uint(couponId), coupon)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// 搜索优惠券
func (service *CouponService) Search(ctx context.Context) serializer.Response {
	code := e.SUCCESS
	if service.PageSize == 0 {
		service.PageSize = 15
	}

	couponDao := dao.NewCouponDao(ctx)
	coupons, err := couponDao.SearchCoupon(service.Info, service.BasePage)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildCoupons(coupons), uint(len(coupons)))
}

//领取优惠券
func (service *CouponService) Receive(ctx context.Context, pId string, uId uint) serializer.Response {
	code := e.SUCCESS

	//输入string转换成uint
	pIduint := util.StringToUint(pId)

	//获取优惠券信息
	couponDao := dao.NewCouponDao(ctx)
	coupon, err := couponDao.GetCouponById(pIduint)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//生成优惠券code
	codeRandstring := util.RandomString(4)
	couponcode := strconv.FormatInt(int64(pIduint), 16) + strconv.FormatInt(int64(uId), 16) + codeRandstring
	//code 生成规则  前八位是16进制的优惠券id，后八位是16进制的用户id，后四位是随机生成的字符串

	//用户优惠券表需要插入的数据，分别是领的用户和领的什么券，券code和code类型，类型默认为1，表示没有生效
	usercoupon := &model.UserCoupon{
		UserID:     uId,
		CouponID:   pIduint,
		CouponCode: couponcode,
		Type:       1,
	}
	query := map[string]interface{}{"id": pId, "num": coupon.Num} //更新时的查询条件,这里添加了num即每次更新需要和上次查询的数据进行对比
	//这个就是乐观锁cas

	//开启事务
	connection := dao.NewDBClient(ctx)
	err = connection.Transaction(func(tx *gorm.DB) error {
		//要注意这里要用tx了，不能用上面的dao了
		//插入用户优惠券表
		err = tx.Model(&model.UserCoupon{}).Create(usercoupon).Error
		if err != nil {
			return err
		}
		//更新优惠券数量
		if coupon.Num > 1 {
			coupon.Num -= 1

			res := tx.Model(&model.Coupon{}).Where(query).
				Updates(coupon).RowsAffected
			if res == 0 {
				return errors.New("数据库没有更新成功可能是乐观锁冲突")
			}

		} else if coupon.Num == 1 { //为什么要这里加一个，因为gorm0值不会更新，所以要手动设置null
			couponzero := map[string]interface{}{
				"Name":          coupon.Name, //优惠券名称
				"Coupontype":    coupon.Coupontype,
				"Title":         coupon.Title,
				"Info":          coupon.Info,
				"FinishTime":    coupon.FinishTime, //优惠券结束时间
				"CreateTime":    coupon.CreateTime, //优惠券创建时间
				"DiscountPrice": coupon.DiscountPrice,
				"Effect":        coupon.Effect,
				"Num":           sql.NullInt32{},
				"BossID":        coupon.BossID,
				"BossName":      coupon.BossName,
			}
			res := tx.Model(&model.Coupon{}).Where(query).
				Updates(couponzero).RowsAffected
			if res == 0 {
				return errors.New("数据库没有更新成功可能是乐观锁冲突")
			}

		} else {
			err = errors.New("优惠券数量不足")
		}
		if err != nil {
			return err
		}

		return nil

	})
	//事物结束
	if err != nil {
		code = e.ErrorTransaction
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}
