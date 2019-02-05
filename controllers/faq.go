package controllers

import (
	"fmt"
	"net/http"
	"os"

	"../views"
)

func NewFAQ() *FAQUsers {
	return &FAQUsers{
		NewView: views.NewView("bootstrap", "users/new"),
	}

}

type FAQUsers struct {
	NewView *views.View
}

// GET /signup
func (u *FAQUsers) New(w http.ResponseWriter, r *http.Request) {

	if err := u.NewView.Render(w, nil); err != nil {
		os.Exit(9)
	}
}

type SignupFormFaq struct {
	Quastion string `schema:"faq"`
}

//POST /signup

func (u *FAQUsers) Create(w http.ResponseWriter, r *http.Request) {

	var form SignupFormFaq
	if err := parseForm(r, &form); err != nil {
		os.Exit(82)
	}
	fmt.Fprintln(w, form)

}
