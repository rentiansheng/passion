package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/rentiansheng/passion/client/redis/define"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/11/9
    @desc:

***************************/

var redisLib = make(map[string]InitFn, 0)

// Redis key expire is PX
//   redis command 参考 	https://redis.io/commands/
type Redis interface {

	//********  list  ********//

	BLPop(ctx context.Context, timeout time.Duration, keys ...string) define.String
	BRPop(ctx context.Context, timeout time.Duration, key string) define.String
	BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) define.String
	// LIndex https://redis.io/commands/lindex/
	LIndex(ctx context.Context, key string, index int64) define.String
	// LInsert https://redis.io/commands/linsert/
	// LInsert(ctx context.Context, key string, value interface{}, isBefore bool) define.Int
	LLen(ctx context.Context, key string) define.Int
	// TODO: LPOP
	LPush(ctx context.Context, key string, values ...interface{}) define.Int
	// TODO: LPUSHX
	LRange(ctx context.Context, key string, start, stop int64) define.Strings
	LRem(ctx context.Context, key string, count int64, value interface{}) define.Int
	// TODO: LSET
	LTrim(ctx context.Context, key string, start, stop int64) define.Bool
	RPop(ctx context.Context, key string) define.Strings
	RPopLPush(ctx context.Context, source, destination string) define.String
	RPush(ctx context.Context, key string, values ...interface{}) define.Int
	// TODO: RPUSHX

	//********  hash  ********//

	HDel(ctx context.Context, key string, fields ...string) define.Int
	// TODO: HExists
	HGet(ctx context.Context, key, field string) define.String
	HGetAll(ctx context.Context, key string) define.MapStr
	HIncrBy(ctx context.Context, key, field string, incr int64) define.Int
	// TODO: HINCRBYFLOAT
	HKeys(ctx context.Context, key string) define.Strings
	// TODO: HLEN
	HMGet(ctx context.Context, key string, fields ...string) define.MapStr
	// TODO: HMSET
	HSet(ctx context.Context, key string, values ...interface{}) define.Int
	HSetNX(ctx context.Context, key string, value interface{}) define.Bool
	// TODO: HSTRLEN, HAVLS,HSCAN

	//********  sets  ********//

	SAdd(ctx context.Context, key string, members ...interface{}) define.Int
	//  TODO: SCARD,SDIFF,SDIFFSTORE,SINTER,SINTERSTOR,SISMEMERE,
	SMembers(ctx context.Context, key string) define.Any
	//  TODO: SMOVE,SPOP,SRANDMEMERR
	SRem(ctx context.Context, key string, members ...interface{}) define.Int
	// TODO: SUNION SUNIONSTORE,SSCAN

	//******** SORTED   sets  ********//

	//******** key  ********//

	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) define.Bool
	Expire(ctx context.Context, key string, expire time.Duration) define.Int
	Get(ctx context.Context, key string) define.Any
	MGet(ctx context.Context, keys ...string) define.Any
	MSet(ctx context.Context, values ...interface{}) define.Bool
	Rename(ctx context.Context, key, newKey string) define.Bool
	RenameNX(ctx context.Context, key, newKey string) define.Bool
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) define.Bool
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) define.Bool
	MSetNX(ctx context.Context, key string, values ...interface{}) define.Bool
	TTL(ctx context.Context, key string) define.Int
	Incr(ctx context.Context, key string) define.Int
	IncrBy(ctx context.Context, key string, incr int64) define.Int

	//******** client  ********//

	Ping(ctx context.Context) error
	Close() error
	IsNil(err error) bool
}

type InitFn func(cfg define.RedisConfig) (Redis, error)

func RegisterRedis(name string, fn InitFn) {
	if _, ok := redisLib[name]; ok {
		panic(fmt.Sprintf("redis lib name %s duplicate.", name))
	}
	redisLib[name] = fn
}

func GetRedisInit(name string) InitFn {

	return redisLib[name]
}

func Client() InitFn {
	return redisLib["default"]
}

func InitRedisIterator(cfg define.RedisConfig, fn func(name string, r Redis) error) error {
	for name, initFn := range redisLib {
		r, err := initFn(cfg)
		if err != nil {
			return err
		}
		if err := fn(name, r); err != nil {
			return err
		}
	}
	return nil
}
