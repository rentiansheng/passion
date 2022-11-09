package gomodule

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	redisGo "github.com/gomodule/redigo/redis"

	"github.com/rentiansheng/passion/client/redis"
	"github.com/rentiansheng/passion/client/redis/define"
	"github.com/rentiansheng/passion/lib/array"
)

/***
Desc:
	expire: 时间必须用 int64，要不然回提示ERR value is not an integer or out of range
****/
type instance struct {
	prefix string
	client *redisGo.Pool
}

func init() {
	redis.RegisterRedis("default", initInstance)
}

func initInstance(cfg define.RedisConfig) (redis.Redis, error) {

	pool := initPool(cfg)
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	return &instance{prefix: cfg.KeyPrefix, client: pool}, nil

}

//********  list  ********//

// BLPop doc: https://redis.io/commands/blpop/
func (i instance) BLPop(ctx context.Context, timeout time.Duration, keys ...string) define.String {
	rConn := i.client.Get()
	defer rConn.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, key)
	}

	args = append(args, int64(timeout))

	return replyAny(rConn.Do("BLPOP", args...)).ToValString()
}

// BRPop https://redis.io/commands/brpop/
func (i instance) BRPop(ctx context.Context, timeout time.Duration, key string) define.String {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("BRPOP", key, int64(timeout))).ToValString()
}

//BRPopLPush doc: https://redis.io/commands/brpoplpush/
// deprecated:  Redis version 6.2.0, this command is regarded as deprecated.
func (i instance) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) define.String {
	rConn := i.client.Get()
	defer rConn.Close()
	return replyAny(rConn.Do("BRPOPLPUSH", source, destination, int64(timeout))).ToValString()
}

//LIndex https://redis.io/commands/lindex/
func (i instance) LIndex(ctx context.Context, key string, index int64) define.String {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("LINDEX", key, index)).ToString()
}

// LLen https://redis.io/commands/llen/
func (i instance) LLen(ctx context.Context, key string) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("LLEN", key)).ToInt()
}

func (i instance) LPush(ctx context.Context, key string, values ...interface{}) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	// 第一个参数是key, 这里使用convToString 因为gomodule 对key 会做类型解析
	args := mergeStrArgs(key, convToString(values...))
	return replyAny(rConn.Do("LPUSH", args...)).ToInt()
}

func (i instance) LRange(ctx context.Context, key string, start, stop int64) define.Strings {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("LRANGE", key, start, stop)).ToStrings()
}

func (i instance) LRem(ctx context.Context, key string, count int64, value interface{}) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("LREM", key, count, value)).ToInt()
}

func (i instance) LTrim(ctx context.Context, key string, start, stop int64) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("LTRIM", key, start, stop)).ToBool()
}

// RPop doc: https://redis.io/commands/rpop/
//     reply:  Bulk string reply, Array reply
func (i instance) RPop(ctx context.Context, key string) define.Strings {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("RPOP", key)).ToValStrings()
}

// RPopLPush dock: https://redis.io/commands/rpoplpush/
//  reply: Bulk string reply
func (i instance) RPopLPush(ctx context.Context, source, destination string) define.String {
	rConn := i.client.Get()
	defer rConn.Close()
	return replyAny(rConn.Do("RPOPLPUSH", source, destination)).ToString()
}

// RPush doc: https://redis.io/commands/rpush/
func (i instance) RPush(ctx context.Context, key string, values ...interface{}) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	// 第一个参数是key, 这里使用convToString 因为gomodule 对key 会做类型解析
	args := mergeStrArgs(key, convToString(values...))
	return replyAny(rConn.Do("RPUSH", args...)).ToInt()

}

//********  hash  ********//

// Hdel https://redis.io/commands/hdel/
func (i instance) HDel(ctx context.Context, key string, fields ...string) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	args := array.Strings.ToInterface(append([]string{key}, fields...))
	return replyAny(rConn.Do("HDEL", args...)).ToInt()
}

// HGet https://redis.io/commands/hget/
func (i instance) HGet(ctx context.Context, key, field string) define.String {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("HGET", key, field)).ToString()
}

// HGetAll https://redis.io/commands/hgetall/
func (i instance) HGetAll(ctx context.Context, key string) define.MapStr {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("HGETALL", key)).ToMapStr()
}

func (i instance) HIncrBy(ctx context.Context, key, field string, incr int64) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("HINCRBY", key, field, incr)).ToInt()

}

