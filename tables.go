package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/monopolly/errors"
)

// insert
func (a *Conn) CreateTable(table string, keys map[string]string) (err errors.E) {
	var list []string
	for k, v := range keys {
		list = append(list, fmt.Sprintf("%s  %s", k, v))
	}
	q := fmt.Sprintf("create table if not exists %s (\n%s\n)", table, strings.Join(list, ",\n"))
	_, er := a.Pool.Exec(context.Background(), q)
	if er != nil {
		err = errors.Database(er)
	}
	return
}

// insert
func (a *Conn) DropTable(table string) (err errors.E) {
	_, er := a.Pool.Exec(context.Background(), "drop table "+table)
	if er != nil {
		err = errors.Database(er)
	}
	return
}
