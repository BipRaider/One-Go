package controllers

import (
	"fmt"
	"net/http"
	"os"

	"../models"
	"../rand"
	"../views"
)

const bs string = "bootstrap"

func NewUser(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView(bs, "users/new"),
		NewFaq:    views.NewView(bs, "users/faq"),
		LoginView: views.NewView(bs, "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	NewFaq    *views.View
	LoginView *views.View
	us        models.UserService
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
	type Alert struct {
		Level   string
		Message string
	}
	type Data struct {
		Alert Alert
		Yield interface{}
	}
	a := Alert{
		Level:   "warning",                                // можно вписать такие имена  класа success,info,danger,warning
		Message: "Successfully rendered a dynamic alert!", // выводин на экран данное сообщение
	}
	d := Data{
		Alert: a,
		Yield: "Hello!",
	}
	if err := u.NewView.Render(w, d); err != nil {
		os.Exit(9)
	}

}

//This is used to render the form where  can create
// a new FAQ message
// GET /signup
func (u *Users) NewFaqGet(w http.ResponseWriter, r *http.Request) {
	if err := u.NewFaq.Render(w, nil); err != nil {
		os.Exit(91)
	}
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
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
	err := u.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)

	//-----1.1---1.2---
	// if err := r.ParseForm(); err != nil { //----1.1----
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
	// 	os.Exit(81)
	// }
	// fmt.Fprintln(w, form)

}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//Login is used to verify the provided email address
//and password and then log  the user in if they are correct
//POST/login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFaund:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrPasswordInCorrect:
			fmt.Fprintln(w, "Invalid password provided.")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

//signIn is used to sign the given user  in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true, ////Когда вы помечаете куки с флагом HttpOnly, он сообщает браузеру,
		// что этот конкретный куки должен быть доступен только серверу. Любая попытка доступа к куки из клиентского скрипта строго запрещена.
	}
	http.SetCookie(w, &cookie)
	return nil
}

//CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookei, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookei.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)

}

//https://github.com/gorilla/mux
//https://golang.org/pkg/net/http/#pkg-examples -----1.1-----
