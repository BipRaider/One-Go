package views

import (
	"log"

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
type Data struct {
	Alert *Alert
	User  *models.User
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
