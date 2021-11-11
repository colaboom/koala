package main

import (
	"fmt"
	"os"
	"text/template"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Printf("parse file err : %v\n", err)
		return
	}

	p := Person{
		Name: "Tom",
		Age:  27,
	}
	err = t.Execute(os.Stdout, p)
	if err != nil {
		fmt.Printf("execute err : %v\n", err)
		return
	}

	return
}
