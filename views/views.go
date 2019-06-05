package views

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"../context"
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
		log.Println(err)
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
		log.Println(err)
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

//Функция  которая выводит в браузер нужную файл и какого шаблона , 1.2.1
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Функция  которая выводит в браузер нужную файл и какого шаблона , 1.2
//Render is used to render the viewwith the predefind layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // надо указывать кодировку ;charset=utf-8
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
		// do nothing
	default:
		vd = Data{
			Yield: data,
		}
	}

	vd.User = context.User(r.Context())

	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. If the problem  persists , please email support@thebipus.com ", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

//------------------------------------------------------------
//Template errors
func NotFound() *View {
	Name := "bootstrap"
	files := []string{"views/nodfound.gohtml"}
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
	}
	return &View{
		Template: t,
		Layout:   Name,
	}
}
