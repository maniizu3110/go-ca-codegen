package querybuilder_test

import (
	"go-ca-codegen/codegen/template/api/repository/util/querybuilder"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEmbedEmbedUser struct {
	EmbedEmbedName   string
	EmbedEmbedStruct *testEmbedUser
}

type testEmbedUser struct {
	EmbedName string
	testEmbedEmbedUser
	EmbedStruct *testEmbedUser
}
type testStruct struct {
	String string
}
type testUser struct {
	Name         string
	Int          int
	Struct       *testStruct
	Struct2      testStruct
	UnExpandable testStruct `api:"unexpandable"`
	testEmbedUser
}

func Test_parseQueryQuery(t *testing.T) {
	s := func(args ...interface{}) []interface{} {
		res := make([]interface{}, len(args))
		for i := range args {
			res[i] = args[i]
		}
		return res
	}
	ss := func(args ...[]interface{}) [][]interface{} {
		var res [][]interface{}
		for i := range args {
			res = append(res, args[i])
		}
		return res
	}
	for _, user := range []interface{}{testUser{}, &testUser{}, []testUser{}, []*testUser{}} {
		queries := [][]string{
			{"Name=foo"},
			{"EmbedName=foo", "EmbedEmbedName=bar"},
			{"Name<foo", "Name<=foo", "Name>foo", "Name>=foo", "Name!=foo"},
			{"Name = foo", "Name    <=   foo"},
			{"Name in foo,bar,hoge", "Name in foo, bar, hoge"},
			{"Name not in foo,bar,hoge"},
			{"Name is NULL", "Name is NOTNULL"},
			{"Name includes foo"},
		}
		expected := [][][]interface{}{
			ss(s("name = ?", "foo")),
			ss(s("embed_name = ?", "foo"), s("embed_embed_name = ?", "bar")),
			ss(s("name < ?", "foo"), s("name <= ?", "foo"), s("name > ?", "foo"), s("name >= ?", "foo"), s("name != ?", "foo")),
			ss(s("name = ?", "foo"), s("name <= ?", "foo")),
			ss(s("name IN (?)", []string{"foo", "bar", "hoge"}), s("name IN (?)", []string{"foo", "bar", "hoge"})),
			ss(s("name NOT IN (?)", []string{"foo", "bar", "hoge"})),
			ss(s("name is NULL"), s("name is NOT NULL")),
			ss(s("name LIKE ?", "%foo%")),
		}
		for i, query := range queries {
			value, err := querybuilder.ParseQueryQuery(user, query)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, expected[i], value)
		}
		badQueries := [][]string{
			{"Struct=foo", "Struct2=goo"},
			{"Name includes foo,bar"},
		}
		for _, query := range badQueries {
			_, err := querybuilder.ParseQueryQuery(user, query)
			if !assert.Error(t, err) {
				return
			}
		}
	}
}
