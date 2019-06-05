package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Image is NOT stored in the database
type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	return "/" + i.RelativePath()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

//----------------------------------------------------------------
type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

//--------------------------------------------------------------------------------

type imageService struct{}

func NewImageService() ImageService {
	return &imageService{}
}

///-----------------------------------------
func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}

	//Create a destination file
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close() // в конце закрываем открытый файл

	//Copy reader  data to the destination
	_, err = io.Copy(dst, r) // копируем данные в dst
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	pathImage := is.imagePath(galleryID)
	files, err := filepath.Glob(pathImage + "*")
	if err != nil {
		panic(err)
	}
	var filesPath []string

	for _, str := range files {
		str = strings.Replace(str, "\\", "/", -1) // Функция для замены символов в строке ( поменял \ на / из-за Windows меняет направление слеш)
		filesPath = append(filesPath, str)
	}

	ret := make([]Image, len(filesPath))
	for i := range filesPath {
		filesPath[i] = strings.Replace(filesPath[i], pathImage, "", 1)
		ret[i] = Image{
			Filename:  filesPath[i],
			GalleryID: galleryID,
		}
	}

	return ret, nil
}
func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

//Create the directory to contain our images
func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	// Функция MkdirAll создаст любую из этих папок, которые не существуют.
	// Если вы вообще знакомы с Linux, это похоже на использование mkdir -p.
	// MkdirAll также требует FileMode, который в основном является просто набором разрешений
	// для каталога. Мы будем использовать 0755, который дает текущему пользователю разрешение
	// что-либо делать с файлом, а другим разрешает чтение и выполнение.
	galleryPath := is.imagePath(galleryID) //присваеваем пусть(строку)  с именем галери и айди
	err := os.MkdirAll(galleryPath, 0755)  //MkdirAll создает каталог с именем пути вместе со всеми необходимыми родителями и возвращает nil
	if err != nil {
		return "", err
	}
	//galleryPath := filepath.Join("images", "galleries",fmt.Sprintf(gallery.ID))
	// Мы используем filepath.Join вместо построения пути вручную,
	// потому что косые черты и другие символы могут различаться в разных операционных системах.
	return galleryPath, nil
}

//Страничка 578
