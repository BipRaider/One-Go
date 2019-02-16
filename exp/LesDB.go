package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // подсоединить библиотеку mysql
)

type UserGorm struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Color string
}

func getInfp() (name, email, color string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What id your name?")
	name, _ = reader.ReadString('\n')
	fmt.Println("Whatewr idewr yourer nameqweqwe?")
	email, _ = reader.ReadString('\n')
	fmt.Println("Whatewr color?")
	color, _ = reader.ReadString('\n')
	name = strings.TrimSpace(color)
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	return name, email, color
}

func main() {
	db, err := gorm.Open("mysql", "root:alfadog1@/bipusdb") //Соединение с базой данных

	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	if err := db.DB().Ping(); err != nil {
		os.Exit(133)
	}
	db.LogMode(true)
	db.AutoMigrate(&UserGorm{}) // автомотически отправляет данные в базу даных
	name, email, color := getInfp()
	u := UserGorm{
		Name:  name,
		Email: email,
		Color: color,
	}
	if err = db.Create(&u).Error; err != nil {
		os.Exit(1)
	}
	fmt.Printf("%+v\n", u)
}
