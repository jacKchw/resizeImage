package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
)

//this code rename all image files to number
//for window
//GOOS=windows go build renameIMG.go
//mv renameIMG.exe /mnt/c/Users/aesop

var imagesExt = []string{".png", ".jpg", ".jpeg"}
var width = 1072
var height = 1448

var wg sync.WaitGroup
var files = make(chan string, 5000)

func Include(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func walkFc(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(1)
	files <- path
	return nil
}

func isImageWorker(imgFiles chan string) {
	for fileName := range files {
		go isImage(fileName, imgFiles)
	}
}

func isImage(fileName string, imgFiles chan string) {
	ext := filepath.Ext(fileName)
	_, isImage := Include(imagesExt, ext)
	if isImage {
		imgFiles <- fileName
	} else {
		wg.Done()
	}
}

func resizeWorker(imgFiles chan string) {
	for path := range imgFiles {
		resizeImage(path)
		fmt.Println(path)
		wg.Done()
	}
}

func resizeImage(path string) {
	// Open a test image.
	src, err := imaging.Open(path)
	if err != nil {
		log.Printf("error opening file: %v, %v", err, path)
		return
	}
	// Create a blurred version of the image.
	src = imaging.Fit(src, width, height, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err = imaging.Save(src, path)
	if err != nil {
		log.Printf("error saving file: %v, %v", err, path)
	}
}

func main() {
	//write log to testlogfile
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, err := filepath.Abs("/Users/jack/dev/jackchw/resizeImage")
	fmt.Println(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dir)

	var imgFiles = make(chan string, 5000)
	go isImageWorker(imgFiles)
	go resizeWorker(imgFiles)

	err = filepath.Walk(dir, walkFc)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
