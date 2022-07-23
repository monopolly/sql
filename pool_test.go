package sql

//testing

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Nil(i interface{}) {
	if i != nil {
		fmt.Println("FAIL:", i)
		return
	}
}

func NotNil(i interface{}) {
	if i == nil {
		fmt.Println("FAIL:", i)
		return
	}
}

func TestConn(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	a := assert.New(ggggg)
	_ = a

}
