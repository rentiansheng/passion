package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rentiansheng/mapper"
	cache "github.com/rentiansheng/passion/client/redis/client"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/11/9
    @desc:

***************************/

type FetchCacheDataFn func(ctx context.Context) (interface{}, error)
type AlarmFn func(ctx context.Context, err error)

type CacheLogic struct {
	// 获取需要缓存数据的方法
	fn FetchCacheDataFn
	// 缓存的key
	key string
	// 缓存的时间
	expire  time.Duration
	alarmFn AlarmFn
}

func NewCacheLogic(key string, expire time.Duration, fn FetchCacheDataFn) *CacheLogic {
	return &CacheLogic{
		fn:     fn,
		key:    key,
		expire: expire,
	}
}

func NewCacheLogicAlarm(key string, expire time.Duration, fn FetchCacheDataFn, alarmFn AlarmFn) *CacheLogic {
	return &CacheLogic{
		fn:      fn,
		key:     key,
		expire:  expire,
		alarmFn: alarmFn,
	}
}

func (c *CacheLogic) Get(ctx context.Context, result interface{}) error {
	noFound, err := c.fetchCache(ctx, result)
	if err == nil {
		if !noFound {
			return nil
		}
		// 从缓存获取数据失败或者缓存中数据不存在。降级到redis 执行
	}
	if err := c.fetchRaw(ctx, &result); err != nil {
		// 获取数据失败
		return err
	}
	if err := c.cache(ctx, result); err != nil {
		// 存储失败，不影响后续结果，缓存无法生效
		c.alarm(ctx, err)
	}
	return nil
}

func (c *CacheLogic) fetchCache(ctx context.Context, result interface{}) (bool, error) {

	byteVals, err := cache.Redis().Get(ctx, c.key).Bytes()
	if err != nil {
		if cache.Redis().IsNil(err) {
			return true, nil
		}
		return false, err
	}
	if err := json.Unmarshal(byteVals, result); err != nil {
		return false, err
	}
	return false, nil
}

func (c CacheLogic) alarm(ctx context.Context, err error) {
	if c.alarmFn != nil {
		c.alarmFn(ctx, err)
	}
}

func (c *CacheLogic) fetchRaw(ctx context.Context, result interface{}) error {
	data, err := c.fn(ctx)
	if err != nil {
		return err
	}
	if err := mapper.Mapper(ctx, data, result); err != nil {
		return err
	}
	return nil
}

func (c *CacheLogic) cache(ctx context.Context, data interface{}) error {

	cacheByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := cache.Redis().Set(ctx, c.key, cacheByte, c.expire).Err(); err != nil {
		return err
	}
	return nil
}
