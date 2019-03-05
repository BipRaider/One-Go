package main

import "fmt"

type Cat struct{}
type Dog struct {
	Gaf string
} //-------------------------Dog  в обеи случаех есть --------
type Ask struct {
	one int
	tow string
}

func (d Cat) Speak() {
	fmt.Println("gav")
}

func (d Dog) Speak(s string) { //-------------------------Dog  в обеи случаех есть --------
	d.Gaf = s
	fmt.Println(d.Gaf + s)
}

// type Husky struct {    //equal to h.Dog.Speak()
// 	Dog
// }

func (a *Dog) La(g int) { //-------------------------Dog  в обеи случаех есть --------

	fmt.Println(a.Gaf, g)
}
func (a *Ask) La(g int) {
	a = &Ask{
		one: g,
		tow: "2",
	}
	fmt.Println(a)
}

type Husky struct {
	Speaker
}
type SpeakerPrefer struct {
	Speaker
}

// если методы не индетичны запрашиваемого  интерфейсу  то они не будут работать
type Speaker interface { //-----------Соотвецтвует только Dog
	Speak(string) // имена методов/функций  не могут быть с одинаковыми именами и с разными типами(string,int ....)
	La(int)       // но type ... struct  могут с люым набором данных
}

func (sp SpeakerPrefer) Speak(g string) { // при использования  в функций\метотада : func (sp SpeakerPrefer) ,
	// можно обнавлять или добовлять или изменять данные что поступают функцию\метода :func (a *Dog)

	gg := "Dsfdggg - " + g   // вносим изменения
	fmt.Print("Gaga :" + gg) // выводим результат// данный результат будет будет отоброжаться в функций  что мы использовали
	sp.Speaker.Speak(gg)     /// дублирует сообщения //если не водить данную строку то будет вывод в терменале в одну строку
}
func (sp SpeakerPrefer) La(h int) {
	switch h {
	case 3:
		fmt.Print("La 3=:", h)
	case 5:
		fmt.Print("La 5:")
	default:
		fmt.Print("La :")
	}

	sp.Speaker.La(h) //если не водить данную строку то будет вывод в терменале в одну строку
}
func main() {

	h := Husky{SpeakerPrefer{&Dog{}}} // интерфейс удолитворяет  любой интерфейс с индентичными данными
	h.La(3)                           //equal to h.Speaker.Speak()
	h.Speak("sad")
	fmt.Println(h)
}
