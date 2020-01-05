package main

import "C"
import (
	"bufio"
	"bytes"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sync"
	"testing"
)

func main() {

	//fmt.Println("hello, you hillbilly!")

	// 1. Convert a image into a byte array
	var byteBuf = jpgToByte("WhatsAppImage2020-01-05at17.49.59.jpeg")
	//fmt.Println(byteBuf)

	// 2. Convert that byte array back to a image
	// --> Did it work??
	byteToJpg(byteBuf, "TheVeryNewFuckingFile.jpg")

	// 3. Erasure Code your byte array
	var encodedByteBuf = erasureCodeBytes(byteBuf)

}

func erasureCodeBytes (byteBuf []byte) (encodedByteBuf []byte) {

	code := intoCode(byteBuf)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 2, 3, 4}

	corrupted := corrupt(append(source, encoded...), errList, shardLength)

	recovered := code.Decode(corrupted, errList, false)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}

func intoCode(byteBuf []byte) *Code {
	m := float64(len(byteBuf))
	rate := 0.75 // how many percent of all shards do we need to reconstruct the original data
	k := math.Round(rate * m)
	shardLength := 1.25 * m
	size := k * shardLength

	if m <= 0 || k <= 0 || k >= m || k > 127 || m > 127 || size < 0 {
		panic("Invalid erasure code params")
	}
	if size%k != 0 {
		panic("Size to encode is not divisable by k and therefore cannot be encoded into shards")
	}

	encodeMatrix := make([]byte, m*k)
	galoisTables := make([]byte, k*(m-k)*32)

	if k > 5 {
		C.gf_gen_cauchy1_matrix((*C.uchar)(&encodeMatrix[0]), C.int(m), C.int(k))
	} else {
		C.gf_gen_rs_matrix((*C.uchar)(&encodeMatrix[0]), C.int(m), C.int(k))
	}

	C.ec_init_tables(C.int(k), C.int(m-k), (*C.uchar)(&encodeMatrix[k*k]), (*C.uchar)(&galoisTables[0]))
	return &Code{
		M:            m,
		K:            k,
		ShardLength:  size / k,
		EncodeMatrix: encodeMatrix,
		galoisTables: galoisTables,
		decode: &decodeNode{
			children: make([]*decodeNode, m),
			mutex:    &sync.Mutex{},
		},
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
