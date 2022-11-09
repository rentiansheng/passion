package array

import "testing"

/***************************
    @author: tiansheng.ren
    @date: 2022/5/10
    @desc:

***************************/

func TestStringIntersection(t *testing.T) {
	type inputItem struct {
		arr1 []string
		arr2 []string
	}
	inputs := []struct {
		input  inputItem
		result []string
	}{
		{
			input:  inputItem{arr1: []string{"1", "2", "3"}, arr2: []string{}},
			result: []string{},
		},
		{
			input:  inputItem{arr1: []string{}, arr2: []string{"1", "2", "3"}},
			result: []string{},
		},
		{
			input:  inputItem{arr1: []string{"1", "2", "3"}, arr2: []string{"1", "2", "3"}},
			result: []string{"1", "2", "3"},
		},
		{
			input:  inputItem{arr1: []string{"1"}, arr2: []string{"1", "2", "3"}},
			result: []string{"1"},
		},
		{
			input:  inputItem{arr1: []string{"1", "2", "3"}, arr2: []string{"1"}},
			result: []string{"1"},
		},
		{
			input:  inputItem{arr1: []string{"2", "3", "3", "3", "3", "3", "4", "6", "7"}, arr2: []string{"1", "7", "3"}},
			result: []string{"7", "3"},
		},
		{
			input:  inputItem{arr1: []string{"2", "3", "3", "3", "3", "3", "4", "6", "7"}, arr2: []string{"2", "2", "2", "2", "2", "2"}},
			result: []string{"2", "2", "2", "2", "2", "2"},
		},
	}

	for idx, testInput := range inputs {
		testRes := Strings.Intersection(testInput.input.arr1, testInput.input.arr2)
		if len(testRes) != len(testInput.result) {
			t.Errorf("test index: %d, expect: %#v, actual: %#v", idx, testInput.result, testRes)
			continue
		}
		testResMap := make(map[string]struct{}, 0)
		for _, res := range testRes {
			testResMap[res] = struct{}{}
		}
		for _, expectItem := range testInput.result {
			if _, ok := testResMap[expectItem]; !ok {
				t.Errorf("test index: %d, expect: %#v, actual: %#v", idx, testInput.result, testRes)
				continue
			}
		}
	}

}

func TestStringToIntArray(t *testing.T) {
	type inputItem struct {
		string string
		split  string
		output []int64
		hasErr bool
	}
	suits := []inputItem{
		{
			string: "1,22,333",
			split:  ",",
			output: []int64{1, 22, 333},
		},
		{
			string: "10000,010,0000",
			split:  ",",
			output: []int64{10000, 10, 0},
		},
		{
			string: "1000",
			split:  ",",
			output: []int64{1000},
		},
		{
			string: "1#1",
			split:  "#",
			output: []int64{1, 1},
		},
		{
			string: "1121a",
			split:  ",",
			output: []int64{},
			hasErr: true,
		},
	}
	for idx, suit := range suits {
		ret, err := Strings.StringToIntArray(suit.string, suit.split)
		if err != nil {
			if !suit.hasErr {
				t.Errorf("test index %d error. has err: %s", idx, err.Error())
			}
			continue
		}
		if len(ret) != len(suit.output) {
			t.Errorf("test index %d error. actual len: %d, expect len: %d", idx, len(ret), len(suit.output))
			continue
		}
		tmpVal := numberIntersection(ret, suit.output)
		if len(tmpVal) != len(suit.output) {
			t.Errorf("test index %d error. actual: %#v, expect: %#v", idx, ret, suit.output)
		}
	}

}

func numberIntersection(arr1, arr2 []int64) []int64 {
	tmpMap := make(map[int64]struct{}, len(arr1))
	for _, item := range arr1 {
		tmpMap[item] = struct{}{}
	}
	var ret []int64
	for _, item := range arr2 {
		if _, ok := tmpMap[item]; ok {
			ret = append(ret, item)
		}
	}

	return ret

}

func TestStringDiff(t *testing.T) {
	type inputItem struct {
		arr1 []string
		arr2 []string
	}
	inputs := []struct {
		input  inputItem
		result []string
	}{
		{
			input:  inputItem{arr1: []string{"1", "2", "3"}, arr2: []string{}},
			result: []string{"1", "2", "3"},
		},
		{
			input:  inputItem{arr1: []string{}, arr2: []string{"1", "2", "3"}},
			result: []string{"1", "2", "3"},
		},
		{
			input:  inputItem{arr1: []string{"1", "2", "3"}, arr2: []string{"1", "2", "3"}},
			result: []string{},
		},
		{
			input:  inputItem{arr1: []string{"1"}, arr2: []string{"1", "2", "3"}},
			result: []string{"2", "3"},
		},
		{
			input:  inputItem{arr1: []string{"1", "2", "3"}, arr2: []string{"1"}},
			result: []string{"2", "3"},
		},
		{
			input:  inputItem{arr1: []string{"2", "3", "3", "3", "3", "3", "4", "6", "7"}, arr2: []string{"1", "7", "3"}},
			result: []string{"2", "4", "6", "1"},
		},
		{
			input:  inputItem{arr1: []string{"2", "3", "3", "3", "3", "3", "4", "6", "7"}, arr2: []string{"2", "2", "2", "2", "2", "2"}},
			result: []string{"3", "4", "6", "7"},
		},
	}

	for idx, testInput := range inputs {
		testRes := Strings.Diff(testInput.input.arr1, testInput.input.arr2)
		if len(testRes) != len(testInput.result) {
			t.Errorf("test index: %d, expect: %#v, actual: %#v", idx, testInput.result, testRes)
			continue
		}
		testResMap := make(map[string]struct{}, 0)
		for _, res := range testRes {
			testResMap[res] = struct{}{}
		}
		for _, expectItem := range testInput.result {
			if _, ok := testResMap[expectItem]; !ok {
				t.Errorf("test index: %d, expect: %#v, actual: %#v", idx, testInput.result, testRes)
				continue
			}
		}
	}

}
