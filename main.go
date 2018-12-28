package main

import (
	"fmt"
	"net/http"
)

func habdlerFunc(w http.ResponseWriter, r *http.Request) { //Данную функцию можно назвать как угодно
	fmt.Fprint(w, "<h1>Welcom to my awesone site!</h1>")
}

func main() {
	http.HandleFunc("/", habdlerFunc) // спомощью этой функций отправляется выше созданая функция адрес
	http.ListenAndServe(":3000", nil) // это адрес сервера  куда будет отправляться данные
}

//git remote add origin https://github.com/BipRaider/One-Go.git
//git push -u origin master
