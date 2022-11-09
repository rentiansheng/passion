package array

import (
	"encoding/json"
	"fmt"
	"reflect"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/6/23
    @desc:

***************************/

var (
	Int = intType("int array tools")
)

type intType string

func (i intType) Join(items []int, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) JoinU(items []uint, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) Join64(items []int64, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) JoinU64(items []uint64, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) Join32(items []int32, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) JoinU32(items []uint32, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) Join16(items []int32, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) JoinU16(items []uint32, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) Join8(items []int32, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) JoinU8(items []uint32, sep string) string {
	itemsLen := len(items)
	switch itemsLen {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", items[0])
	}

	result := fmt.Sprintf("%d", items[0])
	for _, val := range items[1:] {
		result += fmt.Sprintf("%s%d", sep, val)
	}

	return result
}

func (i intType) Exist64(existId int64, ids []int64) bool {
	for _, id := range ids {
		if existId == id {
			return true
		}
	}
	return false
}

// ToIntSlice input: slice or array of int64, uint, uint64.
// output: slice of int
func (i intType) ToIntSlice(v interface{}) ([]int, error) {
	res := make([]int, 0)
	switch list := v.(type) {
	case []int64:
		for _, num := range list {
			res = append(res, int(num))
		}
	case []uint:
		for _, num := range list {
			res = append(res, int(num))
		}
	case []uint64:
		for _, num := range list {
			res = append(res, int(num))
		}
	default:
		return nil, fmt.Errorf("ToIntSlice does not support type of %v", reflect.TypeOf(v).Kind().String())
	}
	return res, nil
}

// ToIntSlice input: slice or array of int64, uint, uint64.
// output: slice of int
func (i intType) ToUint64Slice(v interface{}) ([]uint64, error) {
	res := make([]uint64, 0)
	switch list := v.(type) {
	case []int64:
		for _, num := range list {
			res = append(res, uint64(num))
		}
	case []uint:
		for _, num := range list {
			res = append(res, uint64(num))
		}
	case []int:
		for _, num := range list {
			res = append(res, uint64(num))
		}
	default:
		return nil, fmt.Errorf("ToIntSlice does not support type of %v", reflect.TypeOf(v).Kind().String())
	}
	return res, nil
}

func (i intType) ToUint64(v interface{}) (uint64, error) {
	var res uint64
	switch val := v.(type) {
	case int:
		res = uint64(val)
	case int64:
		res = uint64(val)
	case float32:
		res = uint64(val)
	case float64:
		res = uint64(val)
	case json.Number:
		valInt64, err := val.Int64()
		if err != nil {
			return 0, err
		}
		res = uint64(valInt64)
	default:
		return 0, fmt.Errorf("ToUint64ErrCode does not support type of %v", reflect.TypeOf(v).Kind().String())
	}

	return res, nil
}
