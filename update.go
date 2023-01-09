package sql

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/monopolly/errors"
	"github.com/monopolly/jsons"
)

// update fields
func (a *Conn) Update(table string, id, k, v any) (err errors.E) {
	q := fmt.Sprintf("update %s set %s = $1 where id = $2", table, k)
	_, er := a.Pool.Exec(context.Background(), q, v, id)
	if er != nil {
		b, _ := er.(*pgconn.PgError)
		if b != nil {
			switch b.Code {
			case pgerrcode.UniqueViolation:
				err = errors.Exist()
			default:
				err = errors.Database(er)
				err.Set("sqlcode", b.Code)
			}
		} else {
			err = errors.Database(er)
		}
	}
	return
}

// update jsonb: profiles set meta = meta || '{"name":"James Wood"}'::jsonb where id = 1;
func (a *Conn) UpdateJsonbMap(table, mapfield, k string, id, v any) (err errors.E) {
	res := jsons.Creates(k, v).String()
	q := fmt.Sprintf(`update %[1]s set %[2]s = %[2]s || $$%[3]s$$::jsonb where id = $1;`, table, mapfield, JsonbEscape(res))
	_, er := a.Pool.Exec(context.Background(), q, id)
	if er != nil {
		err = errors.Database(er)
	}
	return
}

// update jsonb batch: profiles set meta = meta || '{"name":"James Wood"}'::jsonb where id = 1;
func (a *Conn) UpdateJsonbMapList(table, mapfield string, id any, keys map[string]any) (err errors.E) {
	res := jsons.Create()
	for k, v := range keys {
		res.Add(k, v)
	}
	q := fmt.Sprintf(`update %[1]s set %[2]s = %[2]s || $$%[3]s$$::jsonb where id = $1;`, table, mapfield, JsonbEscape(res.String()))
	_, er := a.Pool.Exec(context.Background(), q, id)
	if er != nil {
		err = errors.Database(er)
	}
	return
}
