package v1

import (
	util "mall/pkg/utils"
	"mall/service"

	"github.com/gin-gonic/gin"
)

func BossRegister(c *gin.Context) {
	var bossRegisterService service.BossService //相当于创建了一个UserRegisterService对象，调用这个对象中的Register方法。
	if err := c.ShouldBind(&bossRegisterService); err == nil {
		res := bossRegisterService.Register(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

//UserLogin 用户登陆接口
func BossLogin(c *gin.Context) {
	var bossLoginService service.BossService
	if err := c.ShouldBind(&bossLoginService); err == nil {
		res := bossLoginService.Login(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

func BossUpdate(c *gin.Context) {
	var bossUpdateService service.BossService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&bossUpdateService); err == nil {
		res := bossUpdateService.Update(c.Request.Context(), claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
