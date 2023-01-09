package sql

//testing

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	dockername = "testdocker"
	table      = "testtable"
	connstring = "postgres://postgres:12345@localhost:9111"
)

func TestConnection(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	a := assert.New(ggggg)
	_ = a

	// postgres
	cmd("docker", "stop", dockername)
	cmd("docker", "rm", dockername)
	defer cmd("docker", "rm", dockername)
	defer cmd("docker", "stop", dockername)
	cmd("docker", "run", "-d", "--name", dockername, "-p", "9111:5432", "-e", "POSTGRES_PASSWORD=12345", "postgres")
	time.Sleep(time.Second * 2)

	conn, er := NewString(connstring)
	if er != nil {
		panic(er)
	}

	er = conn.Ping()
	if er != nil {
		panic(er)
	}

	fmt.Println("connection ok")

	er = conn.CreateTable(table, map[string]string{"id": "bigserial", "title": "text", "bools": "boolean", "created": "bigint", "meta": "jsonb default '{}'::jsonb"})
	if er != nil {
		panic(er)
	}

	id, er := conn.Insert(table, map[string]any{"title": "text", "created": time.Now().Unix()})
	if er != nil {
		panic(er)
	}
	fmt.Println("return id", id)

	er = conn.Update(table, 1, "title", "changed!")
	if er != nil {
		panic(er)
	}

	er = conn.Update(table, 1, "bools", true)
	if er != nil {
		panic(er)
	}

	er = conn.UpdateJsonbMapList(table, "meta", 1, map[string]any{"jsontext": "OK", "num": 111, "bools": true})
	if er != nil {
		panic(er)
	}

	has, er := conn.HasID(table, 1)
	if er != nil {
		panic(er)
	}

	fmt.Println("has", has)

	r, er := conn.RowJson("select * from " + table)
	if er != nil {
		panic(er)
	}
	fmt.Println(string(r))

}

func cmd(name string, v ...string) {
	r, err := exec.Command(name, v...).Output()
	fmt.Println(string(r), err)
}
