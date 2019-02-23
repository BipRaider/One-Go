package controllers

import (
	"../views"
)

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView(bs, "static/home"),
		Contact: views.NewView(bs, "static/contact"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
