package main

import (
	"html/template"
)

type Book struct {
	Name  string
	Words []Word
}

type Word struct {
	Value     string
	Usage     template.HTML
	IsEnglish bool
	Language  string
	Frecuency int
}
