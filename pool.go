package sql

import (
	"github.com/jackc/pgx"
	"github.com/monopolly/errors"
	"github.com/monopolly/weighted"
)

type Pool struct {
	write    *weighted.SW
	read     *weighted.SW
	list     *weighted.SW
	readonly bool
}

func New() (a *Pool) {
	a = new(Pool)
	a.write = new(weighted.SW)
	a.read = new(weighted.SW)
	a.list = new(weighted.SW)
	return
}

//Также добавляет и slave
func (a *Pool) Master(host, dbname, user, pass string, port int, maxconnect int, weight int) (err errors.E) {

	// init
	var conf pgx.ConnPoolConfig
	conf.Host = host
	conf.Database = dbname
	conf.Port = uint16(port)
	conf.User = user
	conf.Password = pass
	conf.MaxConnections = maxconnect

	// conn
	p, er := pgx.NewConnPool(conf)
	if er != nil {
		err = errors.Connection(er.Error())
		return
	}

	a.write.Add(p, weight)
	a.read.Add(p, weight)
	a.list.Add(p, 1)
	return
}

//Также добавляет и slave
func (a *Pool) Slave(host, dbname, user, pass string, port int, maxconnect int, weight int) (err errors.E) {

	// init
	var conf pgx.ConnPoolConfig
	conf.Host = host
	conf.Database = dbname
	conf.Port = uint16(port)
	conf.User = user
	conf.Password = pass
	conf.MaxConnections = maxconnect

	// conn
	p, er := pgx.NewConnPool(conf)
	if er != nil {
		err = errors.Connection(er.Error())
		return
	}

	a.read.Add(p, weight)
	a.list.Add(p, 1)
	return
}

//нужно делать c.Release() после операций
func (a *Pool) Write() (c *Conn, err errors.E) {
	//если readonly значит что мы закрываем все коннекты
	if a.readonly {
		err = errors.Try("readonly")
		return
	}
	pool := a.write.Next()
	if pool == nil {
		err = errors.Connection()
		return
	}

	c = new(Conn)
	p, ok := pool.(*pgx.ConnPool)
	if !ok {
		err = errors.Connection()
		return
	}
	c.pool = p
	return
}

//нужно делать c.Release() после операций
func (a *Pool) Read() (c *Conn, err errors.E) {
	//если readonly значит что мы закрываем все коннекты
	if a.readonly {
		err = errors.Try("readonly")
		return
	}
	pool := a.read.Next()
	if pool == nil {
		err = errors.Connection()
		return
	}
	c = new(Conn)
	p, ok := pool.(*pgx.ConnPool)
	if !ok {
		err = errors.Connection()
		return
	}
	c.pool = p
	return
}

//закрывает и удаляет все коннекты
func (a *Pool) Stop() {
	a.readonly = true

	a.list.Iterate(func(pool interface{}) {
		if pool == nil {
			return
		}
		pool.(*pgx.ConnPool).Close()
	})

	a.read.RemoveAll()
	a.write.RemoveAll()
}
