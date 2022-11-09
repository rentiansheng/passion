package client

import (
	"fmt"

	"github.com/rentiansheng/passion/client/redis"
	"github.com/rentiansheng/passion/client/redis/define"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/11/9
    @desc:

***************************/

const defaultRedisName = "default"

var (
	instance redis.Redis
)

func InitRedis(cfg define.RedisConfig) error {

	// cfg 是就是默认值 default
	redisInitName := defaultRedisName
	fnInit := redis.GetRedisInit(redisInitName)
	if fnInit == nil {
		return fmt.Errorf("redis %s not implement", redisInitName)
	}
	r, err := fnInit(cfg)
	if err != nil {
		return err
	}
	instance = r

	return nil
}

func Redis() redis.Redis {
	return instance
}
