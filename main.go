package main

import (
	"fmt"
	"net/http"
	"os"

	"./views"
	"github.com/gorilla/mux"
)

var (
	homeView    *views.View
	conatctView *views.View
)

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
	fmt.Fprint(w, "<h1> NotFoud file</h1>")
}
func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	conatctView = views.NewView("bootstrap", "views/contact.gohtml")

	r := mux.NewRouter() //https://www.gorillatoolkit.org/pkg/mux

	r.NotFoundHandler = http.HandlerFunc(notFound) //Заменили вид выводящейся ошибки на своё -----1.4
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)

	http.ListenAndServe(":3000", r) // это адрес сервера  куда будет отправляться данные
}
