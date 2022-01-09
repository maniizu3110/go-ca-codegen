package querybuilder

import (
	"errors"
	"go-ca-codegen/codegen/template/api/repository/util"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func ParseOrderQuery(model interface{}, members []string) ([]string, error) {
	t := util.GetElementType(reflect.TypeOf(model))
	var res []string
	for i := range members {
		s, err := ParseOrder(t, members[i])
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func ParseOrder(t reflect.Type, order string) (string, error) {
	var key string
	var dir string
	if strings.HasPrefix(order, "-") {
		key = strings.TrimLeft(order, "-")
		dir = "desc"
	} else {
		key = strings.TrimLeft(order, " ")
		dir = "asc"
	}
	field, found := util.FindFieldByNameDeep(t, key)
	if found && isSupportedType(field.Type) {
		//gorm.io/gormには下記メソッドがないためToDBNameしない
		//return gorm.ToDBName(key) + " " + dir, nil
		return key + "" + dir, nil
	}
	return "", errors.New(key + " がデータに含まれていません")
}

// BuildOrderQuery
// member = ["Key", "-Key"], -で降順になる
func BuildOrderQuery(model interface{}, member []string, query *gorm.DB) (*gorm.DB, error) {
	s, err := ParseOrderQuery(model, member)
	if err != nil {
		return nil, err
	}
	for i := range s {
		query = query.Order(s[i])
	}
	return query, nil
}
