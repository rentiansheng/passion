package redislock

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	cache "github.com/rentiansheng/passion/client/redis/client"
	"github.com/rentiansheng/passion/lib/bytesconv"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/5/19
    @desc:

***************************/

type Lock interface {
	Lock(ctx context.Context) (bool, error)
	Unlock(ctx context.Context) error
}

type lock struct {
	rid          string
	key          string
	lockedExpire time.Duration
}

func NewLock(key string, lockedExpire time.Duration) Lock {

	return &lock{
		rid:          uuid.NewString(),
		key:          key,
		lockedExpire: lockedExpire,
	}
}

// Lock 获取执行锁, lockedExpire 单位是Nanosecond， 最后回转换为 Millisecond
func (l *lock) Lock(ctx context.Context) (bool, error) {

	// 是否可以执行任务
	ok, err := cache.Redis().SetNX(ctx, l.key, bytesconv.StringToBytes(l.rid), l.lockedExpire/time.Millisecond).Result()
	if err != nil {
		return false, err
	}

	return ok, nil

}

// Unlock 释放锁
func (l *lock) Unlock(ctx context.Context) error {

	val, err := cache.Redis().Get(ctx, l.key).String()
	if err != nil {
		return err
	}
	// 判断是否是否当前任务锁的
	if val == l.rid {
		if err := cache.Redis().Del(ctx, l.key); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unauthorized operation")
	}

	return nil
}
