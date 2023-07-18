package util
import (
	"math/rand"
	"time"
	"gin-mall/cache"

	
)

func GenerateCouponCode() string {
	//生成时间戳
	timestamp := time.Now().Unix()
	//生成redis序列号

	//获取当天的日期精确到天
	today := time.Now().Format("2006-01-02")
	// 使用go-redis的incr命令，每次自增1，返回自增后的值
	redisSerial, err := cache.RedisClient.Incr(today).Result()
	if err != nil {
		logging.Info(err)
		
	}

	redisSerial := rand.Intn(1000000)