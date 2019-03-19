package main

import (
	"bufio"
	"bytes"
	"fmt"
	AES "gitlab.com/MXCFoundation/util/aes_encryption"
	"io"
	"os"
)

func main() {
	// deal with input
	if len(os.Args) != 3 {
		Fatal("Invalid input params.")
	}
	originalFile := os.Args[1]
	encryptedFile := os.Args[2]

	reader, err := os.OpenFile(originalFile, os.O_RDONLY, 0644)
	if err != nil {
		Fatal(err)
	}

	writer, err := os.OpenFile(encryptedFile,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		Fatal(err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			Fatal(err)
		}
		if err := writer.Close(); err != nil {
			Fatal(err)
		}
	}()

	readBuff := bytes.NewBuffer(make([]byte, 1024))
	if _, err = readBuff.ReadFrom(reader); err != nil && err != io.EOF {
		Fatal(err)
	}

	// read password from console
	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter password: ")
	key, err := stdinReader.ReadBytes('\n')
	if err != nil {
		Fatal(err)
	}
	// encrypt file
	fmt.Printf("\nkey: %s \n", key)
	result, err := AES.AesEncrypt(readBuff.Bytes(), key[:len(key)-1])
	if err != nil {
		Fatal(err)
	}
	fmt.Printf("\nResult data:\n %s \n", string(result))

	origData, err := AES.AesDecrypt(result, key[:len(key)-1])
	if err != nil {
		Fatal(err)
	}
	fmt.Printf("\nDecrypted data:\n %s \n", string(origData))

	// write encrypted file
	writeBuff := bytes.NewBuffer(result)
	if _, err = writeBuff.WriteTo(writer); err != nil {
		Fatal(err)
	}

}
