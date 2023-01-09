package sql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/monopolly/errors"
)

func (a *Conn) Exist(sql string, arg ...any) (resp bool, err errors.E) {
	r, er := a.Pool.Query(context.Background(), sql+" limit 1", arg...)
	if er != nil {
		err = errors.Database(er)
		return
	}
	return r.Next(), nil
}

func (a *Conn) HasID(table string, id any) (has bool, err errors.E) {
	q := fmt.Sprintf("select id from %s where id = $1 limit 1", table)
	r := a.Pool.QueryRow(context.Background(), q, id)
	var h any
	er := r.Scan(&h)
	if er != nil {
		if er == pgx.ErrNoRows {
			return
		}
		err = errors.Database(er)
		return
	}
	has = true
	return
}
