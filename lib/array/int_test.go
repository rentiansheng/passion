package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/6/23
    @desc:

***************************/

func TestIntType_Join(t *testing.T) {

	suits := []struct {
		input  []int
		sep    string
		output string
	}{
		{
			nil,
			",",
			"",
		},
		{
			[]int{},
			",",
			"",
		},
		{
			[]int{1},
			",",
			"1",
		},
		{
			[]int{1, 1},
			",",
			"1,1",
		},
		{
			[]int{1, 1, 2, 3, 4, 5, 6, 7, 81000, 1},
			",",
			"1,1,2,3,4,5,6,7,81000,1",
		},
		{
			[]int{1, 1, 2, 3, 4, 5, 6, 7, 81000, 1},
			"#",
			"1#1#2#3#4#5#6#7#81000#1",
		},
		{
			[]int{1, 1, 2, 3, 4, 5, 6, 7, 81000, 1},
			"",
			"11234567810001",
		},
	}
	for idx, suit := range suits[6:] {
		output := Int.Join(suit.input, suit.sep)
		require.Equal(t, suit.output, output, "test suit index %d", idx)
	}

}

func TestToIntSlice(t *testing.T) {
	list := []uint64{1, 2, 3}
	res, err := Int.ToIntSlice(list)
	assert.Nil(t, err)
	_ = res
}
