package main

import "fmt"

type MysqlConfig struct {
	// Hostname string `jsom:"host"` Данные хост можно найти в настройка mysql   (Данные строчки нужны для Postgres)
	// Port     int    `jsom:"port"`Данные port можно найти в настройка mysql (Данные строчки нужны для Postgres)
	Username string `json:"user"`     // Имя поста в Mysql
	Password string `json:"password"` // пароль
	DS       string `json:"ds"`       // Обьединяет парль с юзер нейв для коректного пути
	BDName   string `json:"name"`     // Название basedata
	Decoder  string `json:"decoder"`  // обезательная часть чтобы коректно отображалась информация и читалась из basedata
} //mysqlinfo = "root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local"
func (c MysqlConfig) Dialect() string {
	return "mysql"
}

func (c MysqlConfig) ConnectionInf() string {
	if c.Password == "" {
		return fmt.Sprintf("%s%s%s%s", c.Username, c.DS, c.BDName, c.Decoder)
	}
	return fmt.Sprintf("%s%s%s%s%s", c.Username, c.Password, c.DS, c.BDName, c.Decoder)
}
func DefaultMysqlConfig() MysqlConfig {
	return MysqlConfig{
		Username: "root:",
		Password: "alfadog1",
		DS:       "@/",
		BDName:   "bipusdb",
		Decoder:  "?charset=utf8&parseTime=True&loc=Local",
	}
}

//---------------
type Config struct {
	Port int
	Env  string
}

func (c Config) isProd() bool {
	return c.Env == "prod" // выдаёт false
}

func DefaultConfig() Config {
	return Config{
		Port: 3000,  // локальный порт проэкта
		Env:  "dev", //// Secure -устанавливает флаг безопасности в куки. По умолчанию true
		// Установите  «false» в противном случае файл cookie не будет отправляться по небезопасному каналу
	}
}

//------
// #models/users.go

// const userPwPepper = "secret-random-string" // любую страку написать для усложнения паролей
// const hmacSecretKey = " secret-hmac-key"    // любую страку написать для усложнения паролей

// #models/service.go

// //Соединение с базой данных  !!ВАЖНО ?charset=utf8&parseTime=True&loc=Local  добисывать в конце если надо чтобы выводило время
// 	db, err := gorm.Open("mysql", connectionInfo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	db.LogMode(true) // устанавливаем онаним ведения журнала (True для подробных ),(False - выводит ток ошибки )
