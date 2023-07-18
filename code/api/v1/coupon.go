package v1

import (
	util "mall/pkg/utils"
	"mall/service"

	"github.com/gin-gonic/gin"
)

// 创建优惠券
func CreateCoupon(c *gin.Context) {
	claim, _ := util.ParseToken(c.GetHeader("Authorization")) //获取商家的token判断是否合法
	createCouponService := service.CouponService{}            //实列化商品服务这个类
	//c.SaveUploadedFile()
	if err := c.ShouldBind(&createCouponService); err == nil { //类绑定
		res := createCouponService.Create(c.Request.Context(), claim.ID) //绑定完之后需要传入上下文和商家id
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

//商家发布的优惠券列表
func ListCoupons(c *gin.Context) {
	claim, _ := util.ParseToken(c.GetHeader("Authorization")) //获取商家的token判断是否合法
	listCouponsService := service.CouponService{}
	if err := c.ShouldBind(&listCouponsService); err == nil {
		res := listCouponsService.List(c.Request.Context(), claim.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

//商品详情
func ShowCoupons(c *gin.Context) {
	showCouponService := service.CouponService{}
	res := showCouponService.Show(c.Request.Context(), c.Param("id"))
	c.JSON(200, res)
}

//删除商品
func DeleteCoupon(c *gin.Context) {
	deleteCouponService := service.CouponService{}
	res := deleteCouponService.Delete(c.Request.Context(), c.Param("id"))
	c.JSON(200, res)
}

//更新商品
func UpdateCoupon(c *gin.Context) {
	updateCouponService := service.CouponService{}
	if err := c.ShouldBind(&updateCouponService); err == nil {
		res := updateCouponService.Update(c.Request.Context(), c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

//搜索商品
func SearchCoupons(c *gin.Context) {
	searchCouponsService := service.CouponService{}
	if err := c.ShouldBind(&searchCouponsService); err == nil {
		res := searchCouponsService.Search(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

//领取优惠券
func ReceiveCoupon(c *gin.Context) {
	claim, _ := util.ParseToken(c.GetHeader("Authorization")) //获取用户的token判断是否合法
	receiveCouponService := service.CouponService{}
	res := receiveCouponService.Receive(c.Request.Context(), c.Param("id"), claim.ID)
	if res.Status == 200 {
		c.JSON(200, res)
	} else {
		c.JSON(400, res)
	}
}
