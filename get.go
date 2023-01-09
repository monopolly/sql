package sql

import (
	"context"
	"fmt"

	"github.com/monopolly/errors"
)

// will added limit 1 to the end
func (a *Conn) GetJson(id any, table string) (resp []byte, err errors.E) {
	q := fmt.Sprintf("select * from %s where id = $1 limit 1", table)
	q = JsonResult(q)
	r, er := a.Pool.Query(context.Background(), q, id)
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
