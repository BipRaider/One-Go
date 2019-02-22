package controllers

import (
	"fmt"
	"net/http"
	"os"

	"../models"
	"../views"
)

func NewUser(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		NewFaq:  views.NewView("bootstrap", "static/faq"),
		us:      us,
	}
}

type Users struct {
	NewView *views.View
	NewFaq  *views.View
	us      *models.UserService
}

//GET Reading a resource ПОЛУЧИТЬ Чтение ресурса
// POST Creating a resource POST Создание ресурса
// PUT Updating a resource  PUT Обновление ресурса
// PATCH Updating a resource PATCH Обновление ресурса
// DELETE Deleting a resource УДАЛИТЬ Удаление ресурса

// The difference between PUT and PATCH.
// Both PUT and PATCH are used to represent updating a resource, but they both do it
// in a fundamentally different way. PUT generally is expected to accept an entirely
// new representation of an object, even if some of the fields didn’t change, while
// PATCH was proposed as a way to update resources and also signify that you won’t
// be passing all of the fields for the resource, but instead will only be providing the
// updated fields.
// For all practical purposes you mostly just need to remember that these are both
// used to update resources.

//This is used to render the form where  can create
//a new user account.
//Это используется для визуализации формы, где можно создать
// новая учетная запись пользователя
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) { //обрабатывает Html шаблоны и вывоодит в браузер .
	if err := u.NewView.Render(w, nil); err != nil {
		os.Exit(9)
	}

}

// GET /signup
func (u *Users) NewFaqGet(w http.ResponseWriter, r *http.Request) {
	if err := u.NewFaq.Render(w, nil); err != nil {
		os.Exit(9)
	}
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"emeil"`
	Password string `schema:"password"`
	Quastion string `schema:"faq"`
}

//This is used to process the sign up form when a user tries to
//create a new user account.
// Используется для обработки формы регистрации, когда пользователь пытается
// создать новую учетную запись пользователя.
//
//POST /signup

func (u *Users) Create(w http.ResponseWriter, r *http.Request) { // Обрабатывает в водимые данные в браузере
	//----1.3---
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		os.Exit(82)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // выводит ошибку если в в базе данных есть такой ID
		return
	}
	fmt.Fprintln(w, form)

	//-----1.1---1.2---
	// if err := r.ParseForm(); err != nil { //----1.1----
	// 	err = errors.New("ERROR func user/Create at ParseForm")
	// 	os.Exit(8)
	// }
	//----1.1---
	//r.PostForm = map[string][]string
	// fmt.Fprintln(w, r.PostForm["emeil"]) // выводит срез записи  после в вода даных в sign up
	// fmt.Fprintln(w, r.PostFormValue("emeil")) // выводит первую запись после в вода даных в sign up
	// fmt.Fprintln(w, r.PostForm["password"])
	// fmt.Fprintln(w, r.PostFormValue("password"))
	//----1.2----
	// dec := schema.NewDecoder()
	// form := SignupForm{}
	// if err := dec.Decode(&form, r.PostForm); err != nil {

	// 	err = errors.New("ERROR func user/Create at Decode")
	// 	os.Exit(81)
	// }
	// fmt.Fprintln(w, form)

}

//https://github.com/gorilla/mux
//https://golang.org/pkg/net/http/#pkg-examples -----1.1-----

//152 page
