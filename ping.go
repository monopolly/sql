package sql

import (
	"context"
)

func (a *Conn) Ping() (err error) {
	return a.Pool.Ping(context.Background())
}
