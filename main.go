package main

import (
	"fmt"
	"github.com/signintech/gopdf"
	"github.com/thrawn01/args"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Config struct {
	Random bool
}

const PATH = "result"

var confg Config

func main() {
	parser := args.NewParser()
	parser.AddOption("--random").Alias("-r").IsTrue().StoreTrue(&confg.Random).Help("Random")
	parser.ParseOrExit(nil)
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeLetter})

	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	filePath, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	fontPath := filepath.Dir(filePath)

	os.Mkdir(root+"\\"+PATH, 0666)
	files, err := ioutil.ReadDir(root)
	check(err)
	legalFiles := []os.FileInfo{}
	for _, file := range files {
		if !file.IsDir() {
			fileType := GetFileType(file.Name())
			if fileType == "image/jpeg" || fileType == "image/png" || fileType == "image/jpg" || fileType == "image/bmp" || fileType == "image/gif" || fileType == "image/x-bmp" {
				legalFiles = append(legalFiles, file)
			}
		}
	}
	if confg.Random {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(legalFiles), func(i, j int) { legalFiles[i], legalFiles[j] = legalFiles[j], legalFiles[i] })
	}
	for _, file := range legalFiles {
		fmt.Println(root + "\\" + file.Name())
		exec.Command("magick", root+"\\"+file.Name(), root+"\\"+PATH+"\\"+file.Name()+".ppm").Output()
		exec.Command("potrace", root+"\\"+PATH+"\\"+file.Name()+".ppm", "-b", "pdfpage", "-o", root+"\\"+PATH+"\\"+file.Name()+".pdf").Output()
		os.Remove(root + "\\" + PATH + "\\" + file.Name() + ".ppm")
		pdf.AddPage()
		tpl1 := pdf.ImportPage(root+"\\"+PATH+"\\"+file.Name()+".pdf", 1, "/MediaBox")
		pdf.UseImportedTemplate(tpl1, 0, 0, 0, 0)
		pdf.AddPage()
		err := pdf.AddTTFFont("Roboto", fontPath+"\\Roboto-Regular.ttf")
		if err != nil {
			log.Print(err.Error())
			return
		}
		fontSize := 2
		err = pdf.SetFont("Roboto", "", fontSize)
		if err != nil {
			log.Print(err.Error())
			return
		}
		//pdf.SetGrayFill(0.5)
		//pdf.Cell(nil, "Áa")

		//Measure Width
		text := " "
		pdf.Cell(nil, text)

	}
	pdf.WritePdf("pdk-result.pdf")
	os.RemoveAll(root + "\\" + PATH)
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
