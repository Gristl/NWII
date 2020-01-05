package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	fmt.Println("hello, you hillbilly!")
	var byteBuf= jpgToByte("WhatsAppImage2020-01-05at17.49.59.jpeg")
	fmt.Println(byteBuf)
	byteToJpg(byteBuf, "TheVeryNewFuckingFile.jpg")
}

func byteToJpg (imageBuf []byte, imageName string) {
	err := ioutil.WriteFile(imageName, imageBuf, 0644)
	check(err)
}

// source: https://socketloop.com/tutorials/golang-convert-an-image-file-to-byte
func jpgToByte (imageName string) (imageBuf []byte) {
	file, err := os.Open(imageName)

	check(err)

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	imageBuf = make([]byte, size)

	// read file into bytes
	var buffer = bufio.NewReader(file)
	_, err = buffer.Read(imageBuf)

	check(err)
	// then we need to determine the file type
	// see https://www.socketloop.com/tutorials/golang-how-to-verify-uploaded-file-is-image-or-allowed-file-types
	//filetype := http.DetectContentType(bytes)
	//fmt.Println(filetype)
	return imageBuf
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}
