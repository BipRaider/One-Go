package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var ggg int = 333

// func habdlerFunc(w http.ResponseWriter, r *http.Request) { //Данную функцию можно назвать как угодно
// 	print("1")
// 	w.Header().Set("Content-Type", "text/html")        //https://golang.org/pkg/net/http/#Handler //Указывает тип страницы и стиль (Html amd CSS  и тд)
// 	if r.URL.Path == "/wqe" || r.URL.Path == "/wqe/" { //https://golang.org/pkg/net/url/#URL //r.URL.Path   имя страницы и что на ней будет
// 		fmt.Fprint(w, "<h1>Welcom to my as2 site!</h1>")
// 		fmt.Fprint(w, "<a href=\"contact\">Open contact</a>")
// 	} else if r.URL.Path == "/contact" || r.URL.Path == "/contact/" {
// 		fmt.Fprint(w, "<a href=\"thebipus@gmail.com\">theBipus@gmail.com</a>")
// 		fmt.Fprint(w, "\n", ggg)
// 	} else {
// 		w.WriteHeader(http.StatusNotFound) // https://golang.org/pkg/net/http/#ResponseWriter  указывает состояние страницы
// 		fmt.Fprint(w, "<h1> We could not the page you were loking for	</h1><p>Please emaul us if you keep being sent to am invalide page.</p>")
// 	}

// }

//--------------------------1.3
//Прописываешь wr Tabl и вывордится это :w http.ResponseWriter, r *http.Request
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/wqe" || r.URL.Path == "/wqe/" { //https://golang.org/pkg/net/url/#URL //r.URL.Path   имя страницы и что на ней будет
		fmt.Fprint(w, "<h1>Welcom to my as2 site!</h1>")
		fmt.Fprint(w, "<a href=\"contact\">Open contact</a>")
	} else {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "<h1>Welcom to my HOME!</h1>")
	}
}
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<a href=\"thebipus@gmail.com\">theBipus@gmail.com</a>")

}

//---------------1.2
// func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	fmt.Fprint(w, "Welcome!\n")
// }

// func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
// }

func main() {
	//---------------1.3
	r := mux.NewRouter() //https://www.gorillatoolkit.org/pkg/mux
	r.HandleFunc("/wqe/", home)
	r.HandleFunc("/contact/", contact)
	//----------------------1.2
	// router := httprouter.New() //"github.com/julienschmidt/httprouter"
	// router.GET("/", Index)
	// router.GET("/hello/:name", Hello)

	//  log.Fatal(http.ListenAndServe(":8080", router))
	//-------------------------------1.1
	// mux := &http.ServeMux{} //https://golang.org/pkg/net/http/#ServeMux
	// mux.HandleFunc("/", habdlerFunc)

	// http.HandleFunc("/", habdlerFunc) // спомощью этой функций отправляется выше созданая функция адрес
	//"/"  открываться будет отсюда  и все следующие страницы http://localhost:3000
	//"/dog/" Строгая ,Страница будет открываться http://localhost:3000/dog/

	http.ListenAndServe(":3000", r) // это адрес сервера  куда будет отправляться данные
}

//git remote add origin https://github.com/BipRaider/One-Go.git
//git push -u origin master
