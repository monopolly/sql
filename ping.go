package sql

import (
	"context"

	"github.com/monopolly/errors"
)

func (a *Conn) Ping() (err errors.E) {
	er := a.Pool.Ping(context.Background())
	if er != nil {
		return errors.Connection(er)
	}
	return
}
