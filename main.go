package main

import (
	"fmt"
	"net/http"
)

var ggg int = 333

func habdlerFunc(w http.ResponseWriter, r *http.Request) { //Данную функцию можно назвать как угодно
	print("1")
	w.Header().Set("Content-Type", "text/html") //https://golang.org/pkg/net/http/#Handler //Указывает тип страницы и стиль
	if r.URL.Path == "/" {                      //https://golang.org/pkg/net/url/#URL //r.URL.Path   имя страницы и что на ней будет
		fmt.Fprint(w, "<h1>Welcom to my as2 site!</h1>")
		fmt.Fprint(w, "<a href=\"contact\">Open contact</a>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "<a href=\"thebipus@gmail.com\">theBipus@gmail.com</a>")
		fmt.Fprint(w, "\n", ggg)
	} else {
		w.WriteHeader(http.StatusNotFound) // https://golang.org/pkg/net/http/#ResponseWriter  указывает состояние страницы
		fmt.Fprint(w, "<h1> We could not the page you were loking for	</h1><p>Please emaul us if you keep being sent to am invalide page.</p>")
	}

}

func main() {
	http.HandleFunc("/", habdlerFunc) // спомощью этой функций отправляется выше созданая функция адрес
	//"/"  открываться будет отсюда  и все следующие страницы http://localhost:3000
	//"/dog/" Строгая ,Страница будет открываться http://localhost:3000/dog/

	http.ListenAndServe(":3000", nil) // это адрес сервера  куда будет отправляться данные
}

//git remote add origin https://github.com/BipRaider/One-Go.git
//git push -u origin master
