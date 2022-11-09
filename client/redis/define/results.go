package define

import "time"

/***************************
    @author: tiansheng.ren
    @date: 2022/9/4
    @desc:

***************************/

// base is the base result for redis commands
type base interface {
	Err() error
}

// Result is the common result for redis commands
type Result interface {
	base
	Val() interface{}
	Result() (interface{}, error)
}

// String is the string result for redis commands
type String interface {
	base
	Val() string
	Result() (string, error)
}

// Float  is the float result for redis commands
type Float interface {
	base
	Val() float64
	Result() (float64, error)
}

// Int is the int result for redis commands
type Int interface {
	base
	Val() int64
	Result() (int64, error)
}

// Slice  is the slice result for redis commands
type Slice interface {
	base
	Val() []interface{}
	Result() ([]interface{}, error)
}

// Status  is the status result for redis commands
type Status interface {
	base
	Val() string
	Result() (string, error)
}

// Bool the bool result for redis commands
type Bool interface {
	base
	Val() bool
	Result() (bool, error)
}

// Ints is the int slice result for redis commands
type Ints interface {
	base
	Val() []int64
	Result() ([]int64, error)
}

// Strings is the string slice result for redis commands
type Strings interface {
	base
	Val() []string
	Result() ([]string, error)
}

// Bools is the bool slice result for redis commands
type Bools interface {
	base
	Val() []bool
	Result() ([]bool, error)
}

// MapStr is the kv. key and value is string for redis commands
type MapStr interface {
	base
	Val() map[string]string
	Result() (map[string]string, error)
}

// MapStringInt is the string int map result for redis commands
type MapStringInt interface {
	base
	Val() map[string]int64
	Result() (map[string]int64, error)
}

// MapStringObject is the string struct map result for redis commands
type MapStringObject interface {
	base
	Val() map[string]struct{}
	Result() (map[string]struct{}, error)
}

// Duration is the duration result for redis commands
type Duration interface {
	base
	Val() time.Duration
	Result() (time.Duration, error)
}

// Scan is the duration result for redis commands
type Scan interface {
	base
	Val() (keys []string, cursor uint64)
	Result() (keys []string, cursor uint64, err error)
}

type Any interface {
	base
	any
}

type any interface {
	Bytes() ([]byte, error)
	String() (string, error)
	Int64() (int64, error)
	Uint64() (uint64, error)
	Float64() (float64, error)
	Bool() (bool, error)
	Int64s() ([]int64, error)
	Map() (map[string]interface{}, error)
	MapStr() (map[string]string, error)
	MapInt() (map[string]int, error)
	MapInt64() (map[string]int64, error)
	Decode(val interface{}) error
	Float64s() ([]float64, error)
	Strings() ([]string, error)
}
