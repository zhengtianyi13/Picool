package service

import (
	"context"
	"mall/conf"
	"mall/dao"
	"mall/model"
	"mall/pkg/e"
	util "mall/pkg/utils"
	"mall/serializer"
	"mime/multipart"
	"time"

	logging "github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

// UserService 管理用户服务
type BossService struct {
	NickName string `form:"nick_name" json:"nick_name"`
	BossName string `form:"boss_name" json:"boss_name"`
	Password string `form:"password" json:"password"`
	Key      string `form:"key" json:"key"` // 前端进行判断
}

type BossSendEmailService struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	//OpertionType 1:绑定邮箱 2：解绑邮箱 3：改密码
	OperationType uint `form:"operation_type" json:"operation_type"`
}

type BossValidEmailService struct {
}

func (service BossService) Register(ctx context.Context) serializer.Response { //不传指针说明只读方法

	code := e.SUCCESS //e中定义了返回码和返回内容

	if service.Key == "" || len(service.Key) != 16 {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密钥长度不足",
		}
	}
	//进行加密操作

	//设置密钥
	util.Encrypt.SetKey(service.Key)

	bossDao := dao.NewBossDao(ctx)
	_, exist, err := bossDao.ExistOrNotByBossName(service.BossName)

	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	boss := &model.Boss{
		NickName: service.NickName,
		BossName: service.BossName,
		Status:   model.Active,
		Money:    util.Encrypt.AesEncoding("10000"), // 初始金额,金额使用前面的key生成的加密码加密
	}
	//加密密码
	if err = boss.SetPassword(service.Password); err != nil {
		logging.Info(err)
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	boss.Avatar = "http://q1.qlogo.cn/g?b=qq&nk=294350394&s=640"
	//创建用户
	err = bossDao.CreateBoss(boss)

	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

//Login 用户登陆函数
func (service BossService) Login(ctx context.Context) serializer.Response {
	code := e.SUCCESS
	bossDao := dao.NewBossDao(ctx)
	boss, exist, err := bossDao.ExistOrNotByBossName(service.BossName)
	if !exist { //如果查询不到，返回相应的错误
		logging.Info(err)
		code = e.ErrorUserNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if boss.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	token, err := util.GenerateToken(boss.ID, service.BossName, 0)
	if err != nil {
		logging.Info(err)
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.TokenData{User: serializer.BuildBoss(boss), Token: token},
		Msg:    e.GetMsg(code),
	}
}

//Update 用户修改信息
func (service BossService) Update(ctx context.Context, uId uint) serializer.Response {
	var err error
	code := e.SUCCESS
	//找到用户
	bossDao := dao.NewBossDao(ctx)
	boss, err := bossDao.GetBossById(uId)
	if service.NickName != "" {
		boss.NickName = service.NickName
	}

	err = bossDao.UpdateBossById(uId, boss)
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
		Data:   serializer.BuildBoss(boss),
		Msg:    e.GetMsg(code),
	}
}

func (service *BossService) Post(ctx context.Context, uId uint, file multipart.File, fileSize int64) serializer.Response {
	code := e.SUCCESS
	var err error

	// path, err := UploadToQiNiu(file, fileSize)

	// if err != nil {
	// 	code = e.ErrorUploadFile
	// 	return serializer.Response{
	// 		Status: code,
	// 		Data:   e.GetMsg(code),
	// 		Error:  path,
	// 	}
	// }

	bossDao := dao.NewBossDao(ctx)        //拿到数据库返回对象，dao是跟数据库进行交互的，dao中是写sql的
	boss, err := bossDao.GetBossById(uId) //通过dao去获取用户信息，byid

	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//保存图片到本地
	path, err := UploadAvatarToLocalStatic(file, uId, boss.BossName)
	if err != nil {
		//
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  path,
		}

	}
	boss.Avatar = path                      //user对象中的Avatar就等于我们现在保存的地址
	err = bossDao.UpdateBossById(uId, boss) //通过userdao去更新数据库

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
		Data:   serializer.BuildBoss(boss),
		Msg:    e.GetMsg(code),
	}

}

// Send 发送邮件
func (service *BossSendEmailService) Send(ctx context.Context, id uint) serializer.Response {
	code := e.SUCCESS
	var address string

	token, err := util.GenerateEmailToken(id, service.OperationType, service.Email, service.Password) //生成密码token
	if err != nil {
		logging.Info(err)
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	noticeDao := dao.NewNoticeDao(ctx) //邮件模板，通过数据库获取
	notice, err := noticeDao.GetNoticeById(service.OperationType)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	address = conf.ValidEmail + token    //接口地址+token
	mailStr := notice.Text               //邮件模板，现在还没有。。。，你要自己去加
	mailText := mailStr + "\n" + address //邮件内容

	m := mail.NewMessage() //第三方包gopkg的mail

	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "youjian test")
	m.SetBody("text/html", mailText)

	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass) //

	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		logging.Info(err)
		code = e.ErrorSendEmail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}

// Valid 验证内容
func (service BossValidEmailService) Valid(ctx context.Context, token string) serializer.Response {
	var bossID uint
	var email string
	var password string
	var operationType uint
	code := e.SUCCESS

	//验证token
	if token == "" {
		code = e.InvalidParams
	} else {
		claims, err := util.ParseEmailToken(token)

		if err != nil {
			logging.Info(err)
			code = e.ErrorAuthCheckTokenFail

		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		} else {

			bossID = claims.UserID
			email = claims.Email
			password = claims.Password
			operationType = claims.OperationType

		}
	}
	if code != e.SUCCESS {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	//获取该用户信息
	bossDao := dao.NewBossDao(ctx)
	boss, err := bossDao.GetBossById(bossID)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if operationType == 1 {
		//1:绑定邮箱
		boss.Email = email
	} else if operationType == 2 {
		//2：解绑邮箱
		boss.Email = ""
	} else if operationType == 3 {
		//3：修改密码
		err = boss.SetPassword(password)
		if err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}
	err = bossDao.UpdateBossById(bossID, boss)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 成功则返回用户的信息
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildBoss(boss),
	}
}
