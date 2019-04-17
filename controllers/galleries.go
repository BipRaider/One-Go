package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"../context"
	"../models"
	"../views"
	"github.com/gorilla/mux"
)

const (
	ShowGallery = "show_gallery"
)

func NewGalleries(gs models.GalleryService, mr *mux.Router) *Galleries {
	return &Galleries{
		New:      views.NewView(bs, "galleries/new"),
		ShowView: views.NewView(bs, "galleries/show"),
		gs:       gs,
		r:        mr,
	}
}

type Galleries struct {
	New      *views.View
	ShowView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

//
type GalleryForm struct {
	Title string `schema:"title"`
}

// переходи по сылка по в "/galleries/{id:[0-9]+}"
//Get /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFaund:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoooops!!!  Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)

}

//POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) { // Обрабатывает в водимые данные в браузере

	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	// This is our new code
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	url, err := g.r.Get(ShowGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		//TODO: Make this go to the index page
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}
