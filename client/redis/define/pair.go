package define

import (
	"encoding/json"
	"fmt"

	"github.com/rentiansheng/passion/lib/bytesconv"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/9/4
    @desc:

***************************/

type Pair interface {
	Pair() (string, []byte)
}

func StrPair(key, val string) Pair {
	return pairValueString{key: key, val: val}
}

func BytesPair(key string, b []byte) Pair {
	return pairValueBytes{key: key, val: b}
}

func AnyJSONPair(key string, val interface{}) (Pair, error) {
	return anyJSONPair{}.Pair(key, val)
}

func MapJSONPair(data map[string]interface{}) ([]Pair, error) {

	return mapJSONPair{}.Pair(data)
}

type pairValueString struct {
	key string
	val string
}

func (p pairValueString) Pair() (string, []byte) {
	return p.key, bytesconv.StringToBytes(p.val)
}

type pairValueBytes struct {
	key string
	val []byte
}

func (k pairValueBytes) Pair() (string, []byte) {
	return k.key, k.val
}

type anyJSONPair struct {
}

func (a anyJSONPair) Pair(key string, value interface{}) (Pair, error) {

	bytesVal, err := AnyJSONBytes(value)
	if err != nil {
		return nil, fmt.Errorf("key %s convert to json byte error.  %s", key, err.Error())
	}

	return BytesPair(key, bytesVal[1:len(bytesVal)-1]), nil
}

type mapJSONPair struct {
}

func (k mapJSONPair) Pair(data map[string]interface{}) ([]Pair, error) {
	pairs := make([]Pair, 0, len(data))
	for key, val := range data {
		pair, err := AnyJSONPair(key, val)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}

func AnyJSONBytes(value interface{}) ([]byte, error) {
	switch val := value.(type) {
	case []byte:
		return val, nil
	case string:
		return []byte(val), nil
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		str := fmt.Sprintf("%v", val)
		return []byte(str), nil
	default:
		bytesVal, err := json.Marshal([]interface{}{value})
		if err != nil {
			return nil, err
		}

		return bytesVal[1 : len(bytesVal)-1], nil
	}
}
