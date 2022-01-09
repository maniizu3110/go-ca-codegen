package main

import "gorm.io/gorm"

//go:generate go run ../codegen/main.go -file ${GOFILE} -dest ..

type Test struct {
	gorm.Model
}
