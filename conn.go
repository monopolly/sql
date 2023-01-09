package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monopolly/errors"
)

type Conn struct {
	Pool *pgxpool.Pool
}

// default 4
func (a *Conn) SetMaxConnections(max int) {
	a.Pool.Config().MaxConns = int32(max)
}

func New(host, db, user, pass string, port ...int) (res *Conn, err errors.E) {
	ports := 5432
	if port != nil {
		ports = port[0]
	}
	return NewString(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, ports, db))
}

func NewLocal(user, pass string, db ...string) (res *Conn, err errors.E) {
	switch db == nil {
	case true:
		return NewString(fmt.Sprintf("postgres://%s:%s@localhost:5432", user, pass))
	default:
		return NewString(fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", user, pass, db[0]))
	}
}

// postgres://user:pass@localhost:5432/dbname
// postgres://user:pass@localhost:5432/dbname?application_name=pgxtest&search_path=myschema&connect_timeout=5
func NewString(connstring string) (res *Conn, err errors.E) {

	// conn
	p, er := pgxpool.New(context.Background(), connstring)
	if er != nil {
		err = errors.Connection(er)
		return
	}

	res = new(Conn)
	res.Pool = p
	res.Pool.Config().HealthCheckPeriod = time.Second * 30
	return
}
