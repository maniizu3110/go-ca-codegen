package main

import (
	"flag"
	"fmt"
	"go-ca-codegen/util"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

func GeneratePackage(destDir string, placeHolder string, model string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	getRel := func(f string) string {
		res, err := filepath.Rel(wd, f)
		if err != nil {
			panic(err)
		}
		return res
	}
	templateDir, err := filepath.Abs("../codegen/template")
	if err != nil {
		panic(err)
	}
	destDir, err = filepath.Abs(destDir)
	if err != nil {
		panic(err)
	}

	for _, file := range util.FindFiles(templateDir, ".go") {
		dest := filepath.Join(
			strings.Replace(filepath.Dir(file), templateDir, destDir, 1),
			strings.Replace(filepath.Base(file), placeHolder, strcase.ToSnake(model), 1),
		)
		prefix := util.SamePrefix(templateDir, dest)
		//prefix = prefix[:len(prefix)-1]
		suffix := util.SameSuffix(filepath.Dir(file), filepath.Dir(dest))[1:]

		fromImportPath := strings.TrimSuffix(strings.TrimPrefix(templateDir, prefix), suffix)
		toImportPath := strings.TrimSuffix(strings.TrimPrefix(filepath.Dir(dest), prefix), suffix)
		if len(toImportPath) == 0 {
			fromImportPath += "/"
		}

		isAutogen := strings.HasPrefix(filepath.Base(file), "autogen_")

		if strings.HasPrefix(filepath.Base(file), "ignore_") {
			// ignore _*
			continue
		} else if strings.Contains(filepath.Base(file), placeHolder) {
			if isAutogen {
				// 				headerMessage := `
				// // Code generated automatically. DO NOT EDIT.
				// // このファイルはプログラムによって自動生成されています
				// // 編集したい場合は，元のファイル {} を編集してください
				// `
				// 				headerMessage = strings.Replace(headerMessage, "{}", getRel(file), 1)
				headerMessage := ""
				funcMessage := `
この関数はプログラムによって自動生成されています

編集したい場合は，元のファイル {} を編集してください`
				funcMessage = strings.Replace(funcMessage, "{}", getRel(file), 1)
				util.CopyFileWithReplacePlaceHolder(
					file, dest,
					placeHolder, model,
					fromImportPath, toImportPath,
					headerMessage,
					funcMessage)
			} else if _, err := os.Stat(dest); os.IsNotExist(err) {
				// dest が存在しないなら作る
				// ただし，dest には prefix autogen_ をつける
				// 				headerMessage := `
				// // Code generated automatically. DO NOT EDIT.
				// // この関数はプログラムによって自動生成されています
				// // すべてのモデルに対する変更をしたい場合は，元のファイル {} を編集してください
				// // このモデル専用のカスタマイズしたい場合は，{} という名前のファイルを作って，以下の内容をコピーしてください
				// // （次回自動生成時に，このファイルが削除されます）
				// `
				// headerMessage = strings.Replace(headerMessage, "{}", getRel(file), 1)
				// headerMessage = strings.Replace(headerMessage, "{}", filepath.Base(dest), 1)
				headerMessage := ""
				funcMessage := `
この関数はプログラムによって自動生成されています

編集したい場合は，元のファイル {} を編集してください`
				funcMessage = strings.Replace(funcMessage, "{}", getRel(file), 1)
				dest2 := filepath.Join(filepath.Dir(dest), "autogen_"+filepath.Base(dest))
				util.CopyFileWithReplacePlaceHolder(
					file, dest2,
					placeHolder, model,
					fromImportPath, toImportPath,
					headerMessage, funcMessage)
			} else {
				// dest が存在するので作らない
				// autogen_ がついた dest があれば消す
				dest2 := filepath.Join(filepath.Dir(dest), "autogen_"+filepath.Base(dest))
				if _, err := os.Stat(dest2); !os.IsNotExist(err) {
					fmt.Println("[Delete]", dest2)
					if err := os.Remove(dest2); err != nil {
						panic(err)
					}
				}
			}
		} else {
			if _, err := os.Stat(dest); os.IsNotExist(err) || isAutogen {
				headerMessage := `
//


// Code generated automatically. DO NOT EDIT.
//
// このファイルはプログラムによって自動生成されています
//
// 編集したい場合は，元のファイル {} を編集してください
`
				headerMessage = strings.Replace(headerMessage, "{}", getRel(file), 1)
				util.CopyFile(file, dest, headerMessage)
			}
		}
	}
}

func main() {
	fmt.Println("hello")
	modelName := ""
	file := flag.String("file", "", "file name")
	dest := flag.String("dest", "", "destination root direction path")
	flag.Parse()
	if len(*file) > 0 {
		fileName := filepath.Base(*file)
		modelName = strcase.ToCamel(strings.TrimSuffix(fileName, filepath.Ext(fileName)))
	} else {
		panic("specify model or file")
	}
	GeneratePackage(*dest, "PlaceHolder", modelName)
}
