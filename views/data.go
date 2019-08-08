package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"../models"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	//AlertMsgGeneric is displayed when any random errror
	//is encountered by our backend.
	AlertMsgGeneric = "Something went wrong .Please try again ,and contact us if the problem persists"
)

//Alert is used to render Bootstrap Alert messager in the
//bootstrap.gohtml template
type Alert struct {
	Level   string
	Message string
}

// Data is the top leval structure that views expect data
//to come in
type Data struct {
	Alert *Alert
	User  *models.User
	CSRF  template.HTML //Кодируем страницы
	Yield interface{}
}

//----------------------------
// SerAler  возвращает ошибки в барузер
func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PiblicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		log.Println(err)
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}
}

// AlertError  выводит ошибки в сообщений на экран ,что мы задаём.
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

type PiblicError interface {
	error
	Public() string
}

//Если пользователь не загрузит перенаправление через 5 минут, оно просто истечем.
func persistAlert(w http.ResponseWriter, alert Alert) {
	// We don't want alerts showing up days later. If the
	// user doesnt load the redirect in 5 minutes we will
	// just expire it.

	expiresAt := time.Now().Add(5 * time.Minute)
	fmt.Println("cookie expires at..", expiresAt)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func clearAlert(w http.ResponseWriter) {
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

// Если один из файлов cookie отсутствует, мы будем считать, что предупреждение недействительно, и вернем ноль
func getAlert(r *http.Request) *Alert {
	// If either cookie is missing we will assume the alert
	// is invalid and return nil
	lvl, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}
	alert := Alert{
		Level:   lvl.Value,
		Message: msg.Value,
	}
	return &alert
}

// RedirectAlert accepts all the normal params for an
// http.Redirect and performs a redirect, but only after
// persisting the provided alert in a cookie so that it can
// be displayed when the new page is loaded.
func RedirectAlert(w http.ResponseWriter, r *http.Request, urlStr string, code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, urlStr, code)

}
