package main

import (
	"html/template" //https://golang.org/pkg/text/template/
	"os"
)

type Dof struct {
	Name1 string
}
type User struct { // создаём тип  и  указываем имена переменных для определения черер t.Execute(os.Stdout,data)
	Dof  Dof
	Name string
	Age  int

	Float  float64
	Slise  []string
	Slise1 []int
	Map    map[string]string
}

func main() {
	//https://golang.org/pkg/text/template/
	//https://golang.org/pkg/html/template/#ParseFiles
	t, err := template.ParseFiles("exp/hello.gohtml") // дайёт адрес най файл   и определяет по типу имени в вочто оборачиваться будет файл ,функция ParseFiles  для анализа нашего файла шаблона

	if err != nil {
		panic(err)
	}

	data := User{ //  обьевляем переменную с сзначениями
		Dof: Dof{
			Name1: "pikul",
		},

		Name:  "Jom Smith",
		Age:   355,
		Float: 55.34,
		Slise: []string{"A", "b", "sdf"},

		Map: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": " value3",
		},
	}
	err = t.Execute(os.Stdout, data) // обрабатываем значения  переменой data  и по именам присваем им вид отображения (в файле hello.gohtml)
	if err != nil {
		panic(err)

	}

}
