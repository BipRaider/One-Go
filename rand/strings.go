package rand

/// для генераций произвольной строики
import (
	"crypto/rand"
	"encoding/base64"
)

const RemeberTokenBytes = 32

//Bytes will help us generate n random bytes, or will
//return an error if there was one. This uses the crypto/rand
// package so it safe to use with things like remeber tokens.
func Bytes(n int) ([]byte, error) { // генерирует произвольные числа  в заданом количестве (n)
	b := make([]byte, n)   // создает срез байтов этой длины.
	_, err := rand.Read(b) // fункцию Read из пакета crypto/rand,проверяет наличие ошибок и возвращает срез байтов, если ошибок нет.
	if err != nil {
		return nil, err
	}
	return b, nil
}

//NBytes  return the nember of bytes used the base64
//URL encoded string.
func NBytes(bese64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(bese64String) //возвращает количество байтов, использованных в кодированной строке base64URL
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

//String will generate a byte  slise   of size  nBytes and then
//return a string that is the base64 URL encoded version
//of that byte slice.
func String(nBytes int) (string, error) { //c генерируемый рэндом токен преобразуется в строку
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil // приоброзовали срез байтов в строку
}

//RememberToken is a helper function designed to generate
//remember tokens of a predeterined byte size.
func RememberToken() (string, error) { //c генерируется рэндом токен с количеством символов что мы в ведём
	return String(RemeberTokenBytes)
}
