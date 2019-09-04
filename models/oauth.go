package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

//OAuth -  тип данных отправляемых  и принемаемых из DrobBox
type OAuth struct {
	gorm.Model
	UserID  uint   `gorm:"not null;unique_index:user_id_service"`
	Service string `gorm:"not null;unique_index:user_id_service"`
	oauth2.Token
}

// //queryOpt  тип для запрос данных из базы даных
// type queryOpt func(db *gorm.DB) *gorm.DB

// // byID используеться для выдачи ID адреса пользователя из database
// func byID(id uint) queryOpt {
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Where("id=?", id) // возвращаем даные и з базыданых
// 	}
// }
// OAuthDB.Find(byID(123),byUserID(321))

func NewOAuthService(db *gorm.DB) OAuthService {
	return &oauthValidatir{&oauthGorm{db}}
}

type OAuthService interface {
	OAuthDB
}

type OAuthDB interface {
	Find(userID uint, service string) (*OAuth, error)
	Create(oauth *OAuth) error
	Delete(id uint) error
}

///--------------------------
type oauthValidatir struct {
	OAuthDB
}

func (ov *oauthValidatir) Create(oauth *OAuth) error {
	err := runOAuthValFuncs(oauth,
		ov.userIDRequired,
		ov.serviceRequired,
	)
	if err != nil {
		return err
	}
	return ov.OAuthDB.Create(oauth)
}

//Delete  will delete the oauth with the provided ID
func (ov *oauthValidatir) Delete(id uint) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return ov.OAuthDB.Delete(id)
}

//------------------------------
func (ov *oauthValidatir) userIDRequired(oauth *OAuth) error {
	if oauth.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}
func (ov *oauthValidatir) serviceRequired(oauth *OAuth) error {
	if oauth.Service == "" {
		return ErrServiceRequired
	}
	return nil
}

//-------------------------------
var _ OAuthDB = &oauthGorm{}

type oauthGorm struct {
	db *gorm.DB
}

func (og *oauthGorm) Find(userID uint, service string) (*OAuth, error) {
	var oauth OAuth
	db := og.db.Where("user_id = ? ", userID).Where("service = ?", service)
	err := first(db, &oauth)
	return &oauth, err
}

func (og *oauthGorm) Create(oauth *OAuth) error {
	return og.db.Create(oauth).Error
}
func (og *oauthGorm) Delete(id uint) error {
	oauth := OAuth{Model: gorm.Model{ID: id}}
	return og.db.Unscoped().Delete(&oauth).Error
}

//-------------------------------

type oauthValFunc func(*OAuth) error

//проверяет на ошибкию если есть ошибки выводит их в браузере
func runOAuthValFuncs(oauth *OAuth, fns ...oauthValFunc) error {
	for _, fn := range fns {
		if err := fn(oauth); err != nil {
			return err //выводит ошибку
		}
	}
	return nil
}

//---------------------------------
