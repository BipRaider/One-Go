package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // подсоединить библиотеку mysql
)

type UserGorm struct {
	gorm.Model
	Name   string
	Email  string `gorm:"not null;unique_index"` // добовления типа данных в базуданых
	Color  string
	Orders []Orders `gorm:"ForeignKey:UserId"` // `gorm:"ForeignKey:RecordId"` надо указать для обьядинения библиотек
	// чтобы можно было их присоединять  через  db.Preload() http://doc.gorm.io/crud.html#preloading-eager-loading
}
type Orders struct {
	gorm.Model
	UserID      uint
	Amount      int    //количество
	Description string // описание
}

func main() {

	db, err := gorm.Open("mysql", "root:alfadog1@/bipusdb?charset=utf8&parseTime=True&loc=Local") //Соединение с базой данных  !!ВАЖНО ?parseTime=true  добисывать в конце если надо чтобы выводило время
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.LogMode(true)                       // отображает лог строки запроса в терменале
	db.AutoMigrate(&UserGorm{}, &Orders{}) // автомотически отправляет данные в базу даных

	var u UserGorm

	if err := db.Preload("Orders").Find(&u).Error; err != nil {
		os.Exit(3)
	}
	fmt.Println(u)

	// createOrder(db, u, 1001, "Fake description #4")
	// createOrder(db, u, 231, "Fake description #5")
	// createOrder(db, u, 101, "Fake description #6")

}
func createOrder(db *gorm.DB, user UserGorm, amount int, dest string) {
	err := db.Create(&Orders{
		UserID:      user.ID,
		Amount:      amount,
		Description: dest,
	}).Error
	if err != nil {
		os.Exit(3)
	}
}
