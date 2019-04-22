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
		EditView: views.NewView(bs, "galleries/edit"),

		gs: gs,
		r:  mr,
	}
}

type Galleries struct {
	New      *views.View
	ShowView *views.View
	EditView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

//
type GalleryForm struct {
	Title string `schema:"title"` // тип даных что в водеться в браузере  и берется из gallery/new.gohtml для передачи в функций
}

// переходи по сылка по в "/galleries/{id:[0-9]+}"
//Get /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	// Мы создадим объект views.Data и установим нашу галерею
	// как поле Yield, но технически нам не нужно
	// делать это и просто передать галерею в
	// Визуализация метода из-за переключателя типа, который мы закодировали в
	// метод Render.
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)

}

//Get /galleries/:id/edit  // переходим на строничку для редакций  данных в галерей  по айди юзира
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found  for edit", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd) ////вывод данных на экран
}

//POST/galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found  for update", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, vd) ////вывод данных на экран
		return
	}
	// в водимые даные перезаписываем в базу данных
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}
	// выводит на экран зеленую надпись что галереея обнавилась
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery succesfully update",
	}
	//вывод данных на экран
	g.EditView.Render(w, vd)
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
		Title:  form.Title, //  отправить строку для записи в базу данных из в вода на на странице браузера
		UserID: user.ID,    // запишит номер id user in DB
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

// используется для индефикаций по id user
func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r) //Vars возвращает переменные маршрута для текущего запроса, если таковые имеются.
	idStr := vars["id"] // Далее нам нужно получить переменную "id" из наших переменных.

	id, err := strconv.Atoi(idStr) // Наш idStr является строкой, поэтому мы используем функцию Atoi предоставляется пакетом strconv для преобразования его в целое число
	if err != nil {                // Эта функция также может возвращать ошибку, поэтому нам
		http.Error(w, "Invalid gallery ID", http.StatusNotFound) // нужно проверить на наличие ошибок и отобразить ошибку.
		return nil, err
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFaund:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoooops!!!  Something went wrong", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}
