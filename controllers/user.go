package controllers

import (
	"fmt"
	"net/http"
	"os"

	"../views"
)

func NewUser() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}

}

type Users struct {
	NewView *views.View
}

//This is used to render the form where  can create
//a new user account.
//Это используется для визуализации формы, где можно создать
// новая учетная запись пользователя
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {

	if err := u.NewView.Render(w, nil); err != nil {
		os.Exit(9)
	}
}

//This is used to process the sign up form when a user tries to
//create a new user account.
// Используется для обработки формы регистрации, когда пользователь пытается
// создать новую учетную запись пользователя.
//
//POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { //----1.1----
		os.Exit(8)
	}
	//r.PostForm = map[string][]string
	fmt.Fprintln(w, r.PostForm["emeil"]) // выводит срез записи  после в вода даных в sign up
	// fmt.Fprintln(w, r.PostFormValue("emeil")) // выводит первую запись после в вода даных в sign up
	fmt.Fprintln(w, r.PostForm["password"])
	// fmt.Fprintln(w, r.PostFormValue("password"))

}

//https://github.com/gorilla/mux
//https://golang.org/pkg/net/http/#pkg-examples -----1.1-----
