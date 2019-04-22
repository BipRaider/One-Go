package models

import (
	"github.com/jinzhu/gorm"
)

//Gallery is our image container resources that visitors view
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gotm:"not_null"`
}

//----------------------------------------------------------------------------
type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
}

//--------------------------------------------------------------------------------
func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

///-----------------------------------------
type galleryValFunc func(*Gallery) error

//проверяет на ошибкию если есть ошибки выводит их в браузере
func runGalleryValFuncs(gallery *Gallery, fns ...galleryValFunc) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err //выводит ошибку
		}
	}
	return nil
}

///------------------------------------------
func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}
func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

//Delete  will delete the gallery with the provided ID
func (gv *galleryValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return gv.GalleryDB.Delete(id)
}

//проверка на  соотвецтвие   необходимого id пользователя
func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}
func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

///------------------------------------------
var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

//
func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ? ", id)
	err := first(db, &gallery)
	return &gallery, err
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}
func (gg *galleryGorm) Delete(id uint) error {
	gallery := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&gallery).Error
}
