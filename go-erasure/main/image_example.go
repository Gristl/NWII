package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	fmt.Println("hello, you hillbilly!")
	fmt.Println(jpgIntoByte("WhatsAppImage2020-01-05at17.49.59.jpeg"))
}

// https://socketloop.com/tutorials/golang-convert-an-image-file-to-byte
func jpgIntoByte(imageName string) (imageBuf []byte) {
	file, err := os.Open(imageName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	imageBuf = make([]byte, size)

	// read file into bytes
	var buffer = bufio.NewReader(file)
	_, err = buffer.Read(imageBuf)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// then we need to determine the file type
	// see https://www.socketloop.com/tutorials/golang-how-to-verify-uploaded-file-is-image-or-allowed-file-types
	//filetype := http.DetectContentType(bytes)
	//fmt.Println(filetype)

	return imageBuf
}