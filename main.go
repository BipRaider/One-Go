package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"./views"

	"github.com/gorilla/mux"
)

var (
	homeTeplate    *template.Template
	conatctTeplate *template.Template
)

func home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	if err := homeTeplate.Execute(w, nil); err != nil {
		os.Exit(1)
		panic(err)

	}
}

func contact(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	if err := conatctTeplate.Execute(w, nil); err != nil {
		os.Exit(4)
		panic(err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1> NotFoud file</h1>")
}
func main() {
	homeView = views.NewView("views/home.gohtml", "views/home_banner.gohtml")

	var err error
	homeTeplate, err = template.ParseFiles(
		"views/home.gohtml",
		"views/footer/footer.gohtml",
	)
	if err != nil {
		os.Exit(2)
		panic(err)
	}

	conatctTeplate, err = template.ParseFiles(
		"views/contact.gohtml",
		"views/footer/footer.gohtml",
	)
	if err != nil {
		os.Exit(3)
		panic(err)
	}

	r := mux.NewRouter() //https://www.gorillatoolkit.org/pkg/mux

	r.NotFoundHandler = http.HandlerFunc(notFound) //Заменили вид выводящейся ошибки на своё -----1.4

	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)

	http.ListenAndServe(":3000", r) // это адрес сервера  куда будет отправляться данные
}
