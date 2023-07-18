package dao

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var (
	_db *gorm.DB //gorm的DB对象  全局变量
)

func Database(connRead, connWrite string) { //连接配置数据库

	var ormLogger logger.Interface //orm日志的打印

	if gin.Mode() == "debug" { //日志处理
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connRead, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据版本自动配置
	}), &gorm.Config{ //gorm的配置
		Logger: ormLogger, //日志
		NamingStrategy: schema.NamingStrategy{ //表名的命名策略
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
		},
	})

	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)                  //设置连接池，空闲
	sqlDB.SetMaxOpenConns(100)                 //设置连接池，最大
	sqlDB.SetConnMaxLifetime(time.Second * 30) //设置连接池，最大连接时间
	_db = db

	//主从配置
	_ = _db.Use(dbresolver. //使用dbresolver，gorm分装的一个插件
				Register(dbresolver.Config{ //配置
			// `db2` 作为 sources，`db3`、`db4` 作为 replicas
			Sources:  []gorm.Dialector{mysql.Open(connWrite)},                      // 写操作
			Replicas: []gorm.Dialector{mysql.Open(connRead), mysql.Open(connRead)}, // 读操作
			Policy:   dbresolver.RandomPolicy{},                                    // sources/replicas 负载均衡策略
		}))
	Migration()  //这个会自动的按照你设置的类型进行迁移，即创建对应的表（migration中指定）
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
