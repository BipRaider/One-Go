package main

import (
	"fmt"
	"net/http"
	"os"

	"./controllers"
	"./middleware"
	"./models"
	"./views"
	"github.com/gorilla/mux"
)

var NotF *views.View

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
	usersC := controllers.NewUser(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, r)
	///-----------------
	requireUserMw := middleware.RequireUser{
		UserService: services.User, // запрос из БД, данные на пользователя
	}
	///------------
	NotF = views.NotFound()                        //2
	r.NotFoundHandler = http.HandlerFunc(notFound) //3 //Заменили вид выводящейся ошибки на своё
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/home", staticC.Home).Methods("GET")       //3
	r.Handle("/contact", staticC.Contact).Methods("GET") //3
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")     //3
	r.HandleFunc("/signup", usersC.Create).Methods("POST") //3//Выводит сообщение от функций Create
	r.HandleFunc("/faq", usersC.NewFaqGet).Methods("GET")
	r.HandleFunc("/faq", usersC.Create).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	//Gallery routes

	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

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
