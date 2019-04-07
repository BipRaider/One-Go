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
	NotF *views.View
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	NotF.Render(w, nil)

}

const (
	mysqlinfo = "root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local"
)

func main() {
	//соединение сайта с базойданых
	services, err := models.NewServices(mysqlinfo)
	must(err, 3)

	defer services.Close()
	//services.DestructiveReset() // удаляет из бд
	services.AutoMigrate() //записует в бд

	r := mux.NewRouter() //1 begin

	staticC := controllers.NewStatic()
	usersC := controllers.NewUser(services.User) //2

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

//https://www.gorillatoolkit.org/pkg/mux
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
