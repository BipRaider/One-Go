package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvSuccess  = "success"
)

//Alert is used to render Bootstrap Alert messager in the
//bootstrap.gohtml template
type Alert struct {
	Level   string
	Message string
}

// Data is the top leval structure that views expect data
type Data struct {
	Alert *Alert
	Yield interface{}
}

// a := Alert{
// 	Level:   "warning",                                // можно вписать такие имена  класа success,info,danger,warning
// 	Message: "Successfully rendered a dynamic alert!", // выводин на экран данное сообщение
// }
// d := Data{
// 	Alert: a,
// 	Yield: "Hello!",
// }
