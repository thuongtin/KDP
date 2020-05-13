package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const PATH = "result"

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	os.Mkdir(root+"\\"+PATH, 0666)
	files, err := ioutil.ReadDir(root)
	check(err)
	for _, file := range files {
		if !file.IsDir() {
			fileType := GetFileType(file.Name())
			if fileType == "image/jpeg" || fileType == "image/png" || fileType == "image/jpg" || fileType == "image/bmp" || fileType == "image/x-bmp" {
				fmt.Println(root + "\\" + file.Name())
				exec.Command("magick", root+"\\"+file.Name(), root+"\\"+PATH+"\\"+file.Name()+".ppm").Output()
				exec.Command("potrace", root+"\\"+PATH+"\\"+file.Name()+".ppm", "-b", "pdfpage", "-o", root+"\\"+PATH+"\\"+file.Name()+".pdf").Output()
				os.Remove(root + "\\" + PATH + "\\" + file.Name() + ".ppm")
			}
		}
	}
}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetFileType(filePath string) string {
	f, err := os.Open(filePath)
	if err == nil {
		defer f.Close()
		contentType, err := GetFileContentType(f)
		if err != nil {
			panic(err)
		}
		return contentType
	}
	return ""
}

func GetFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}