package main

import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	origin "github.com/Gristl/NWII/go-erasure/main"
	"io/ioutil"
	"math/rand"
	"os"
)

func main() {
	// 1. Convert a image into a byte array
	var byteBuf = jpgToByte("Image.jpeg")
	//fmt.Println(byteBuf)

	// 2. Convert that byte array back to a image
	// --> Did it work??
	byteToJpg(byteBuf, "ImageDuplicate.jpg")

	// 3. Erasure Code your byte array

	// Make sure that (size%k == 0)
	k := 8
	var a byte = 0
	var toBeDeletedAtTheEnd = 0

	for len(byteBuf)%k != 0 {
		byteBuf = append(byteBuf, a)
		toBeDeletedAtTheEnd++
	}

	size := len(byteBuf) //k * shardLength
	shardLength := size / k
	m := 12

	code := origin.NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 2, 3, 4}

	corrupted := corrupt(append(byteBuf, encoded...), errList, shardLength)

	recovered := code.Decode(corrupted, errList, false)
	for toBeDeletedAtTheEnd > 0 {
		//delete(recovered, byte)
	}
	byteToJpg(corrupted, "ImageCorrupt.jpg")
	byteToJpg(recovered, "ImageRecovered.jpg")
	if !bytes.Equal(byteBuf, recovered) {
		fmt.Println("Source was not successfully recovered with 4 errors")
	}
}

func byteToJpg (byteBuf []byte, imageName string) {
	err := ioutil.WriteFile(imageName, byteBuf, 0644)
	check(err)
}

// source: https://socketloop.com/tutorials/golang-convert-an-image-file-to-byte
func jpgToByte (imageName string) (byteBuf []byte) {
	file, err := os.Open(imageName)

	check(err)

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	byteBuf = make([]byte, size)

	// read file into bytes
	var buffer = bufio.NewReader(file)
	_, err = buffer.Read(byteBuf)

	check(err)
	// then we need to determine the file type
	// see https://www.socketloop.com/tutorials/golang-how-to-verify-uploaded-file-is-image-or-allowed-file-types
	//filetype := http.DetectContentType(bytes)
	//fmt.Println(filetype)
	return byteBuf
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
func corrupt(source, errList []byte, shardLength int) []byte {
	corrupted := make([]byte, len(source))
	copy(corrupted, source)
	for _, err := range errList {
		for i := 0; i < shardLength; i++ {
			corrupted[int(err)*shardLength+i] = 0x00
		}
	}
	return corrupted
}*/