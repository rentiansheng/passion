package gomodule

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/require"

	"github.com/rentiansheng/passion/client/redis"
	"github.com/rentiansheng/passion/client/redis/define"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/9/6
    @desc:

***************************/

var (
	r              redis.Redis
	ctx            = context.Background()
	expire1Second  = time.Duration(1000)
	expire2Second  = expire1Second * 2
	expire5Second  = expire1Second * 5
	expire10Second = expire1Second * 10
	mockRedis      *miniredis.Miniredis
)

func init() {

	// Notice: miniredis expire not work. but using a physical Redis server is ok
	//    reason: https://github.com/alicebob/miniredis#ttls-key-expiration-and-time

	var err error
	mockRedis, err = miniredis.Run()
	if err != nil {
		panic(err)
	}

	cfg := define.RedisConfig{
		MaxIdleConn:   10,
		MaxActiveConn: 10,
		//DialAddress:   "127.0.0.1:6379",
		DialAddress: mockRedis.Addr(),
		DBIndex:     0,
		Password:    "",
		ClusterMode: false,
	}
	rr, err := initInstance(cfg)
	if err != nil {
		panic(err)
	}
	r = rr

}

func fnName() string {
	counter, _, _, success := runtime.Caller(1)

	if !success {
		println("functionName: runtime.Caller: failed")
		os.Exit(1)
	}
	names := strings.Split(runtime.FuncForPC(counter).Name(), "/")
	names = strings.Split(names[len(names)-1], ".")

	return names[len(names)-1]
}

func testBlockPop(t *testing.T, key string, inputs []interface{}, outputs []string, testFn func(key string) define.String) {
	suit := struct {
		inputs  []interface{}
		outputs []string
	}{
		inputs:  inputs,
		outputs: outputs,
	}

	if err := r.LPush(ctx, key, suit.inputs...).Err(); err != nil {
		t.Errorf("push error. %s", err.Error())
		return
	}
	for idx := 0; idx < len(suit.outputs); idx++ {
		val, err := testFn(key).Result()
		require.NoError(t, err, key+" execute error")
		require.Equalf(t, suit.outputs[idx], val, key+" error")
	}

	blockTestChn := make(chan struct{}, 1)
	ticker := time.NewTicker(time.Millisecond * 10)
	go func() {
		_, err := testFn(key).Result()
		if !r.IsNil(err) {
			require.NoError(t, err, key+" block execute error")
		} else {
			t.Errorf("block test failure. fnc name: %s", key)
		}
		close(blockTestChn)
	}()
	// test block
	select {
	case <-blockTestChn:
		t.Errorf("block test failure, fnc name: %s", key)
	case <-ticker.C:
	}
}
func TestInstanceBLPop(t *testing.T) {
	key := fnName()
	inputs := []interface{}{1, 2, 3, "str", true, []string{"1", "2"}}
	outputs := []string{`["1","2"]`, "true", "str", "3", "2", "1"}
	fn := func(key string) define.String {
		return r.BLPop(ctx, 1, key)
	}
	testBlockPop(t, key, inputs, outputs, fn)
}

func TestInstanceBRPop(t *testing.T) {
	key := fnName()
	inputs := []interface{}{1, 2, 3, "str", true, []string{"1", "2"}}
	outputs := []string{"1", "2", "3", "str", "true", `["1","2"]`}
	fn := func(key string) define.String {
		return r.BRPop(ctx, 1, key)
	}
	testBlockPop(t, key, inputs, outputs, fn)
}

