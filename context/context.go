package context

import (
	"context"

	"../models"
)

const (
	userKey privateKey = "user"
)

type privateKey string

//This function that accepts an existing context and a
//user, and then returns a new context with that user set as a value.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user) // отправка данных в контекс
}

// Используеться как афтаризация для всех страниц.
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil { // Получить значение, хранящееся в userKye /// Retrieve the value stored at "my-key"
		if user, ok := temp.(*models.User); ok { // type conversion !
			return user
		}
	}
	return nil
}

//https://www.youtube.com/watch?v=LSzR0VEraWw
