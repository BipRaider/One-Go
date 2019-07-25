package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"./controllers"
	"./middleware"
	"./models"
	"./rand"
	"./views"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

var NotF *views.View

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	NotF.Render(w, r, nil)

}

func main() {
	boolPtr := flag.Bool("Service", false,
		"Provide this flag in production.This ensures"+
			"that a .config file is provided the application starts.") //переменая -определяем запущены ли  файлa (в нашем случае Server)
	flag.Parse()                // возвращаем  определёные файлы при запуске
	cfg := LoadConfig(*boolPtr) // переменная запуск конфига  который бреобразуется из докумета и обрабатуеться  JSON методом
	dbCfg := cfg.Database       // переменая с данными о сапуске сервера БД

	//соединение с базойданых
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInf()),
		models.WithLogMode(!cfg.isProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	must(err)

	defer services.Close()
	//services.DestructiveReset() // удаляет из бд
	services.AutoMigrate() //записует в бд

	r := mux.NewRouter() //1 begin

	staticC := controllers.NewStatic()
	usersC := controllers.NewUser(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	///----------------- Защита сайта от поделлки ,копирования

	b, err := rand.Bytes(32) // кодируем даные
	must(err)
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.isProd())) // кодируем страницу
	// Protect - это промежуточное ПО HTTP, которое обеспечивает защиту CSRF на маршрутах, подключенных к маршрутизатору или суб-маршрутизатору.
	// Secure -устанавливает флаг безопасности в куки. По умолчанию true // Установите  «false» в противном случае файл cookie не будет отправляться по небезопасному каналу
	///-----------------
	userMw := middleware.User{
		UserService: services.User, // запрос из БД, данные на пользователя
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}
	///----------------
	NotF = views.NotFound()                        //2
	r.NotFoundHandler = http.HandlerFunc(notFound) //3 //Заменили вид выводящейся ошибки на своё
	//------
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/home", staticC.Home).Methods("GET")       //3
	r.Handle("/contact", staticC.Contact).Methods("GET") //3
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")     //3
	r.HandleFunc("/signup", usersC.Create).Methods("POST") //3//Выводит сообщение от функций Create
	//Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))    //путь папки в системе
	assetHandler = http.StripPrefix("/assets/", assetHandler) // удаляя указанный префикс из пути URL запроса и вызывая обработчик
	r.PathPrefix("/assets/").Handler(assetHandler)            // Для вывода и обработки урл адресов в браузер
	//Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))    //путь папки в системе
	imageHandler = http.StripPrefix("/images/", imageHandler) // удаляя указанный префикс из пути URL запроса и вызывая обработчик
	r.PathPrefix("/images/").Handler(imageHandler)            // Для вывода и обработки урл адресов в браузер

	//Gallery routes
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	//_POST/galleries/:id/images/:filename/delete
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r))) // это адрес сервера  куда будет отправляться данные и закодированый от поделки

}

//Функция вывода ошибоки
func must(err error) {
	if err != nil {
		log.Println(err)
	}

}

//https://www.gorillatoolkit.org/pkg/mux
//https://getbootstrap.com/docs/3.3/components/#nav

//go build . && ./Go -Service  сохранить и запустить сервер

//https://dev.mysql.com/doc/workbench/en/wb-mysql-connections-navigator-management-users-and-privileges.html

//https://tproger.ru/translations/go-web-server/amp/ получени сертифеката https
