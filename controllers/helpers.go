package controllers

import (
	"net/http"

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
	dec := schema.NewDecoder()                          // NewDecoder возвращает новый декодер.
	dec.IgnoreUnknownKeys(true)                         // Для сохранения обратной совместимости значением по умолчанию является false.
	if err := dec.Decode(dst, r.PostForm); err != nil { //Decode- декодирует значения из map[string][]string в struct.
		return err
	}
	return nil
}

// Если размер тела запроса еще не был ограничен MaxBytesReader,размер ограничен 10 МБ.

// ParseMultipartForm - вызывает ParseForm автоматически.
