/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */

package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
)

type NString struct {
	sql.NullString
}

// -----------------------------------------------------------------------------

func ToNString(s string) NString {
	return NString{sql.NullString{String: s, Valid: s != ""}}
}

func TFields(structure interface{}, exclude []string) []string {
	res := []string{}
	s := reflect.ValueOf(structure).Elem()
	for i := 0; i < s.NumField(); i++ {
		name := s.Type().Field(i).Name
		tag := s.Type().Field(i).Tag.Get("db")
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = strings.ToLower(name)
		}

		excluded := false
		for _, a := range exclude {
			if a == tag {
				excluded = true
				break
			}
		}
		if excluded == false {
			res = append(res, tag)
		}
	}
	return res
}

func TFieldsT(table string, structure interface{}, exclude []string) []string {
	res := []string{}
	s := reflect.ValueOf(structure).Elem()
	for i := 0; i < s.NumField(); i++ {
		name := s.Type().Field(i).Name
		tag := s.Type().Field(i).Tag.Get("db")
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = strings.ToLower(name)
		}

		excluded := false
		for _, a := range exclude {
			if a == tag {
				excluded = true
				break
			}
		}
		if excluded == false {
			res = append(res, table+"."+tag)
		}
	}
	return res
}

func TInsertSql(table string, fields []string, returning string) string {
	sql := "INSERT INTO " + table + " (" + strings.Join(fields, ", ") + ") "
	sql += "VALUES (" + strings.Join(Map(fields, func(s string) string { return ":" + s }), ", ") + ") "
	if returning != "" {
		sql += "RETURNING " + returning
	}
	return sql + ";"
}

func TUpdateSql(table string, fields []string) string {
	sql := "UPDATE " + table + " SET updated = now(), "
	cnt := len(fields)
	i := 0
	for _, f := range fields {
		if f != "id" {
			sql += f + "=:" + f
			if i != cnt-1 {
				sql += ", "
			}
		}
		i++
	}
	sql += " WHERE id = :id"
	return sql + ";"
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// -----------------------------------------------------------------------------

func (ns *NString) MarshalJSON() ([]byte, error) {
	if ns.String == "" && !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NString) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		ns.String = ""
		ns.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ns.String)
	if err != nil {
		return err
	}
	ns.Valid = true
	return nil
}
