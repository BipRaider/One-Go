package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]string, error)
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

func (is *imageService) ByGalleryID(galleryID uint) ([]string, error) {
	pathImage := is.imagePath(galleryID)
	files, err := filepath.Glob(pathImage + "*")
	if err != nil {
		panic(err)
	}
	var filesPath []string
	for _, str := range files {
		str = strings.Replace(str, "\\", "/", -1)
		filesPath = append(filesPath, str)
	}

	return filesPath, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

//Create the directory to contain our images
func (is *imageService) mkImagePath(galleryID uint) (string, error) {

	galleryPath := is.imagePath(galleryID) //Создаёт паку с именем галери айди
	err := os.MkdirAll(galleryPath, 0755)  //MkdirAll создает каталог с именем path(путь) вместе со всеми необходимыми родителями и возвращает nil
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
