package gomodule

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"

	"github.com/rentiansheng/passion/client/redis/define"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/9/6
    @desc:

***************************/

type anyStruct struct {
	val interface{}
	err error
}

func (a anyStruct) Err() error {
	return a.err
}

func (a anyStruct) Bytes() ([]byte, error) {
	return redis.Bytes(a.val, a.err)
}

func (a anyStruct) String() (string, error) {
	return redis.String(a.val, a.err)
}

func (a anyStruct) Int64() (int64, error) {
	return redis.Int64(a.val, a.err)
}

func (a anyStruct) Int() (int, error) {
	return redis.Int(a.val, a.err)
}

func (a anyStruct) Uint64() (uint64, error) {
	return redis.Uint64(a.val, a.err)

}

func (a anyStruct) Float64() (float64, error) {
	return redis.Float64(a.val, a.err)

}

func (a anyStruct) Bool() (bool, error) {
	if a.err != nil {
		return false, a.err
	}
	if strVal, ok := a.val.(string); ok {
		if strings.ToLower(strVal) == "ok" {
			return true, nil
		}
	}
	ok, err := redis.Bool(a.val, a.err)
	if err != nil {
		if IsNil(err) {
			return false, nil
		}
		return false, err
	} else {
		return ok, nil
	}

	return false, nil
}

func (a anyStruct) Int64s() ([]int64, error) {
	return redis.Int64s(a.val, a.err)
}

func (a anyStruct) Map() (map[string]interface{}, error) {
	bytes, err := redis.Bytes(a.val, a.err)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]interface{}, 0)
	if err := json.Unmarshal(bytes, &ret); err != nil {
		return nil, err
	}
	return ret, err
}

func (a anyStruct) MapStr() (map[string]string, error) {
	return redis.StringMap(a.val, a.err)
}

func (a anyStruct) MapInt64() (map[string]int64, error) {
	return redis.Int64Map(a.val, a.err)
}

func (a anyStruct) MapInt() (map[string]int, error) {
	return redis.IntMap(a.val, a.err)
}

func (a anyStruct) Decode(result interface{}) error {
	bytes, err := redis.Bytes(a.val, a.err)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, result)
}

func (a anyStruct) Float64s() ([]float64, error) {
	return redis.Float64s(a.val, a.err)
}

func (a anyStruct) Strings() ([]string, error) {
	return redis.Strings(a.val, a.err)
}

func (a anyStruct) ToString() (result stringStruct) {
	result.val, result.err = a.String()
	return result
}

func (a anyStruct) ToValString() (result stringStruct) {

	values, err := a.Strings()
	if err != nil {
		result.err = err
		return
	} else if len(values) == 2 {
		result.val = values[1]
		return
	} else {
		result.err = fmt.Errorf("not value string protocol")
	}

	return result
}

func (a anyStruct) ToValStrings() (result stringsStruct) {

	values, err := a.Strings()
	if err != nil {
		result.err = err
		return
	}
	for idx := 1; idx < len(values); idx += 2 {
		result.val = append(result.val, values[idx])
	}

	return result
}

func (a anyStruct) ToStrings() (result stringsStruct) {
	result.val, result.err = a.Strings()
	return result
}

func (a anyStruct) ToStatus() (result statusStruct) {
	result.val, result.err = a.String()
	return result
}

func (a anyStruct) ToBool() (result replyBoolStruct) {

	result.val, result.err = a.Bool()

	return result
}

func (a anyStruct) ToInt() (result intStruct) {

	result.val, result.err = a.Int64()
	return
}

func (a anyStruct) ToMapStr() (result mapStrStruct) {
	result.val, result.err = a.MapStr()
	return
}

var _ define.Any = (*anyStruct)(nil)
