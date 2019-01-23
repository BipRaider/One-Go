package views

import (
	"html/template"
	"os"
)

type View struct {
	Template *template.Template
	Layout   string
}

func NewView(layout string, files ...string) *View {
	//Добовляем данные из файла
	files = append(files,

		"views/layouts/bootstrap.gohtml",
		"views/layouts/navbar.gohtml",
		"views/layouts/footer.gohtml",
	)

	t, err := template.ParseFiles(files...) // Записываем в переменую t значаения(срез  файлов ) с  файла по сылке

	if err != nil {

		os.Exit(5)

	}
	return &View{
		Template: t, //Присваиваем шаблону   новый шаблон t  и возвращаем даные
		Layout:   layout,
	}
}

type NotView struct {
	Template *template.Template
	Layout   string
}

func NotFound() *NotView {

	Name := "bootstrap"

	files := []string{
		"views/nodfound.gohtml",
		"views/layouts/bootstrap.gohtml",
		"views/layouts/navbar.gohtml",
		"views/layouts/footer.gohtml",
	}
	t, err := template.ParseFiles(files...)

	if err != nil {
		os.Exit(6)
	}
	return &NotView{
		Template: t,
		Layout:   Name,
	}
}
