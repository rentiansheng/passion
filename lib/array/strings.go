package array

import (
	"strconv"
	"strings"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/5/7
    @desc:
		处理[]string 常见方法

***************************/

var (
	Strings = strType("string array tools")
)

type strType string

// Intersection 求两个数据交集
func (strType) Intersection(arr1, arr2 []string) []string {
	tmpMap := make(map[string]struct{}, len(arr1))
	for _, item := range arr1 {
		tmpMap[item] = struct{}{}
	}
	var ret []string
	for _, item := range arr2 {
		if _, ok := tmpMap[item]; ok {
			ret = append(ret, item)
		}
	}

	return ret

}

func (strType) StringToIntArray(str string, split string) ([]int64, error) {
	strItems := strings.Split(str, split)
	ret := make([]int64, 0, len(strItems))
	for _, chartConfigIdStr := range strItems {
		itemVal, err := strconv.ParseInt(chartConfigIdStr, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, itemVal)
	}
	return ret, nil
}

func (strType) InStrings(str string, arr []string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

func (strType) SQLINArray(arr []string) string {
	results := strings.Join(arr, "','")
	if results != "" {
		results = "'" + results + "'"
	}
	return results
}

// Diff 求两个数组有差异数据
func (strType) Diff(arr1, arr2 []string) []string {

	tmp1Map := make(map[string]struct{}, len(arr1))
	for _, item := range arr1 {
		tmp1Map[item] = struct{}{}
	}
	tmp2Map := make(map[string]struct{}, len(arr2))
	for _, item := range arr2 {
		tmp2Map[item] = struct{}{}
	}
	existMap := make(map[string]struct{})
	var ret []string

	fn := func(arr []string, extraArrMap map[string]struct{}) {
		for _, key := range arr {
			if _, ok := extraArrMap[key]; !ok {
				if _, ok := existMap[key]; ok {
					continue
				}
				existMap[key] = struct{}{}
				ret = append(ret, key)
			}
		}
	}
	fn(arr1, tmp2Map)
	fn(arr2, tmp1Map)
	return ret

}

func (strType) ToInterface(args []string) []interface{} {
	results := make([]interface{}, len(args))
	for idx, val := range args {
		results[idx] = val
	}
	return results
}