// HKeys https://redis.io/commands/hkeys/
func (i instance) HKeys(ctx context.Context, key string) define.Strings {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("HKEYS", key)).ToStrings()
}

func (i instance) HMGet(ctx context.Context, key string, fields ...string) define.MapStr {
	rConn := i.client.Get()
	defer rConn.Close()
	args := []interface{}{key}
	for _, field := range fields {
		args = append(args, field)
	}

	return replyAny(rConn.Do("HMGET", args...)).ToMapStr()
}

func (i instance) HSet(ctx context.Context, key string, values ...interface{}) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	args := append([]interface{}{key}, values...)
	return replyAny(rConn.Do("HSET", args...)).ToInt()
}

func (i instance) HSetNX(ctx context.Context, key string, value interface{}) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("HSETNX", key, value)).ToBool()
}

//********  sets  ********//

// SAdd doc: https://redis.io/commands/sadd/
func (i instance) SAdd(ctx context.Context, key string, members ...interface{}) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("SADD", mergeArgs(key, members)...)).ToInt()

}

func (i instance) SMembers(ctx context.Context, key string) define.Any {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("SREM", key))
}

func (i instance) SRem(ctx context.Context, key string, members ...interface{}) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("SREM", mergeArgs(key, members)...)).ToInt()
}

//******** SORTED   sets  ********//

//******** key  ********//

// Del
func (i instance) Del(ctx context.Context, key string) error {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("DEL", key)).Err()
}

// Existsc
func (i instance) Exists(ctx context.Context, key string) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("exists", key)).ToBool()
}

// Expire doc: https://redis.io/commands/expire/
func (i instance) Expire(ctx context.Context, key string, expire time.Duration) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("expire", key, int64(expire))).ToInt()
}

func (i instance) Get(ctx context.Context, key string) define.Any {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("GET", key))
}

func (i instance) MGet(ctx context.Context, keys ...string) define.Any {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("mget", keys))
}

func (i instance) MSet(ctx context.Context, values ...interface{}) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("MSET", values...)).ToBool()
}

func (i instance) Rename(ctx context.Context, key, newKey string) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("RENAME", key, newKey)).ToBool()
}

func (i instance) RenameNX(ctx context.Context, key, newKey string) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("RENAME", key, newKey)).ToBool()
}

// Set https://redis.io/commands/set/
//   reply:
//     	Simple string reply: OK if SET was executed correctly.
//		Null reply: (nil) if the SET operation was not performed because the user specified the NX or XX option but the condition was not met.
//		If the command is issued with the GET option, the above does not apply. It will instead reply as follows, regardless if the SET was actually performed:
//		Bulk string reply: the old string value stored at key.
//		Null reply: (nil) if the key did not exist.
func (i instance) Set(ctx context.Context, key string, value interface{}, expire time.Duration) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("SET", key, value, "PX", int64(expire))).ToBool()
}

func (i instance) SetNX(ctx context.Context, key string, value interface{}, expire time.Duration) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("SET", key, value, "PX", int64(expire), "nx")).ToBool()
}

func (i instance) MSetNX(ctx context.Context, key string, values ...interface{}) define.Bool {
	rConn := i.client.Get()
	defer rConn.Close()
	return replyAny(rConn.Do("MSETNX", mergeArgs(key, values)...)).ToBool()
}

func (i instance) TTL(ctx context.Context, key string) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("TTL", key)).ToInt()
}

// Incr doc: https://redis.io/commands/incr/
func (i instance) Incr(ctx context.Context, key string) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("INCR", key)).ToInt()
}

func (i instance) IncrBy(ctx context.Context, key string, incr int64) define.Int {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("INCRBY", key, incr)).ToInt()
}

//******** client  ********//

func (i instance) Ping(ctx context.Context) error {
	rConn := i.client.Get()
	defer rConn.Close()

	return replyAny(rConn.Do("PING")).Err()
}

func (i instance) Close() error {
	return nil
}

func (i instance) IsNil(err error) bool {
	return IsNil(err)
}

func IsNil(err error) bool {
	return errors.Is(err, redisGo.ErrNil)
}

func (i instance) getKey(key string) string {
	if !strings.HasPrefix(key, i.prefix) {
		return fmt.Sprintf("%s%s", i.prefix, key)
	}
	return key
}

var _ redis.Redis = (*instance)(nil)
