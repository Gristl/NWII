package erasure

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	
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
	bytes := make([]byte, size)

	// read file into bytes
	var buffer = bufio.NewReader(file)
	_, err = buffer.Read(bytes)    // <--------------- here!

	// then we need to determine the file type
	// see https://www.socketloop.com/tutorials/golang-how-to-verify-uploaded-file-is-image-or-allowed-file-types

	filetype := http.DetectContentType(bytes)

	err = bucket.Put(path, bytes, filetype, s3.ACL("public-read"))
}