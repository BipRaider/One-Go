package middleware

///Требования к юзеру на афторизацию тд
import (
	"net/http"

	"../context"
	"../models"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will return an http.HandlerFunc that will
// check to see if a user is logged in and then either
// call next(w, r) if they are, or redirect them to the
// login page if they are not.

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookei, err := r.Cookie("remember_token") // если пользователь не зарегестрирован(нету куки) то перенаправляет на регистрацию
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mw.UserService.ByRemember(cookei.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
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
