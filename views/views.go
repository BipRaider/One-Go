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
	//Добовляем данные из вайла
	files = append(files,
		"views/layouts/bootstrap.gohtml",
		"views/layouts/footer.gohtml",
	)

	t, err := template.ParseFiles(files...) // Записываем в переменую t значаения(срез  файлов ) с  файла по сылке

	if err != nil {
		println(files)
		os.Exit(5)

	}
	return &View{
		Template: t, //Присваиваем шаблону   новый шаблон t  и возвращаем даные
		Layout:   layout,
	}
}
