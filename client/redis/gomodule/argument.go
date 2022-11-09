package gomodule

import (
	"encoding/json"
	"fmt"
	"reflect"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/9/5
    @desc:

***************************/

func mergeArgs(key interface{}, args []interface{}) []interface{} {
	return append([]interface{}{key}, args...)
}

func mergeStrArgs(key string, args []string) []interface{} {
	results := make([]interface{}, len(args)+1)
	results[0] = key
	for idx, key := range args {
		results[idx+1] = key
	}
	return results
}

type errorFunc interface {
	Error() string
}
type stringFunc interface {
	String() string
}

func convToString(args ...interface{}) []string {
	results := make([]string, len(args))
	for idx, arg := range args {

		if f, ok := arg.(errorFunc); ok {
			results[idx] = f.Error()
			continue
		}
		if f, ok := arg.(stringFunc); ok {
			results[idx] = f.String()
			continue
		}
		if f, ok := arg.(string); ok {
			results[idx] = f
			continue
		}
		if f, ok := arg.([]byte); ok {
			results[idx] = string(f)
			continue
		}

		if arg == nil {
			results[idx] = "null"
			continue
		}

		valTyp := reflect.TypeOf(arg)
		for valTyp.Kind() == reflect.Ptr {
			valTyp = valTyp.Elem()
		}
		kind := valTyp.Kind()

		if kind == reflect.Struct || kind == reflect.Interface ||
			kind == reflect.Array || kind == reflect.Map || kind == reflect.Slice {
			if out, err := json.Marshal(arg); err != nil {
				results[idx] = fmt.Sprintf("%#v", arg)
			} else {
				results[idx] = string(out)
			}
		} else {
			results[idx] = fmt.Sprintf("%#v", arg)
		}
	}

	return results
}

func convToStringInterface(args []string) []interface{} {
	results := make([]interface{}, len(args))
	for idx, key := range args {
		results[idx] = key
	}
	return results
}
