package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	CreateUsers = "insert into bipusdb.users (Age, FirstName, LastName, Email) values (?, ?, ?, ?)"
)

type UserDB struct {
	Id        int    `schema:"id"`
	Age       int    `schema:"Age"`
	FirstName string `schema:"FirstName"`
	LastName  string `schema:"LastName"`
	Email     string `schema:"Email"`
}

var database *sql.DB

// функция добавления данных
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		Age := r.FormValue("Age")
		FirstName := r.FormValue("FirstName")
		LastName := r.FormValue("LastName")
		Email := r.FormValue("Email")
		// Функция для записи данных в базу данных
		// insert into - NameDatabase.NameTable
		// (age,firstname...) -- имена стобцов в базе данных
		//values (?) -- создаётся псевдоним  , данные что надо записать
		_, err = database.Exec(CreateUsers, Age, FirstName, LastName, Email)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "create.html")
	}
}

// запрос в базу даных и вывод на экран
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from bipusdb.users") // запрос вывода из базы данных
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	products := []UserDB{} //

	for rows.Next() {
		p := UserDB{}                                                        // куда записываются данные
		err := rows.Scan(&p.Id, &p.Age, &p.FirstName, &p.LastName, &p.Email) //сканирования и записи данных
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	tmpl, _ := template.ParseFiles("index.html") // полученные данные отправляет в html template
	tmpl.Execute(w, products)                    //добовляет данные и выводит  на страницу html
}

func main() {

	db, err := sql.Open("mysql", "root:alfadog1@/bipusdb") //Соединение с базой данных

	if err != nil {
		log.Println(err)
	}
	database = db //присваивает значения из базы данных

	// err = database.Ping() //проверяет живой ли сервер  если надо и соединяет с сервером
	// if err != nil {
	// 	os.Exit(122)
	// }

	defer db.Close()
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/create", CreateHandler)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":3000", nil)
}

//Как выглядит база данных
// use bipusdb;
// CREATE TABLE users
// (
//     Id INT AUTO_INCREMENT PRIMARY KEY,
//     Age INT DEFAULT 18 CHECK(Age >0 AND Age < 100),
//     FirstName VARCHAR(20) NOT NULL,
//     LastName VARCHAR(20) NOT NULL,
//     Email VARCHAR(30) CHECK(Email !='')
// );
