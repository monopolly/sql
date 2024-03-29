package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/monopolly/errors"
)

// insert, ignore = no conflicts on insert
func (a *Conn) Insert(table string, keys map[string]any, ignore ...bool) (id any, err errors.E) {
	/* insert into channels (sid,lang,category,sub,title) values ('en.startups.preseed', 'en','startups','preseed','Pre-Seed'); */
	lens := len(keys)
	klist := make([]string, lens)
	qlist := make([]string, lens)
	vlist := make([]any, lens)
	var c int
	for k, v := range keys {
		klist[c] = k
		vlist[c] = v
		qlist[c] = fmt.Sprintf("$%d", c+1)
		c++
	}
	q := fmt.Sprintf("insert into %s (%s) values (%s)", table, strings.Join(klist, ","), strings.Join(qlist, ","))
	if ignore != nil {
		q = q + " on conflict do nothing"
	}
	q += " returning id"
	pp := a.Pool.QueryRow(context.Background(), q, vlist...)
	er := pp.Scan(&id)
	if er != nil {
		if er == pgx.ErrNoRows {
			return
		}
		err = errors.Database(er)
	}
	return
}

// insert, ignore = no conflicts on insert
func (a *Conn) InsertTuple(sql string, v ...any) (id any, err errors.E) {
	/* insert into channels (sid,lang,category,sub,title) values ('en.startups.preseed', 'en','startups','preseed','Pre-Seed'); */
	sql += " returning id"
	pp := a.Pool.QueryRow(context.Background(), sql, v...)
	er := pp.Scan(&id)
	if er != nil {
		if er == pgx.ErrNoRows {
			return
		}
		err = errors.Database(er)
	}
	return
}
