package models

import (
	"regexp"
	"strings"

	"../hash"
	"../rand"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"golang.org/x/crypto/bcrypt"
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
	Delete(id uint) error
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

func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}
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
			return nil, ErrPasswordInCorrect
		default:
			return nil, err
		}
	}
	return foundUser, err
}

//-------------------------------------------------------------------------------
type userValFunc func(*User) error

//проверяет на ошибкию если есть ошибки выводит их в браузере
func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err //выводит ошибку
		}
	}
	return nil
}

//-----------------------------------------------------------------------------
var _ UserDB = &userValidator{}

// функция для добовления  данных в type userValidetor  для проверки мыла на валидность заполнение
// emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`), // используется для сопоставления адресов электронной почты
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp //https://gobyexample.com/regular-expressions
}

// ByEmail will normalize the email address before calling
// ByEmail an the UserDB field.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

//ByRemember will hash the remember token and  then call
//Byremember on the subsequent UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

//Create will create the provided user and backfill data
// like the ID, CreatedAt and UpdateAt fileds.
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

//Update will hash a remembeer token if it is  provided .
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

//bcyrpt  password wiil hash a user is password with a
//predefined pappe (userPwPepper) and bcrypt if the
// Password  field is not the empty  string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)                              // этим услажнили просто пароль
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost) // используется для хеширования пароля
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

//
func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

///
func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken() //c генерирует рэндом токен
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

////
func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

/////
func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}

//////
func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrInvalidID
		}
		return nil
	})
}

///////
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)   // выводит каждой буквы уникодовское  цифры нижнего регистра
	user.Email = strings.TrimSpace(user.Email) // уберает все (\t \r \n) из строки
	return nil
}

////////
func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

// проверяет правельность заполнения мыла
func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		//если мыло не соотвецтвует данной форме выводится ошибка
		return ErrEmailInvalid
	}
	return nil

}

//проверяет на наличие мыла в базе данных и выыодит ошибку если такое мыло есть
func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFaund {
		//Email address is noy taken
		return nil
	}
	if err != nil {
		return err
	}
	//We found a user w/ this email address...
	//If the found user  has the same ID as this user , it is
	//an update and this is the same user
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

// Проверка пароля на длину
func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

//
func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}
func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}

	return nil
}

///----------------------------------------------------------------------------------------------------------------------------
var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
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
// and returns that user .This method expects the remember token to already be hashed
//Errors are  the same as ByEmail
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User

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
	return ug.db.Create(user).Error
}

//Update will update the provided user with all of the database
//in the  provaided user object
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

//Delete will delete the user with the proveided ID
func (ug *userGorm) Delete(id uint) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	return ug.db.Delete(&user).Error
}
