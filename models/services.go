package models

import "github.com/jinzhu/gorm"

func NewServices(connectionInfo string) (*Services, error) {
	//Соединение с базой данных  !!ВАЖНО ?parseTime=true  добисывать в конце если надо чтобы выводило время
	db, err := gorm.Open("mysql", connectionInfo) //"root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local"
	if err != nil {
		return nil, err
	}
	db.LogMode(true) // устанавливаем онэим ведения журнала (True для подробных ),(False - выводит ток ошибки )
	return &Services{
		User: NewUserService(db),
		db:   db,
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

//Closse the database connection
func (s *Services) Close() error { return s.db.Close() }

//DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error { // удалит таблицу если существует
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()

}

//AutoMigrate will attempt to autonatically migrate all tables
//Добовляет в базу данных нехватающих полей
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}
