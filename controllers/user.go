package controllers

import (
	"fmt"
	"net/http"

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
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {

	u.NewView.Render(w, nil)
}

//This is used to process the signup form when a user tries to
//create a new user account.
//
//POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a fake message.Pretend that we created the user account")
}

//https://github.com/gorilla/mux
