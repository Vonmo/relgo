/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */

package models

import (
	"strings"
)

type Counter struct {
	Name    string `db:"name"`
	Value   int    `db:"value"`
	Updated string `db:"updated"`
}

func (c *Counter) Get() (err error) {
	sql := `SELECT ` + strings.Join(TFields(&Counter{}, []string{}), ", ") + `
		    FROM counters WHERE name=$1`
	return db.Get(c, sql, c.Name)
}

func (c *Counter) Increment() (err error) {
	sql := `INSERT INTO counters (name, value)
			VALUES (:name, 1) 
			ON CONFLICT (name) 
				DO UPDATE SET value = counters.value + 1,
  							  updated = now()`
	_, err = db.NamedExec(sql, c)
	return err
}

func (c *Counter) Decrement() (err error) {
	sql := `INSERT INTO counters (name, value)
			VALUES (:name, -1) 
			ON CONFLICT (name) 
				DO UPDATE SET value = counters.value - 1,
						      updated = now()`
	_, err = db.NamedExec(sql, c)
	return err
}

func (c *Counter) Reset() (err error) {
	sql := "DELETE FROM counters WHERE name = :name"
	_, err = db.NamedExec(sql, c)
	return err
}
