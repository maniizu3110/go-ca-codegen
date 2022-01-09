package querybuilder

import (
	"go-ca-codegen/codegen/template/api/repository/util"
	"reflect"

	"gorm.io/gorm"
)

// parseExpandQuery : <members> の要素のうち， <model> がメンバ変数 or 埋め込み構造体のメンバ変数として持つ要素のみを返す
//                    タグで api: unexpandable な要素は除外する
func parseExpandQuery(model interface{}, members []string) []string {
	t := util.GetElementTypeDeep(model)
	res := make([]string, 0)
	for i := range members {
		field, found := util.MidDeepFind(t, members[i])
		if found && !field.Anonymous {
			t2 := field.Type
			if t2.Kind() == reflect.Ptr {
				t2 = t2.Elem()
			}
			if t2.Kind() == reflect.Struct || t2.Kind() == reflect.Slice {
				expandable := !util.FindValueFromTag(field, "api", "unexpandable")
				if expandable {
					res = append(res, field.Name)
				}
			}
		}
	}
	return res
}

func BuildExpandQuery(model interface{}, expandMembers []string, query *gorm.DB, preloadConditions ...interface{}) (*gorm.DB, error) {
	members := parseExpandQuery(model, expandMembers)
	for _, m := range members {
		query = query.Preload(m, preloadConditions...)
	}
	return query, nil
}
