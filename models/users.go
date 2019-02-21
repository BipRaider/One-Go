package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	//ErrNorFound is returned when a resource cannot be found
	// in the database.
	ErrNorFaund = errors.New("models: resource not found")
)

func NewUserService(connectionInfo string) (*UserService, error) {
	//Соединение с базой данных  !!ВАЖНО ?parseTime=true  добисывать в конце если надо чтобы выводило время
	db, err := gorm.Open("mysql", connectionInfo) //"root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local"

	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

//ByID will look up by the id provided
//1 - user ,nil
//2 - nil, ErrorNotFound
//3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ? ", id).
		First(&user).
		Error

	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNorFaund
	default:
		return nil, err
	}
}

//Create will create the provided user and backfill data
// like the ID, CreatedAt and UpdateAt fileds.
func (us *UserService) Create(user *User) error {

	return us.db.Create(user).Error
}

//Update will update the provided user with all of the database
//in the  provaided user object
func (us *UserService) Update(user *User) error {
	return us.db.Update(user).Error

}

//Фнкция для closes the Service with database
func (us *UserService) Close() error {
	return us.db.Close()
}

//DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
