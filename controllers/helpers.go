package controllers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

// функция  обработки в водимых данных user
func parseForm(r *http.Request, dst interface{}) error {
	//  ParseForm - Для запросов POST, PUT и PATCH он также анализирует тело запроса как форму и
	//  помещает результаты как в r.PostForm,так и в r.Form. Параметры тела запроса имеют
	//  приоритет над значениями строки запроса URL в r.Form
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}

func parseURLparams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder() // NewDecoder возвращает новый декодер.
	// Вызываем функцию IgnoreUnkownKeys, чтобы сказать декодеру схемы игнорировать ключ токена CSRF
	dec.IgnoreUnknownKeys(true)                     // Для сохранения обратной совместимости значением по умолчанию является false.
	if err := dec.Decode(dst, values); err != nil { //Decode- декодирует значения из map[string][]string в struct.
		return err
	}
	return nil
}

// Если размер тела запроса еще не был ограничен MaxBytesReader,размер ограничен 10 МБ.

// ParseMultipartForm - вызывает ParseForm автоматически.
