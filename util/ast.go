package util

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

// ReplaceImportPath は各 import path に含まれる <from> を <to> に置き換える
func ReplaceImportPath(f *ast.File, from string, to string) {
	for i := range f.Imports {
		f.Imports[i].Path.Value = strings.Replace(f.Imports[i].Path.Value, from, to, 1)
	}
}

// ReplaceIdent は f の変数名，コメントに含まれる <from> をすべて <to> に置き換える
// <from> の Camel, LowerCamel, Snake Case が to の対応する Case に置き換えられる
func ReplaceIdent(f *ast.File, from string, to string) {
	froms := []string{strcase.ToCamel(from), strcase.ToLowerCamel(from), strcase.ToSnake(from)}
	tos := []string{strcase.ToCamel(to), strcase.ToLowerCamel(to), strcase.ToSnake(to)}
	ast.Inspect(f, func(n ast.Node) bool {
		switch aType := n.(type) {
		case *ast.Ident: // 変数名を置き換える
			for i := range froms {
				aType.Name = strings.Replace(aType.Name, froms[i], tos[i], 1)
			}
		case *ast.Comment: // コメントの中に出てくる文字列を置き換える
			for i := range froms {
				aType.Text = strings.ReplaceAll(aType.Text, froms[i], tos[i])
			}
		}
		return true
	})
}

func CopyFileWithReplacePlaceHolder(
	templateFile string, destFile string,
	placeHolder string, model string,
	fromImportPath string, toImportPath string,
	headerMessage string,
) {
	destDir := filepath.Dir(destFile)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		os.MkdirAll(destDir, 0777)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, templateFile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ReplaceImportPath(f, fromImportPath, toImportPath)
	ReplaceIdent(f, placeHolder, model)

	file, err := os.Create(destFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(headerMessage)
	if err := format.Node(file, fset, f); err != nil {
		panic(err)
	}
}
