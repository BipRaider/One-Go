package main

import (
	"html/template" //https://golang.org/pkg/text/template/
	"os"
)

type User struct {
	Name string
	Age  int
	ANme string
}

func main() {
	t, err := template.ParseFiles("exp/hello.gohtml") // дайёт адрес най файл   и определяет по типу имени в вочто оборачиваться будет файл

	if err != nil {
		panic(err)
	}

	data := User{
		Name: "Jom Smith",
		Age:  355,
		ANme: "ASDasfsaf",
	}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)

	}

}
