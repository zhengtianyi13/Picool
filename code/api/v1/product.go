package v1

import (
	util "mall/pkg/utils"
	"mall/service"

	"github.com/gin-gonic/gin"
)

// 创建商品
func CreateProduct(c *gin.Context) {
	form, _ := c.MultipartForm()                              //多文件上传的获取
	files := form.File["file"]                                //获取其中的文件
	claim, _ := util.ParseToken(c.GetHeader("Authorization")) //获取用户的token判断是否合法
	createProductService := service.ProductService{}          //实列化商品服务这个类
	//c.SaveUploadedFile()
	if err := c.ShouldBind(&createProductService); err == nil { //类绑定
		res := createProductService.Create(c.Request.Context(), claim.ID, files) //调用类的方法
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

//商品列表
func ListProducts(c *gin.Context) {
	listProductsService := service.ProductService{}
	if err := c.ShouldBind(&listProductsService); err == nil {
		res := listProductsService.List(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

//商品详情
func ShowProduct(c *gin.Context) {
	showProductService := service.ProductService{}
	res := showProductService.Show(c.Request.Context(), c.Param("id"))
	c.JSON(200, res)
}

//删除商品
func DeleteProduct(c *gin.Context) {
	deleteProductService := service.ProductService{}
	res := deleteProductService.Delete(c.Request.Context(), c.Param("id"))
	c.JSON(200, res)
}

//更新商品
func UpdateProduct(c *gin.Context) {
	updateProductService := service.ProductService{}
	if err := c.ShouldBind(&updateProductService); err == nil {
		res := updateProductService.Update(c.Request.Context(), c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

//搜索商品
func SearchProducts(c *gin.Context) {
	searchProductsService := service.ProductService{}
	if err := c.ShouldBind(&searchProductsService); err == nil {
		res := searchProductsService.Search(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

func ListProductImg(c *gin.Context) {
	var listProductImgService service.ListProductImgService
	if err := c.ShouldBind(&listProductImgService); err == nil {
		res := listProductImgService.List(c.Request.Context(), c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
