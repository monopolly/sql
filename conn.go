package sql

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type Conn struct {
	pool *pgx.ConnPool
}

func (a *Conn) Pool() *pgx.ConnPool {
	return a.pool
}

//will added limit 1 to the end
func (a *Conn) Query(sql string, arg ...interface{}) (resp [][]interface{}, err error) {
	r, er := a.pool.Query(sql, arg...)
	if er != nil {
		b, _ := er.(*pgconn.PgError)
		if b != nil {
			switch b.Code {
			case pgerrcode.UniqueViolation:
				err = errors.New("exist")
			case pgerrcode.NoDataFound:
				err = errors.New("notfound")
			default:
				err = er
			}
		} else {
			err = er
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

//response as json: select only!
//example: select * from users = [][]byte
func (a *Conn) QueryJson(sql string, arg ...interface{}) (resp [][]byte, err error) {
	r, er := a.pool.Query(JsonResult(sql), arg...)
	if er != nil {
		err = errors.New("database")
		return
	}

	for r.Next() {
		raw, _ := r.Values()
		if len(raw) == 0 {
			continue
		}
		//resp = append(resp, raw[0])
	}

	r.Close()
	return
}

//will added limit 1 to the end
func (a *Conn) Row(sql string, arg ...interface{}) (resp []interface{}, err error) {
	r, er := a.pool.Query(sql+" limit 1", arg...)
	if er != nil {
		err = errors.New("database")
		return
	}
	defer r.Close()
	if !r.Next() {
		err = errors.New("notfound")
		return
	}

	resp, er = r.Values()
	if er != nil {
		err = errors.New("database")
		return
	}

	return
}

func (a *Conn) Exist(sql string, arg ...interface{}) (resp bool, err error) {
	r, er := a.pool.Query(sql+" limit 1", arg...)
	if er != nil {
		err = errors.New("database")
		return
	}
	defer r.Close()
	return r.Next(), nil
}

//will added limit 1 to the end
func (a *Conn) RowJson(sql string, arg ...interface{}) (resp []byte, err error) {
	r, er := a.pool.Query(JsonResult(sql+" limit 1"), arg...)
	if er != nil {
		err = errors.New("database")
		return
	}
	defer r.Close()

	if r.Next() {
		r.Scan(&resp)
	} else {
		err = errors.New("notfound")
		return
	}
	return
}

func (a *Conn) Exec(sql string, arg ...interface{}) (err error) {
	_, er := a.pool.Exec(sql, arg...)
	if er != nil {
		b, _ := er.(*pgconn.PgError)
		if b != nil {
			switch b.Code {
			case pgerrcode.UniqueViolation:
				err = errors.New("exist")
			default:
				err = errors.New("database")

			}
		} else {
			err = errors.New("database")
		}
	}
	return
}

/* update events set usr = usr || '{"name":"James Wood"}'::jsonb where uid = 1; */
