package main

import (
	"errors"
	"fmt"
	"regexp"
)

func main() {

	r, _ := regexp.Compile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`) //true
	match1 := "f@f.com"
	fmt.Println(r.MatchString(match1))
	///----------------------------
	r = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	if !r.MatchString("p@s") {
		fmt.Println(errors.New("неправельно заполненый адрес"))
	}
	if r.MatchString("p@s.com") {
		fmt.Println(errors.New("правельно заполненый адрес"))
	}

	//----------------------------
	match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`, "mail@mai") //-false
	fmt.Println(match)

}
