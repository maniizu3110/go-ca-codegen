package querybuilder_test

import (
	"reflect"
	"testing"

	"github.com/maniizu3110/go-ca-codegen/codegen/template/api/repository/util/querybuilder"
	"github.com/stretchr/testify/assert"
)

func TestParseOrder(t *testing.T) {
	type A struct {
		Name string
	}
	type B struct {
		A1 string
		A2 int
		N1 A
		N2 *A
		N3 []A
	}
	var b B
	ok := []string{"A1", "-A2"}
	okdb := []string{"a1", "a2"}
	dir := []string{"asc", "desc"}

	for i, member := range ok {
		res, err := querybuilder.ParseOrder(reflect.TypeOf(b), member)
		assert.NoError(t, err)
		assert.Equal(t, okdb[i]+" "+dir[i], res)
	}
	for _, member := range []string{"N1", "N2", "N3", "a2"} {
		res, err := querybuilder.ParseOrder(reflect.TypeOf(b), member)
		assert.Error(t, err)
		assert.Equal(t, "", res)
	}
}

func TestParseOrderQuery(t *testing.T) {
	type B struct {
		A string
		B []string
	}
	var a B
	_, err := querybuilder.ParseOrderQuery(a, []string{"A", "-A"})
	assert.NoError(t, err)
	_, err = querybuilder.ParseOrderQuery(&a, []string{"A"})
	assert.NoError(t, err)
	_, err = querybuilder.ParseOrderQuery(a, []string{"B"})
	assert.Error(t, err)
}
