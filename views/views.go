package views

import (
	"html/template"
	"os"
	"path/filepath"
)

type View struct {
	Template *template.Template
	Layout   string
}

var (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	//Добовляем данные из файла
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...) // Записываем в переменую t значаения(срез  файлов ) с  файла по сылке

	if err != nil {

		os.Exit(5)

	}
	return &View{
		Template: t, //Присваиваем шаблону   новый шаблон t  и возвращаем даные
		Layout:   layout,
	}
}

// this func переберает  все files ,what есть в папке LayoutDir
// (*)- говорит что все файлы должны соотвецтвовать макету (типу файлу ,что указано )
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		os.Exit(7)
	}
	return files
}

type NotView struct {
	Template *template.Template
	Layout   string
}

//Template errors
func NotFound() *NotView {

	Name := "bootstrap"

	files := []string{"views/nodfound.gohtml"}
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)

	if err != nil {
		os.Exit(6)
	}
	return &NotView{
		Template: t,
		Layout:   Name,
	}
}
