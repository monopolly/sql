package sql

//testing

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestTools(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	fmt.Println(JsonResult("select * from users"))
}
