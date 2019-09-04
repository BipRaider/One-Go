package models

import (
	"strings"
)

const (
	//ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFaund modelError = "models: resource not found"
	//ErrInvalidEmail is returned when an invalid email addres
	ErrInvalidEmail modelError = "models: invalid email address provided"

	//ErrPasswordInCorrect is returned when an invalid password
	//is used when attempting to authenticate a user.
	ErrPasswordInCorrect modelError = "models: invalid password provided"

	//ErrEmailRequired  is returned  when  an email address is
	// not provided when creating a user
	ErrEmailRequired modelError = "models: Email address is required"

	//ErrEmailInvalid is returned  when an email address provaided
	// does not match any of our requirements
	ErrEmailInvalid modelError = "models: Email address is not valid"

	//ErrEmailTaken  is returned  when an update or create is attempted
	//with an email address that is already in use.
	ErrEmailTaken modelError = "models: Email address is already taken"

	//ErrPasswordTooShort is returned when  an update or create is
	//attempted with a user password that is less than 8 characters.
	ErrPasswordTooShort modelError = "models: Passsword must be at least 8 characters ling"

	//ErrTokenInvalid is returned when token provided is not valid
	ErrTokenInvalid modelError = "models: token provided is not valid"

	// ErrPasswordRequired is returned when an create is attempted
	//without a user  password  provided.
	ErrPasswordRequired modelError = "models: Passsword is required"

	//ErrTitleRequired is returned when
	ErrTitleRequired modelError = "models: Title is required"

	//ErrInvalidID is returned when  an invalid ID is provided
	// to a mathod like Delete.
	ErrInvalidID privateError = "models: ID provided was invalid, must be > 0"

	//ErrRememberTooShort is returned when a remember token is not
	//at the least 32 bytes
	ErrRememberTooShort privateError = "models: Remember token must be at bytes"

	// ErrRememberRequired is returned when an create or update is attempted
	//without a valid user remember token hash.
	ErrRememberRequired privateError = "models: Remember token is required"

	//ErrUserIDRequired  is returned when
	ErrUserIDRequired privateError = "models: User ID is required"

	//ErrServiceRequired is returned when
	ErrServiceRequired privateError = "models: Service is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	//(strings.Replace("в водимые даные" , "что замененить", "на что заменять", сколько первых совподений надо заменить(int)если -1 то меняться будут все совподения )
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")     // строку приобразовал в срез
	split[0] = strings.Title(split[0]) // все первые буквы строк будут заглавные
	return strings.Join(split, " ")    // обьеденили срез в одну строку   (первое срез , что вставлять между )
}

//-----------------
type privateError string

func (e privateError) Error() string {
	return string(e)
}
