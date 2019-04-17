package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

// функция  обработки в водимых данных user
func parseForm(r *http.Request, dst interface{}) error {

	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
