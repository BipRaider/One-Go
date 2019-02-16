package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // подсоединить библиотеку mysql
)

type UserGorm struct {
	gorm.Model
	LastName string
	Email    string `gorm:"not null;unique_index"`
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

	// // Create
	// db.Create(&Product{Code: "L1212", Price: 1000})

	// // Read
	// var product Product
	// db.First(&product, 1)                   // find product with id 1
	// db.First(&product, "code = ?", "L1212") // find product with code l1212

	// // Update - update product's price to 2000
	// db.Model(&product).Update("Price", 2000)

	// // Delete - delete product
	// db.Delete(&product)
}
