package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // подсоединить библиотеку mysql
)

type UserGorm struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"` // добовления типа данных в базуданых
	Color string
}

func main() {
	db, err := gorm.Open("mysql", "root:alfadog1@/bipusdb?parseTime=true") //Соединение с базой данных  !!ВАЖНО ?parseTime=true  добисывать в конце если надо чтобы выводило время
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//db.LogMode(true)            // отображает лог строки запроса в терменале
	db.AutoMigrate(&UserGorm{}) // автомотически отправляет данные в базу даных

	var u UserGorm
	name, email, color := getInfp() // выводит данные для дальнейшей записи
	us := UserGorm{
		Name:  name,
		Email: email,
		Color: color,
	}
	if err = db.Create(&us).Error; err != nil { // ЗАписывет данные в базу-данных
		os.Exit(1)
	} else {
		db.Last(&us)
		fmt.Println("Last name created in datanase \n", us)
	}
	// newDB := db.Where("email=?", "blah@blah.com")
	// newDB = newDB.Or("color=?", "wite")             // функция или
	// newDB = newDB.First(&u)
	// db = db.Where("email = ?", "bipus@gmail.com").First(&u)
	if err := db.Where("email = ?", "bipus@gmail.com").First(&u).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			fmt.Println("1111100011111")
		case gorm.ErrInvalidSQL:
			fmt.Println("sql")
		default:
			os.Exit(33)
		}
	}

	fmt.Println(u)

}

// запрос нужных данных для записи в базу данных
func getInfp() (name, email, color string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What id your name?")
	name, _ = reader.ReadString('\n') // Функция для запроса в вода fmt.Scan
	fmt.Println("Whatewr idewr yourer nameqweqwe?")
	email, _ = reader.ReadString('\n')
	fmt.Println("Whatewr color?")
	color, _ = reader.ReadString('\n')
	name = strings.TrimSpace(color) // Возвращает срез строки  без ( \n )
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	return name, email, color
}
