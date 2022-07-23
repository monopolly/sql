package sql

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/monopolly/errors"
)

type Conn struct {
	pool *pgx.ConnPool
}

func (a *Conn) Pool() *pgx.ConnPool {
	return a.pool
}

//will added limit 1 to the end
func (a *Conn) Query(sql string, arg ...interface{}) (resp [][]interface{}, err errors.E) {
	r, er := a.pool.Query(sql, arg...)
	if er != nil {
		b, _ := er.(*pgconn.PgError)
		if b != nil {
			switch b.Code {
			case pgerrcode.UniqueViolation:
				err = errors.Exist(er)
			case pgerrcode.NoDataFound:
				err = errors.NotFound()
			default:
				err = errors.Database(er)
				err.Set("sql", b.Code)
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

//response as json: select only!
//example: select * from users = [][]byte
func (a *Conn) QueryJson(sql string, arg ...interface{}) (resp [][]byte, err errors.E) {
	r, er := a.pool.Query(JsonResult(sql), arg...)
	if er != nil {
		err = errors.Database(er)
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
func (a *Conn) Row(sql string, arg ...interface{}) (resp []interface{}, err errors.E) {
	r, er := a.pool.Query(sql+" limit 1", arg...)
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

func (a *Conn) Exist(sql string, arg ...interface{}) (resp bool, err errors.E) {
	r, er := a.pool.Query(sql+" limit 1", arg...)
	if er != nil {
		err = errors.Database(er)
		return
	}
	defer r.Close()
	return r.Next(), nil
}

//will added limit 1 to the end
func (a *Conn) RowJson(sql string, arg ...interface{}) (resp []byte, err errors.E) {
	r, er := a.pool.Query(JsonResult(sql+" limit 1"), arg...)
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

func (a *Conn) Exec(sql string, arg ...interface{}) (err errors.E) {
	_, er := a.pool.Exec(sql, arg...)
	if er != nil {
		b, _ := er.(*pgconn.PgError)
		if b != nil {
			switch b.Code {
			case pgerrcode.UniqueViolation:
				err = errors.Exist(er)
			default:
				err = errors.Database(er)
				err.Set("sql", b.Code)
			}
		} else {
			err = errors.Database(er)
		}
	}
	return
}

/* update events set usr = usr || '{"name":"James Wood"}'::jsonb where uid = 1; */
