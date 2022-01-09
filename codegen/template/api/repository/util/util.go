package util

import (
	"reflect"
	"strings"
)

const GetAllMaxLimit = 100000000
const GetAllSubLimit = 10000

func ShallowCopy(m interface{}) interface{} {
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	clone := reflect.New(t)
	clone.Elem().Set(reflect.Indirect(reflect.ValueOf(m)))
	return clone.Interface()
}

func FindFieldByNameDeep(t reflect.Type, name string) (reflect.StructField, bool) {
	f, found := t.FieldByName(name)
	if found {
		return f, true
	}
	for i := 0; i < t.NumField(); i++ {
		switch t.Field(i).Type.Kind() {
		case reflect.Struct:
			if t.Field(i).Anonymous {
				f, found = FindFieldByNameDeep(t.Field(i).Type, name)
				if found {
					return f, found
				}
			}
		}
	}
	return reflect.StructField{}, false
}

func GetElementType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func NewInstance(i interface{}) interface{} {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}

// NewInstance returns reference of new slice of i
func NewSliceOf(i interface{}) interface{} {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	ptr := reflect.New(slice.Type())
	ptr.Elem().Set(slice)
	return ptr.Interface()
}

// field が tag: value を持つかチェック
func FindValueFromTag(field reflect.StructField, tag string, value string) bool {
	tags := strings.Split(field.Tag.Get(tag), ",")
	for i := range tags {
		if strings.TrimSpace(tags[i]) == value {
			return true
		}
	}
	return false
}

// MidDeepFind 埋め込まれた構造体までフィールドを探す
func MidDeepFind(t reflect.Type, name string) (reflect.StructField, bool) {
	f, found := t.FieldByName(name)
	if found {
		return f, true
	}
	for i := 0; i < t.NumField(); i++ {
		switch t.Field(i).Type.Kind() {
		case reflect.Struct:
			if t.Field(i).Anonymous {
				f, found = MidDeepFind(t.Field(i).Type, name)
				if found {
					return f, found
				}
			}
		}
	}
	return reflect.StructField{}, false
}

// スライス, ポインタ, ポインタのスライス, スライスのポインタ,とかから本来の型を取り出す
func GetElementTypeDeep(model interface{}) reflect.Type {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// GetElementValue returns ValueOf(model) or ValueOf(model).Elem() if model is Ptr
func GetElementValue(model interface{}) reflect.Value {
	vs := reflect.ValueOf(model)
	if vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	return vs
}
