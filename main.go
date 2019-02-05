package main

import (
	"net/http"
	"os"

	"./controllers"
	"./views"

	"github.com/gorilla/mux"
)

var (
	// homeView    *views.View
	// conatctView *views.View
	NotF *views.View
)

//-----2-----
// func home(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	must(homeView.Render(w, nil), 1)
// }

// func contact(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	must(conatctView.Render(w, nil), 2)
// }

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	must(NotF.Render(w, nil), 404)

}
func main() {
	r := mux.NewRouter() //1 begin

	// homeView = views.NewView("bootstrap", "views/home.gohtml")       //2
	// conatctView = views.NewView("bootstrap", "views/contact.gohtml") //2
	staticC := controllers.NewStatic()
	usersC := controllers.NewUser() //2
	faqC := controllers.NewFAQ()

	//https://www.gorillatoolkit.org/pkg/mux

	NotF = views.NotFound()                        //2
	r.NotFoundHandler = http.HandlerFunc(notFound) //3 //Заменили вид выводящейся ошибки на своё

	r.Handle("/home", staticC.Home).Methods("GET")       //3
	r.Handle("/contact", staticC.Contact).Methods("GET") //3
	r.HandleFunc("/faq", faqC.FAQ).Methods("GET")
	// r.HandleFunc("/faq", faqC.FAQ).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")     //3
	r.HandleFunc("/signup", usersC.Create).Methods("POST") //3//Выводит сообщение от функций Create

	http.ListenAndServe(":3000", r) //end// это адрес сервера  куда будет отправляться данные
}

//http://localhost:3000/
//https://getbootstrap.com/docs/3.3/components/#nav

//Функция вывода ошибоки
func must(err error, n int) {
	if err != nil {
		os.Exit(n)
	}

}
