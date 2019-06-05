package models

import (
	"log"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func testingUserService() (*UserService, error) {
	const (
		mysqlinfo_test = "root:alfadog1@/bipusdb_test?charset=utf8&parseTime=True&loc=Local"
	)

	us, err := NewUserService(mysqlinfo_test)
	if err != nil {
		log.Println(err)
	}
	defer us.Close()
	us.db.LogMode(false)
	//clear the users table between tests
	us.DestructiveReset()
	return us, nil
}
func TestCreateUder(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Name:  "Micheal Scott",
		Email: "Micheal@gmail.com",
	}

	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == 0 {
		t.Errorf("Expected ID > 0. Received %d", user.ID)
	}
	if time.Since(user.CreatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected CreateAt to be recent.Received %s", user.CreatedAt)
	}
	if time.Since(user.UpdatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected UpdatedAt to be recent.Received %s", user.UpdatedAt)
	}
}
