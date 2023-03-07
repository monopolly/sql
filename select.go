package sql

import (
	"bytes"
	"context"
	"fmt"

	"github.com/monopolly/errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// will added limit 1 to the end
func (a *Conn) Query(sql string, arg ...any) (resp [][]any, err errors.E) {

	r, er := a.Pool.Query(context.Background(), sql, arg...)
	if er != nil {
		b, _ := er.(*pgconn.PgError)
		if b != nil {
			switch b.Code {
			case pgerrcode.UniqueViolation:
				err = errors.Exist()
			case pgerrcode.NoDataFound:
				err = errors.NotFound()
			default:
				err = errors.Database(er)
				err.Set("sqlcode", b.Code)
			}
		} else {
			err = errors.Database(er)
		}
		return
	}

	for r.Next() {
		v, err := r.Values()
		if err != nil {
			continue
		}
		resp = append(resp, v)
	}

	r.Close()
	return
}

// response as json: select only!
// example: select * from users = [][]byte
func (a *Conn) QueryJson(sql string, arg ...any) (resp [][]byte, err errors.E) {
	r, er := a.Pool.Query(context.Background(), JsonResult(sql), arg...)
	if er != nil {
		err = errors.Database(er)
		return
	}

	for r.Next() {
		raw := r.RawValues()
		if len(raw) == 0 {
			continue
		}
		x := raw[0]
		resp = append(resp, x)
	}

	return
}

// response as json: select only!
// example: select * from users = [][]byte
func (a *Conn) QueryJsonArray(sql string, arg ...any) (resp []byte, err errors.E) {
	r, er := a.Pool.Query(context.Background(), JsonResult(sql), arg...)
	if er != nil {
		err = errors.Database(er)
		return
	}
	defer r.Close()

	var w bytes.Buffer
	w.WriteString("[")
	var c int
	for r.Next() {
		raw := r.RawValues()
		if len(raw) == 0 {
			continue
		}
		w.Write(raw[0])
		w.WriteString(",")
		c++
	}
	if c == 0 {
		return
	}
	resp = w.Bytes()
	resp[len(resp)-1] = ']'
	return
}

// all records from table in json
func (a *Conn) All(table string) (json [][]byte, err errors.E) {
	q := fmt.Sprintf("select * from %s", table)
	return a.QueryJson(q)
}

// all records from table in json
func (a *Conn) AllArray(table string) (res []byte, err errors.E) {
	q := fmt.Sprintf("select * from %s", table)
	return a.QueryJsonArray(q)
}

// will added limit 1 to the end
func (a *Conn) Row(sql string, arg ...any) (resp []any, err errors.E) {
	r, er := a.Pool.Query(context.Background(), sql+" limit 1", arg...)
	if er != nil {
		err = errors.Database(er)
		return
	}
	defer r.Close()
	if !r.Next() {
		err = errors.NotFound()
		return
	}

	resp, er = r.Values()
	if er != nil {
		err = errors.Database(er)
		return
	}

	return
}

// will added limit 1 to the end
func (a *Conn) RowJson(sql string, arg ...any) (resp []byte, err errors.E) {
	r, er := a.Pool.Query(context.Background(), JsonResult(sql+" limit 1"), arg...)
	if er != nil {
		err = errors.Database(er)
		return
	}
	defer r.Close()

	if r.Next() {
		r.Scan(&resp)
	} else {
		err = errors.NotFound()
		return
	}
	return
}
