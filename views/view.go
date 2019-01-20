package views

import (
	"html/template"
	"os"
)

func NewView(files ...string) *View {

	files = append(files, "views/footer/footer.gohtml")
	t, err := template.ParseFiles(files...)

	if err != nil {
		os.Exit(5)
		panic(err)
	}
	return &View{
		Template: t,
	}
}

type View struct {
	Template *template.Template
}
