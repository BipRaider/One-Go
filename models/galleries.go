package models

import "github.com/jinzhu/gorm"

//Gallery is our image container resources that visitors view
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gotm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

// func (gd *galleryGorm) Create(g *Gallery) error {

// }

type galleryGorm struct {
	db *gorm.DB
}
