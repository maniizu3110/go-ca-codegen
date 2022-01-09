package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/maniizu3110/go-ca-codegen/util"
)

func GeneratePackage(destDir string, placeHolder string, model string) {
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

		if strings.HasPrefix(filepath.Base(file), "ignore_") {
			// ignore _*
			continue
		} else if strings.Contains(filepath.Base(file), placeHolder) {
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				// dest が存在しないなら作る
				headerMessage := ""
				dest2 := filepath.Join(filepath.Dir(dest), filepath.Base(dest))
				util.CopyFileWithReplacePlaceHolder(
					file, dest2,
					placeHolder, model,
					fromImportPath, toImportPath,
					headerMessage)
			} else {
				// dest が存在するときは作成しない
				dest2 := filepath.Join(filepath.Dir(dest), filepath.Base(dest))
				if _, err := os.Stat(dest2); !os.IsNotExist(err) {
					fmt.Println("[Delete]", dest2)
					if err := os.Remove(dest2); err != nil {
						panic(err)
					}
				}
			}
		} else {
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				util.CopyFile(file, dest)
			}
		}
	}
}

func main() {
	modelName := ""
	file := flag.String("file", "", "file name")
	dest := flag.String("dest", "", "destination root direction path")
	flag.Parse()
	if len(*file) > 0 {
		fileName := filepath.Base(*file)
		modelName = strcase.ToCamel(strings.TrimSuffix(fileName, filepath.Ext(fileName)))
	} else {
		panic("please specify file")
	}
	GeneratePackage(*dest, "PlaceHolder", modelName)
}
