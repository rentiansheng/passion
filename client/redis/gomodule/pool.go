package gomodule

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rentiansheng/passion/client/redis/define"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/11/9
    @desc:

***************************/

func initPool(cfg define.RedisConfig) *redis.Pool {
	maxIdle := cfg.MaxIdleConn
	maxActive := cfg.MaxActiveConn
	timeout := time.Duration(cfg.Timeout)
	// 建立连接池
	redisClient := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(cfg.MaxIdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", cfg.DialAddress,
				redis.DialPassword(cfg.Password),
				redis.DialDatabase(cfg.DBIndex),
				redis.DialConnectTimeout(timeout*time.Second),
				redis.DialReadTimeout(timeout*time.Second),
				redis.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
	return redisClient
}
