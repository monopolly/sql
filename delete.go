package sql

import (
	"context"
	"fmt"

	"github.com/monopolly/errors"
)

// will added limit 1 to the end
func (a *Conn) Delete(table string, id any) (err errors.E) {
	q := fmt.Sprintf("delete from %s where id = $1", table)
	_, er := a.Pool.Exec(context.Background(), q, id)
	if er != nil {
		err = errors.Database(er)
		return
	}
	return
}
