package models

import "github.com/jinzhu/gorm"

func NewServices(connectionInfo string) (*Services, error) {
	//Соединение с базой данных  !!ВАЖНО ?parseTime=true  добисывать в конце если надо чтобы выводило время
	db, err := gorm.Open("mysql", connectionInfo) //"root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local"
	if err != nil {
		return nil, err
	}
	return &Services{}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
}
