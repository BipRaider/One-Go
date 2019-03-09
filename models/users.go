package models

import (
	"errors"

	"../hash"
	"../rand"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"golang.org/x/crypto/bcrypt"
)

var (
	//ErrNorFound is returned when a resource cannot be found
	// in the database.
	ErrNotFaund = errors.New("models: resource not found")
	//ErrInvalidI is returned when  an invalid ID is provided
	// to a mathod like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid, must be > 0")

	ErrInvalidEmail = errors.New("models:invalid email address provided")

	ErrInvalidPassword = errors.New("models :invalid password provided")
)

const userPwPepper = "secret-random-string" // любую страку написать для усложнения паролей
const hmacSecretKey = " secret-hmac-key"

// User the user model in our database
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unigue_index"`
}

//UserDB is used to interact with the users database
//
//For pretty much all single user queries:
//if  the user  if found ,we will return a nill error
//if the user if not found ,we will return ErrRecordNotFound
//if these is another error ,we will return an error with
//more information about what went wrong
//
// For single user queries, any error but ErrNotFound should
//probably result in a 500 error
type UserDB interface {
	//Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	//Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id *uint) error

	// Used to close  a DB connectionInfo
	Close() error

	//Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods used to manipulate  and
// work with the user model
type UserService interface {
	//Authenticate will verify the provided email address and
	// password are correct.If they are correct ,the user
	//corresponding to than email will be returned .Otherwise
	// You will  receive either:
	//ErrNotFound,ErrInvalidPassword ,or another error if
	// something goes wrong.
	Authenticate(email, password string) (*User, error)

	UserDB
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

var _ UserService = &userService{}

type userService struct {
	UserDB
}

//Authenticate  can be used to authenticate a user with the
// provaided email address and password

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper)) // функция для разшифрофки  хешированый пороль
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword: ////  выводит ошибку хеша
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, err
}

var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
}

// func (uv *userValidator) ByID(id uint) (*User, error) {
// 	if id <= 0 {
// 		return nil, errors.New("Invalid id")
// 	}
// 	return uv.UserDB.ByID(id)
// }

func newUserGorm(connectionInfo string) (*userGorm, error) {
	//Соединение с базой данных  !!ВАЖНО ?parseTime=true  добисывать в конце если надо чтобы выводило время
	db, err := gorm.Open("mysql", connectionInfo) //"root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local"
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

//ByID will look up by the id provided
//if  the user  if found ,we will return a nill error
//if the user if not found ,we will return ErrRecordNotFound
//if these is another error ,we will return an error with
//more information about what went wrong
//
// as a general rule, any error but ErrNotFound should
//probably result in a 500 error
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ? ", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user with the given email address and
//return that user
//if  the user  if found ,we will return a nill error
//if the user if not found ,we will return ErrRecordNotFound
//if these is another error ,we will return an error with
//more information about what went wrong
//
// as a general rule, any error but ErrNotFound should
//probably result in a 500 error
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err

}

//ByRemember looks up  a user with the given remember token
// and returns that user .This method will handler hashing
//the token for us.
//Errors are  the same as ByEmail
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User

	rememberHash := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//first  will query using  the provided gorm.BD and it will
//get the first item returned and place it into dst. if
// nothing is found in the query , it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFaund
	}
	return err
}

//Create will create the provided user and backfill data
// like the ID, CreatedAt and UpdateAt fileds.
func (ug *userGorm) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)                              // этим услажнили просто пароль
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost) // используется для хеширования пароля
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken() //c генерирует рэндом токен
		if err != nil {
			return err
		}
		user.Remember = token

	}
	user.RememberHash = ug.hmac.Hash(user.Remember)

	return ug.db.Create(user).Error
}

//Update will update the provided user with all of the database
//in the  provaided user object
func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

//Delete will delete the user with the proveided ID
func (ug *userGorm) Delete(id *uint) error {
	if *id == 0 {
		return ErrInvalidID
	}
	user := User{
		Model: gorm.Model{
			ID: *id,
		},
	}

	return ug.db.Delete(&user).Error
}

//Фнкция для closes the Service with database
func (ug *userGorm) Close() error { return ug.db.Close() }

//DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error { // удалит таблицу если существует
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()

}

//AutoMigrate will attempt to autonatically migrate the
//user table
//Добовляет в базу данных нехватающих полей
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