func Test_instance_BRPopLPush(t *testing.T) {
	type args struct {
		ctx         context.Context
		source      string
		destination string
		timeout     time.Duration
	}
	tests := []struct {
		name string
		args args
		want define.String
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.BRPopLPush(tt.args.ctx, tt.args.source, tt.args.destination, tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BRPopLPush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceClose(t *testing.T) {
	if err := r.Close(); err != nil {
		t.Errorf("close error. %s", err.Error())
	}
}

func TestInstanceExistsAndDel(t *testing.T) {
	key := fnName()

	ok, err := r.Exists(ctx, key).Result()
	if err != nil {
		t.Errorf("exists %s", err.Error())
		return
	}
	if ok {
		t.Errorf("key exists")
		return
	}
	if err := r.Set(ctx, key, "", expire10Second).Err(); err != nil {
		t.Errorf("set error. %s", err.Error())
		return
	}
	ok, err = r.Exists(ctx, key).Result()
	if err != nil {
		t.Errorf("exists %s", err.Error())
		return
	}
	if !ok {
		t.Errorf("key not exists")
		return
	}

	if err := r.Del(ctx, key); err != nil {
		t.Errorf("delete error. %s", err.Error())
		return
	}
	ok, err = r.Exists(ctx, key).Result()
	if err != nil {
		t.Errorf("exists %s", err.Error())
		return
	}
	if ok {
		t.Errorf("key exists")
		return
	}

}

func TestInstanceExpireAndTTL(t *testing.T) {
	key := fnName()
	if err := r.Set(ctx, key, key, expire1Second).Err(); err != nil {
		t.Errorf("set key  error. %s", err.Error())
		return
	}
	if err := r.Expire(ctx, key, expire10Second).Err(); err != nil {
		t.Errorf("expire error. %s", err.Error())
		return
	}
	ttl, err := r.TTL(ctx, key).Result()
	if err != nil {
		t.Errorf("ttl error. %s", err.Error())
		return
	}
	if ttl < int64(expire1Second)-2 {
		t.Errorf("expire key ttl error")
		return
	}

	if err := r.Del(ctx, key); err != nil {
		t.Errorf("delete key error. %s", err.Error())
		return
	}

}

func TestInstanceGetAndSet(t *testing.T) {

	key := fnName()
	strVal := "TestInstaceGetAndSet"
	intVal := int64(1000)
	// 写入字符串
	bl, err := r.Set(ctx, key, strVal, expire1Second).Result()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !bl {
		t.Errorf("write string to redis error.")
		return
	}
	// 读入字符串
	if resultStr, err := r.Get(ctx, key).String(); err != nil {
		t.Error(err.Error())
		return
	} else if resultStr != strVal {
		t.Error("get string result error")
		return
	}

	// 写入数字
	bl, err = r.Set(ctx, key, intVal, expire1Second).Result()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !bl {
		t.Errorf("write int to redis error.")
		return
	}
	// 读入写入数字
	any := r.Get(ctx, key)
	if resultInt, err := any.Int64(); err != nil {
		t.Error(err.Error())
		return
	} else if resultInt != intVal {
		t.Error("get int result error")
		return
	}
	if resultIntStr, err := any.String(); err != nil {
		t.Error(err.Error())
		return
	} else if resultIntStr != fmt.Sprintf("%#v", intVal) {
		t.Error("get int string result error")
		return
	}

	time.Sleep(time.Second * 1)
	if err = r.Get(ctx, key).Err(); err == nil {
		t.Errorf("key %s expire error.", key)
	} else if !IsNil(err) {
		t.Errorf("key %s expire validate error. %s", key, err.Error())
	}

	if err := r.Del(ctx, key); err != nil {
		t.Errorf("delete %s  error. %s", key, err.Error())

	}
}

func Test_instance_HDel(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		fields []string
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.HDel(tt.args.ctx, tt.args.key, tt.args.fields...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HDel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_HGet(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
	}
	tests := []struct {
		name string
		args args
		want define.String
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.HGet(tt.args.ctx, tt.args.key, tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_HGetAll(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want define.MapStr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.HGetAll(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HGetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_HIncrBy(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		field string
		incr  int64
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.HIncrBy(tt.args.ctx, tt.args.key, tt.args.field, tt.args.incr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HIncrBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_HKeys(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want define.Strings
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.HKeys(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_HMGet(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		fields []string
	}
	tests := []struct {
		name string
		args args
		want define.MapStr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.HMGet(tt.args.ctx, tt.args.key, tt.args.fields...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HMGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_HSet(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.HSet(tt.args.ctx, tt.args.key, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_HSetNX(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.HSetNX(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HSetNX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceIncrAndIncrBy(t *testing.T) {
	key := fnName()

	if err := r.Del(ctx, key); err != nil {
		t.Errorf("delete key error. %s", err)
		return
	}
	val := int64(0)
	for idx := int64(0); idx < 10; idx++ {
		newVal, err := r.Incr(ctx, key).Result()
		if err != nil {
			t.Errorf("incr error. %v", err)
			return
		}
		if newVal != idx+1 {
			t.Errorf("incr error. idx: %v, value: %d", idx, val)
			return
		}
		val = newVal
	}

	for idx := int64(0); idx < 10; idx++ {
		val += idx
		newVal, err := r.IncrBy(ctx, key, idx).Result()
		if err != nil {
			t.Errorf("incrby error. %v", err)
			return
		}
		if newVal != val {
			t.Errorf("incrby error. idx: %v, value: %d", idx, val)
			return
		}
	}

}

func Test_instance_LLen(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.LLen(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_LPush(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.LPush(tt.args.ctx, tt.args.key, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LPush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_LRange(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name string
		args args
		want define.Strings
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.LRange(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_LRem(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		count int64
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.LRem(tt.args.ctx, tt.args.key, tt.args.count, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LRem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_LTrim(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name string
		args args
		want define.Bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.LTrim(tt.args.ctx, tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_MGet(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name string
		args args
		want define.Any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.MGet(tt.args.ctx, tt.args.keys...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_MSet(t *testing.T) {
	type args struct {
		ctx    context.Context
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.MSet(tt.args.ctx, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_MSetNX(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.MSetNX(tt.args.ctx, tt.args.key, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MSetNX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstancePing(t *testing.T) {
	if err := r.Ping(ctx); err != nil {
		t.Log(err.Error())
	}
}

func Test_instance_RPop(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want define.Strings
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.RPop(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RPop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_RPopLPush(t *testing.T) {
	type args struct {
		ctx         context.Context
		source      string
		destination string
	}
	tests := []struct {
		name string
		args args
		want define.String
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.RPopLPush(tt.args.ctx, tt.args.source, tt.args.destination); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RPopLPush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_RPush(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.RPush(tt.args.ctx, tt.args.key, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RPush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_Rename(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		newKey string
	}
	tests := []struct {
		name string
		args args
		want define.Bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.Rename(tt.args.ctx, tt.args.key, tt.args.newKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_RenameNX(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    string
		newKey string
	}
	tests := []struct {
		name string
		args args
		want define.Bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.RenameNX(tt.args.ctx, tt.args.key, tt.args.newKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RenameNX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_SAdd(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.SAdd(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_SMembers(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want define.Any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.SMembers(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_SRem(t *testing.T) {
	type args struct {
		ctx     context.Context
		key     string
		members []interface{}
	}
	tests := []struct {
		name string
		args args
		want define.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := instance{}
			if got := i.SRem(tt.args.ctx, tt.args.key, tt.args.members...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SRem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceSetNX(t *testing.T) {
	key := fnName()

	if err := r.Del(ctx, key); err != nil {
		t.Errorf("delete key error. %s", err)
		return
	}
	ok, err := r.SetNX(ctx, key, key, expire1Second).Result()
	if err != nil {
		t.Errorf("setnx error. %s", err)
		return
	}
	if !ok {
		t.Errorf("setnx locked failure.")
		return
	}

	ok, err = r.SetNX(ctx, key, key, expire1Second).Result()
	if err != nil {
		t.Errorf("setnx error. %s", err)
		return
	}
	if ok {
		t.Errorf("setnx locked failure. but not locked")
		return
	}

	time.Sleep(time.Second * 1)
	ok, err = r.SetNX(ctx, key, key, expire1Second).Result()
	if err != nil {
		t.Errorf("setnx error. %s", err)
		return
	}
	if !ok {
		t.Errorf("setnx locked failure.")
		return
	}
}

func redisSleep(timeout int64) {
	time.Sleep(time.Duration(timeout) * time.Second)
	// https://github.com/alicebob/miniredis/issues/259
	// https://github.com/alicebob/miniredis#ttls-key-expiration-and-time
	if mockRedis != nil {
		mockRedis.FastForward(time.Duration(timeout))
	}
}
