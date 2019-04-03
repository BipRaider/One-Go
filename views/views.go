package views

import (
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
	addTemplatePath(files)
	addTemlateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...) // Записываем в переменую t значаения(срез  файлов ) с  файла по сылке

	if err != nil {
		panic(err)
	}
	return &View{
		Template: t, //Присваиваем шаблону   новый шаблон t  и возвращаем даные
		Layout:   layout,
	}
}

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// this func переберает  все files ,what есть в папке LayoutDir  1.1.2
// (*)- говорит что все файлы должны соотвецтвовать макету (типу файлу ,что указано )
func layoutFiles() []string { //Добовляем данные из файла

	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

//1.1.3
// addTemplatePath  takes in a slice of strings
// representing file paths for templates and it prepends // представляет файл пути за шаблоном и это добовляет
//the TamplateDir directory to each string in the slice  //каталог TamplateDir для каждой строки в срезе
//
//Eg the input {"home"}would result in the output
//{"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) { //Добавит "views/"
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

//1.1.4
//addTemplateExt  takes in a slice of strings
//representing file paths for templates and it prepends
//the TamplateExt extension to each string in the slice
//
//Eg the input {"home"}would result in the output
//{"home.gohmtl"} if TemplateExt ==".gohtml"
func addTemlateExt(files []string) { //Добавит ".gohml"
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

// Функция  которая выводит в браузер нужную файл и какого шаблона , 1.2
//Render is used to render the viewwith the predefind layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // надо указывать кодировку ;charset=utf-8
	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data{
			Yield: data,
		}
	}

	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

//Функция  которая выводит в браузер нужную файл и какого шаблона , 1.2.1
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

//------------------------------------------------------------
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
