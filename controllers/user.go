package controllers

import (
	"log"
	"net/http"
	"time"

	"../context"
	"../email"
	"../models"
	"../rand"
	"../views"
)

const bs string = "bootstrap"

func NewUser(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		NewView:      views.NewView(bs, "users/new"),
		LoginView:    views.NewView(bs, "users/login"),
		ForgotPwView: views.NewView(bs, "users/forgot_pw"),
		ResetPwView:  views.NewView(bs, "users/reset_pw"),
		us:           us,
		emailer:      emailer,
	}
}

type Users struct {
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
	us           models.UserService
	emailer      *email.Client // клиент отправления писем
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
	var form SignupForm
	parseURLparams(r, &form)
	u.NewView.Render(w, r, form)
}

//This is used to render the form where  can create
// a new FAQ message
// GET /signup

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//This is used to process the sign up form when a user tries to
//create a new user account.
// Используется для обработки формы регистрации, когда пользователь пытается
// создать новую учетную запись пользователя.
//
//POST /signup
///
func (u *Users) Create(w http.ResponseWriter, r *http.Request) { // Обрабатывает в водимые данные в браузере
	var vd views.Data
	var form SignupForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	u.emailer.Welcom(user.Name, user.Email) // Оправим письмо к нам на мыло (то мыло которое мы хочем получать песьмо от регестрируемых пользователей)
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcom to BipGo.pw",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//Login is used to verify the provided email address
//and password and then log  the user in if they are correct
//POST/login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password) // запрос из базы нужных данных
	if err != nil {
		switch err {
		case models.ErrNotFaund:
			vd.AlertError("Invalid email address.")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcom to BipGo.pw",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

//Logout is used to delete a users sessions cookien (remember_token)
//and them will updata the user resource with a new remember token.
//POST/logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	// удаляем старый токен , заменяем старый на новый
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie) // устанавливаем новый токен

	user := context.User(r.Context())          //обновляем пользователя новым токеном запоминания
	token, _ := rand.RememberToken()           //Генерируем токен
	user.Remember = token                      //Запоминаем с генерированый токен
	u.us.Update(user)                          /// обновляем пользователя
	http.Redirect(w, r, "/", http.StatusFound) // отправляем пользователя на домашнюю страницу
}

//ResetPwForm is used to process the forgot password form
//and the reset password form.
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

//
//POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	//TODO: Process the forgot password form and initiate that process
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}
	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	err = u.emailer.ResetPw(form.Email, token) //Send the user an email with their token and password reset  instructions
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{ // отправляем пользователя на /reset страницу с выводом сообщения
		Level:   views.AlertLvlSuccess,
		Message: "Instructions for resetting your password have been emailed to you",
	})

}

//ResetPw displays the reset password for and has a method
//so that we can prefill the from data with a token provided
//via the URL query params
//
//GET /reset
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form

	if err := parseURLparams(r, &form); err != nil {
		vd.SetAlert(err)
	}
	u.ResetPwView.Render(w, r, vd)
}

//CompleteReset processed the rest password form
//
//POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}
	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	u.signIn(w, user)
	views.RedirectAlert(w, r, "/", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Your password has been reset and you have been logged in!",
	})
}

//signIn is used to sign the given user  in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken() //Генерируем токен
		if err != nil {
			return err
		}
		user.Remember = token   //Запоминаем с генерированый токен
		err = u.us.Update(user) // обновляем пользователя
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

//https://github.com/gorilla/mux
//https://golang.org/pkg/net/http/#pkg-examples -----1.1-----
//https://app.mailgun.com/app/account/setup
