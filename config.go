package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type MysqlConfig struct {
	// Hostname string `jsom:"host"` Данные хост можно найти в настройка mysql   (Данные строчки нужны для Postgres)
	// Port     int    `jsom:"port"`Данные port можно найти в настройка mysql (Данные строчки нужны для Postgres)
	Username string `json:"username"` // Имя поста в Mysql
	Password string `json:"password"` // пароль
	DS       string `json:"ds"`       // Обьединяет пароль с юзернейв для коректного пути
	DBName   string `json:"dbname"`   // Название database
	Decoder  string `json:"decoder"`  // обезательная часть чтобы коректно отображалась информация и читалась из basedata
} //mysqlinfo = "root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local"
func (c MysqlConfig) Dialect() string {
	return "mysql"
}

//приобразуем путь в строку
func (c MysqlConfig) ConnectionInf() string {
	if c.Password == "" {
		return fmt.Sprintf("%s%s%s%s", c.Username, c.DS, c.DBName, c.Decoder)
	}
	return fmt.Sprintf("%s%s%s%s%s", c.Username, c.Password, c.DS, c.DBName, c.Decoder)
}

// Данные для в соединения с базой данных
func DefaultMysqlConfig() MysqlConfig {
	return MysqlConfig{
		Username: "root:",    //имя в MySql
		Password: "alfadog1", //пароль в MySql
		DS:       "@/",
		DBName:   "bipusdb",                                // название бд в в MySql
		Decoder:  "?charset=utf8&parseTime=True&loc=Local", // кодировка бд лоя коректного вывода
	}
}

//---------------
type Config struct {
	Port     int         `json:"port"`     // Адрес сервера
	Env      string      `json:"env"`      // для безопасности в куки
	Pepper   string      `json:"pepper"`   // любую стрoку написать для усложнения паролей
	HMACKey  string      `json:"hmac_key"` // любую стрoку написать для усложнения паролей
	Database MysqlConfig `json:"database"`
}

// для вывода значения False
func (c Config) isProd() bool {
	return c.Env == "prod" // выдаёт false
}

func DefaultConfig() Config {
	return Config{
		Port:     3000,                   // локальный порт проэкта
		Env:      "dev",                  //// Secure -устанавливает флаг безопасности в куки. По умолчанию true// Установите  «false» в противном случае файл cookie не будет отправляться по небезопасному каналу
		Pepper:   "secret-random-string", // Для услажнения пароля в веденного пользователем
		HMACKey:  "secret-hmac-key",      // для усложнения хешированя юзера
		Database: DefaultMysqlConfig(),   // Запуск по умолчанию конфига
	}
}

// Приобразуем данные  из конфига с  помощью JSON для отправки данных на сервер для запуска БД
func LoadConfig(configReq bool) Config {
	f, err := os.Open(".config")
	if err != nil {
		if configReq {
			panic(err)
		}
		log.Println(err)
		fmt.Println("Using the default config....")
		return DefaultConfig()
	} // Если не будет найдень другой конфиг то по умолчанию запустится стандартный

	//**
	var c Config
	dec := json.NewDecoder(f) // приобразуем данные  для чтения JSON
	err = dec.Decode(&c)      //Декодируем даные из  приобразованого декотера выше в заполняем тип Config
	if err != nil {
		log.Println(err)
		fmt.Println("Using the default config....")
		return DefaultConfig() // Если не будет найдень другой конфиг то по умолчанию запустится стандартный
	}
	fmt.Println("Successfully loaded .config")
	return c
}
