package main

import (
	"fmt"
	"net/http"
	"os"

	"./controllers"
	"./models"
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

const (
	mysqlinfo = "root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local"
)

func main() {
	//соединение сайта с базойданых
	us, err := models.NewUserService(mysqlinfo)
	must(err, 3)
	defer us.Close()
	//us.DestructiveReset()
	us.AutoMigrate()

	r := mux.NewRouter() //1 begin

	// homeView = views.NewView("bootstrap", "views/home.gohtml")       //2
	// conatctView = views.NewView("bootstrap", "views/contact.gohtml") //2
	staticC := controllers.NewStatic()
	usersC := controllers.NewUser(us) //2
	// faqC := controllers.NewFAQ()

	//https://www.gorillatoolkit.org/pkg/mux

	NotF = views.NotFound()                        //2
	r.NotFoundHandler = http.HandlerFunc(notFound) //3 //Заменили вид выводящейся ошибки на своё

	r.Handle("/home", staticC.Home).Methods("GET")       //3
	r.Handle("/contact", staticC.Contact).Methods("GET") //3

	r.HandleFunc("/faq", usersC.NewFaqGet).Methods("GET")
	r.HandleFunc("/faq", usersC.Create).Methods("POST")

	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/signup", usersC.New).Methods("GET")     //3
	r.HandleFunc("/signup", usersC.Create).Methods("POST") //3//Выводит сообщение от функций Create

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	fmt.Println("Starting the server on :3000...")
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

//https://dev.mysql.com/doc/workbench/en/wb-mysql-connections-navigator-management-users-and-privileges.html

//https://tproger.ru/translations/go-web-server/amp/ получени сертифеката https
