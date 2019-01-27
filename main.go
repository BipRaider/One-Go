package main

import (
	"net/http"
	"os"

	"./views"
	"github.com/gorilla/mux"
)

var (
	homeView    *views.View
	conatctView *views.View
	signupView  *views.View
	NotF        *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil), 1)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(conatctView.Render(w, nil), 2)
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signupView.Render(w, nil), 3)

}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	must(NotF.Render(w, nil), 404)

}
func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	conatctView = views.NewView("bootstrap", "views/contact.gohtml")
	signupView = views.NewView("bootstrap", "views/signup.gohtml")

	r := mux.NewRouter() //https://www.gorillatoolkit.org/pkg/mux

	NotF = views.NotFound()
	r.NotFoundHandler = http.HandlerFunc(notFound) //Заменили вид выводящейся ошибки на своё -----1.4
	r.HandleFunc("/home", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/signup", signup)

	http.ListenAndServe(":3000", r) // это адрес сервера  куда будет отправляться данные
}

//http://localhost:3000/
//https://getbootstrap.com/docs/3.3/components/#nav

//Функция вывода ошибоки
func must(err error, n int) {
	if err != nil {
		os.Exit(n)
	}

}
