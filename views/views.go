package views

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type View struct {
	Template *template.Template
	Layout   string
}

// функция для обьявления какой шаблон использовать и хранящийся  файл (чтобы потом выводить в браузере)   1.1
func NewView(layout string, files ...string) *View {
	//Добовляем данные из файла
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...) // Записываем в переменую t значаения(срез  файлов ) с  файла по сылке

	if err != nil {
		err = errors.New("ERROR func views/NewView at ParseFiles")
		os.Exit(5)
	}
	return &View{
		Template: t, //Присваиваем шаблону   новый шаблон t  и возвращаем даные
		Layout:   layout,
	}
}

var (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
)

// this func переберает  все files ,what есть в папке LayoutDir
// (*)- говорит что все файлы должны соотвецтвовать макету (типу файлу ,что указано )
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		os.Exit(7)
	}
	return files
}

// Функция для которая выводит в браузер нужную файл и какого шаблона , 1.2
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)

}

//Template errors
func NotFound() *View {

	Name := "bootstrap"

	files := []string{"views/nodfound.gohtml"}
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)

	if err != nil {
		os.Exit(6)
	}
	return &View{
		Template: t,
		Layout:   Name,
	}
}
