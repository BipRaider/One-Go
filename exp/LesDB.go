package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // подсоединить библиотеку mysql
)

type UserGorm struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Color string
}

func main() {
	db, err := gorm.Open("mysql", "root:alfadog1@/bipusdb?parseTime=true") //Соединение с базой данных  !!ВАЖНО ?parseTime=true  добисывать в конце если надо чтобы выводило время
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.LogMode(true)
	db.AutoMigrate(&UserGorm{}) // автомотически отправляет данные в базу даных

	var u UserGorm
	// newDB := db.Where("email=?", "blah@blah.com")
	// newDB = newDB.Or("color=?", "wite")
	// newDB = newDB.First(&u)
	db = db.Where("email = ?", "bipus@gma1il.com").First(&u)
	if err := db.Where("email = ?", "bipus1@gmail.com").First(&u).Error; err != nil {
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
