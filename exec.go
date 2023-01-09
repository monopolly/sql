package sql

import (
	"context"

	"github.com/monopolly/errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func (a *Conn) Exec(sql string, arg ...any) (err errors.E) {
	_, er := a.Pool.Exec(context.Background(), sql, arg...)
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
