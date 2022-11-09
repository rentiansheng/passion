package gomodule

import (
	"github.com/rentiansheng/passion/client/redis/define"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/9/5
    @desc:

***************************/

func replyInt64(num int, err error) (result intStruct) {
	result.val = int64(num)
	result.err = err
	return result
}

func replyAny(val interface{}, err error) anyStruct {
	return anyStruct{val: val, err: err}
}

func newReplyBool(val bool, err error) replyBoolStruct {
	return replyBoolStruct{
		val:   val,
		myErr: myErrType(err),
	}
}

type myErr struct {
	err error
}

func myErrType(err error) myErr {
	return myErr{err: err}
}

func (e myErr) Err() error {
	return e.err
}

type mapStrStruct struct {
	val map[string]string
	myErr
}

func (r mapStrStruct) Val() map[string]string {
	return r.val
}

func (r mapStrStruct) Result() (map[string]string, error) {
	return r.val, r.err
}

var _ define.MapStr = (*mapStrStruct)(nil)

type mapStrIntStruct struct {
	val map[string]int64
	myErr
}

func (r mapStrIntStruct) Val() map[string]int64 {
	return r.val
}

func (r mapStrIntStruct) Result() (map[string]int64, error) {
	return r.val, r.err
}

var _ define.MapStringInt = (*mapStrIntStruct)(nil)

type intStruct struct {
	val int64
	myErr
}

func (r intStruct) Val() int64 {
	return r.val
}

func (r intStruct) Result() (int64, error) {
	return r.val, r.err
}

var _ define.Int = (*intStruct)(nil)

type stringStruct struct {
	val string
	myErr
}

func (r stringStruct) Val() string {
	return r.val
}

func (r stringStruct) Result() (string, error) {
	return r.val, r.err
}

var _ define.String = (*stringStruct)(nil)

type stringsStruct struct {
	val []string
	myErr
}

func (r stringsStruct) Val() []string {

	return r.val
}

func (r stringsStruct) Result() ([]string, error) {
	return r.val, r.err
}

var _ define.Strings = (*stringsStruct)(nil)

type statusStruct struct {
	val string
	myErr
}

func (r statusStruct) Val() string {
	return r.val
}

func (r statusStruct) Result() (string, error) {
	return r.val, r.err
}

var (
	_ define.Status = (*statusStruct)(nil)
)

type replyBoolStruct struct {
	val bool
	myErr
}

func (b replyBoolStruct) Val() bool {
	return b.val
}

func (b replyBoolStruct) Result() (bool, error) {
	return b.val, b.err
}

var _ define.Bool = (*replyBoolStruct)(nil)
