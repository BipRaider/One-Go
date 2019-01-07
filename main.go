package main

import (
	"fmt"
	"net/http"
)

func habdlerFunc(w http.ResponseWriter, r *http.Request) { //Данную функцию можно назвать как угодно
	// print("Сhac brauzer ") // log_color_app:     bold_white    при дабовление это в runner.conf   пичает в терменале какждый раз как обновит в браузере страницу
	fmt.Fprint(w, "<h1>Welcom to my as2 site!</h1>")
}

func main() {
	http.HandleFunc("/", habdlerFunc) // спомощью этой функций отправляется выше созданая функция адрес
	//"/"  открываться будет отсюда  и все следующие страницы http://localhost:3000

	//"/dog/" Строгая ,Страница будет открываться после http://localhost:3000/dog/

	http.ListenAndServe(":3000", nil) // это адрес сервера  куда будет отправляться данные
}

//git remote add origin https://github.com/BipRaider/One-Go.git
//git push -u origin master
