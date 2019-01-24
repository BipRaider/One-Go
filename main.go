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
)
var NotF *views.NotView

func home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.ExecuteTemplate(w, homeView.Layout, nil)

	if err != nil {
		os.Exit(1)

	}
}

func contact(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	err := conatctView.Template.ExecuteTemplate(w, conatctView.Layout, nil)

	if err != nil {
		os.Exit(2)

	}
}

func notFound(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	err := NotF.Template.ExecuteTemplate(w, NotF.Layout, nil)

	if err != nil {
		os.Exit(1)

	}

}
func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	conatctView = views.NewView("bootstrap", "views/contact.gohtml")
	NotF = views.NotFound()
	r := mux.NewRouter() //https://www.gorillatoolkit.org/pkg/mux

	r.NotFoundHandler = http.HandlerFunc(notFound) //Заменили вид выводящейся ошибки на своё -----1.4
	r.HandleFunc("/home", home)
	r.HandleFunc("/contact", contact)

	http.ListenAndServe(":3000", r) // это адрес сервера  куда будет отправляться данные
}

//http://localhost:3000/
//https://getbootstrap.com/docs/3.3/components/#nav
