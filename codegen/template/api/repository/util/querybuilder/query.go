package querybuilder

import (
	"errors"
	"github.com/maniizu3110/go-ca-codegen/codegen/template/api/repository/util"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func TrySplitQuery(t reflect.Type, query string, op string) (key string, vals []string, ok bool) {
	split := strings.SplitN(query, op, 2)
	if len(split) != 2 {
		ok = false
		return
	}
	key = strings.TrimSpace(split[0])
	field, found := util.FindFieldByNameDeep(t, key)
	if found && !isSupportedType(field.Type) {
		ok = false
		return
	}
	for _, val := range strings.Split(split[1], ",") {
		vals = append(vals, strings.TrimSpace(val))
	}
	ok = len(vals) > 0
	return
}

func ParseQuery(t reflect.Type, query string) ([]interface{}, error) {
	if key, vals, ok := TrySplitQuery(t, query, " not in "); ok {
		return []interface{}{key + " NOT IN (?)", vals}, nil
	}
	if key, vals, ok := TrySplitQuery(t, query, " in "); ok {
		return []interface{}{key + " IN (?)", vals}, nil
	}
	if key, vals, ok := TrySplitQuery(t, query, " includes "); ok && len(vals) == 1 {
		return []interface{}{key + " LIKE ?", "%" + vals[0] + "%"}, nil
	}
	for _, op := range []string{"<=", ">=", "!=", "=", "<", ">"} {
		if key, vals, ok := TrySplitQuery(t, query, op); ok && len(vals) == 1 {
			return []interface{}{key + " " + op + " ?", vals[0]}, nil
		}
	}

	if key, vals, ok := TrySplitQuery(t, query, " is "); ok && len(vals) == 1 {
		switch vals[0] {
		case "NULL":
			return []interface{}{key + " is NULL"}, nil
		case "NOTNULL":
			return []interface{}{key + " is NOT NULL"}, nil
		}
	}
	return nil, errors.New(query + " は無効なクエリです")
}

func ParseQueryQuery(model interface{}, queryMember []string) ([][]interface{}, error) {
	var res [][]interface{}
	t := util.GetElementType(reflect.TypeOf(model))

	for i := range queryMember {
		q, err := ParseQuery(t, queryMember[i])
		if err != nil {
			return nil, err
		}
		res = append(res, q)
	}
	return res, nil
}

// BuildQueryQuery parses each queryMember and build SQL.
//
// Following queries are supported:
// 	operators : =, !=, <, <=, >, >=
//		ex) name=foo
//  in : in, not in, support multi values (separated by ',')
//		ex) name in a,b,c,d
//			name not in foo,bar
//  null check : is NULL is NOTNULL
//		ex) name is NULL
//			name is NOTNULL
//  like :
//      ex) name like foo%
func BuildQueryQuery(model interface{}, member []string, query *gorm.DB) (*gorm.DB, error) {
	queries, err := ParseQueryQuery(model, member)
	if err != nil {
		return nil, err
	}
	for i := range queries {
		query = query.Where(queries[i][0], queries[i][1:]...)
	}
	return query, nil
}
