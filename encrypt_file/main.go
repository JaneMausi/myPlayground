package main

import (
	"bufio"
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

	readBuff := make([]byte, 1024)
	readLen := 0
	for {
		n, err := reader.Read(readBuff)
		if err != nil && err != io.EOF {
			Fatal(err)
		}
		if n == 0 {
			break
		}
		readLen += n
	}

	// read password from console
	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter password: ")
	key, err := stdinReader.ReadBytes('\n')
	if err != nil {
		Fatal(err)
	}
	//fmt.Printf("\nkey: %s \n", key)

	// encrypt file
	//fmt.Printf("\nLength of original data: %d\n", readLen)
	result, err := AES.AesEncrypt(readBuff[:readLen], key[:len(key)-1])
	if err != nil {
		Fatal(err)
	}
	//fmt.Printf("\nResult data(%d):\n %s \n", len(result), base64.StdEncoding.EncodeToString(result))

	origData, err := AES.AesDecrypt(result, key[:len(key)-1])
	if err != nil {
		Fatal(err)
	}
	fmt.Printf("\nDecrypted data(%d):\n %s \n", len(origData), string(origData))

	// write encrypted file
	writeLen, err := writer.Write(result)
	if err != nil {
		Fatal(err)
	}
	fmt.Printf("\nWrite encrypted data to file (%d) \n", writeLen)
}
