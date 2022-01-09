package querybuilder

import (
	"reflect"
	"testing"
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

func Test_parseExpandQuery(t *testing.T) {
	check := func(i interface{}) bool {
		res := parseExpandQuery(i, []string{"Name", "Struct", "Struct2", "UnExpandable", "testEmbedUser", "EmbedStruct", "EmbedEmbedStruct"})
		return reflect.DeepEqual(res, []string{"Struct", "Struct2", "EmbedStruct", "EmbedEmbedStruct"})
	}
	var user testUser
	if !check(user) {
		t.Error("user")
	}
	if !check(&user) {
		t.Error("&user")
	}
	users := []testUser{user}
	if !check(users) {
		t.Error("users")
	}
	if !check(&users) {
		t.Error("&users")
	}
}
