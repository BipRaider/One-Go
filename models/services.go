package models

//Обьеденяем все микросерверы в один общий
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ServicesConfig func(*Services) error

//Соединяемся с базой данных
func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		//Соединение с базой данных  !!ВАЖНО ?charset=utf8&parseTime=True&loc=Local  добисывать в конце если надо чтобы выводило время
		db, err := gorm.Open(dialect, connectionInfo) //"root:password@/NameDB?charset=utf8&parseTime=True&loc=Local"
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

// LogMode установить режим журнала, `true` для подробных журналов,` false` для отсутствия журнала, по умолчанию, будет печатать только журналы ошибок
func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

// Проверка пользователя  на коодировку
func WithUser(pepper, hmacKye string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKye)
		return nil
	}
}

// запуск галирей
func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

//Запускаем просмотр картинок
func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

//Запуск функций  которые запускают Бд, Проверку на кодироания , запускают галерей и просмотр картинок
func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

//Closse the database connection
func (s *Services) Close() error { return s.db.Close() }

//DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error { // удалит таблицы если существует
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
