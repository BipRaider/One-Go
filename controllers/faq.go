package controllers

import (
	"fmt"
	"net/http"
	"os"

	"../views"
)

func NewFAQ() *FAQUsers {
	return &FAQUsers{
		NewFaq: views.NewView("bootstrap", "static/faq"),
	}
}

type FAQUsers struct {
	NewFaq *views.View
}

// GET /signup
func (u *FAQUsers) NewFaqGet(w http.ResponseWriter, r *http.Request) {
	if err := u.NewFaq.Render(w, nil); err != nil {
		os.Exit(9)
	}
}

type SignupFormFaq struct {
	Quastion string `schema:"faq"`
}

//POST /signup
func (u *FAQUsers) NewFaqCreate(w http.ResponseWriter, r *http.Request) {

	var form SignupFormFaq
	if err := parseForm(r, &form); err != nil {
		os.Exit(82)
	}
	fmt.Fprintln(w, form)
}
