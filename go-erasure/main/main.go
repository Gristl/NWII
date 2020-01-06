package main
// http://drmingdrmer.github.io/tech/distributed/2017/02/01/ec.html
import (
	"C"
	"bufio"
	"bytes"
	"fmt"
	//origin "github.com/Gristl/NWII/go-erasure"
	"io/ioutil"
	origin "main/originalCode"
	"os"
)

func main() {
	// 1. Convert a image into a byte array
	var byteBuf = jpgToByte("Image.jpeg")
	fmt.Println("This is the picture as as byte array: " )
	fmt.Println(byteBuf)

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

	/*source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}*/

	encoded := code.Encode(byteBuf)


	/* ToDo Nice To have: generate random errors in the variabel errList
	var errList []byte

	for i := 0; i < k; i++ {
		errList = append(errList, byte(rand.Intn(len(encoded))))
	}
	*/
	errList := []byte{0, 2, 3, 4} // ToDo Tell me why we can make at least 4 erros --> try with more, does not work

	// The function corrupt deletes the bytes at the indexes that are handed over in the errList and all the following #sharedLength bytes
	// That means that the i' th share is completely deleted
	corrupted := origin.Corrupt(append(byteBuf, encoded...), errList, shardLength)

	stillCorrupt := code.DecodeWOMagic(corrupted, errList, false)
	recovered := code.Decode(corrupted, errList, false)

	// Delete the last bytes in the array that we added so that size%k == 0
	for toBeDeletedAtTheEnd > 0 {
		byteBuf = byteBuf[:len(byteBuf)-1]
		stillCorrupt = corrupted[:len(corrupted)-1]
		recovered = recovered[:len(recovered)-1]
		toBeDeletedAtTheEnd--
	}
	byteToJpg(stillCorrupt, "ImageCorrupt.jpg")
	byteToJpg(recovered, "ImageRecovered.jpg")

	fmt.Println("This is the length of the byte array: %v", len(byteBuf))
	fmt.Println("This is the length of the encoded byte array: %v", len(encoded))
	fmt.Println("This is the length of the currupt encoded byte array: %v", len(corrupted))
	fmt.Println("This is the length of the currupt decoded array: %v", len(stillCorrupt))
	fmt.Println("This is the length of the recovered decoded byte array: %v", len(recovered))


	if !bytes.Equal(byteBuf, recovered) {
		fmt.Println("Source was not successfully recovered with 4 errors")
	} else {
		fmt.Println("This was sooooo incredibly successful! ")
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


func check(err error) {
	if err != nil {
		fmt.Println(err)
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