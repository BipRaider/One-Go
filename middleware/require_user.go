package middleware

///Требования к юзеру на афторизацию тд
import (
	"net/http"
	"strings"

	"../context"
	"../models"
)

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		//if the user is requesting a static asset or image
		// we will not need  to lookup the current user so we skip
		// doing that.
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}

		cookei, err := r.Cookie("remember_token") // если пользователь не зарегестрирован(нету куки) то перенаправляет на регистрацию
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookei.Value)
		if err != nil {
			next(w, r)
			return
		}
		// Get the context from our request
		ctx := r.Context() // This is how we get a request's context

		// Create a new context from the existing one that has
		// our user stored in it with the private user key
		ctx = context.WithUser(ctx, user) // отправка данных в контекс

		// Create a new request from the existing one with our
		// context attached to it and assign it back to `r`.
		r = r.WithContext(ctx)

		// Call next(w, r) with our updated context.
		next(w, r)
	})
}

///-------------------------------------------------------------
//RequireUser assumes that User middleware has already been run
//otherwise it will not work  correctly.
type RequireUser struct {
	User
}

//Apply  assumes that User middleware has already been run
//otherwise it will not work  correctly.
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will return an http.HandlerFunc that will
// check to see if a user is logged in and then either
// call next(w, r) if they are, or redirect them to the
// login page if they are not.

//ApplyFn  assumes that User middleware has already been run
//otherwise it will not work  correctly.
// работает с индефецированым пользователем и обрабатывает все функций с этим индефецированым пользователем
// если пользователь не индефецирован отправляет на страничку логина
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
