package main

import (
	AES "gitlab.com/MXCFoundation/util/aes_encryption"
	"bufio"
	"encoding/base64"
	"fmt"
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
	rbuff := bufio.NewReader(reader)

	writer, err := os.OpenFile(encryptedFile,
		os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		Fatal(err)
	}
	wbuff := bufio.NewWriter(writer)

	defer func() {
		if err := reader.Close(); err != nil {
			Fatal(err)
		}
		if err := writer.Close(); err != nil {
			Fatal(err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := rbuff.Read(buf)
		if err != nil && err != io.EOF {
			Fatal(err)
		}
		if n == 0 {
			break
		}
	}
	// read password from console
	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter password: ")
	key, err := stdinReader.ReadBytes('\n')
	if err != nil {
		Fatal(err)
	}

	fmt.Printf("\nkey: %s \n", key)
	result, err := AES.AesEncrypt(buf, key[:len(key)-1])
	if err != nil {
		Fatal(err)
	}
	fmt.Printf("result:\n %s \n",base64.StdEncoding.EncodeToString(result))
	_, err = wbuff.Write(result)
	if err != nil {
		Fatal(err)
	}
	if err = wbuff.Flush(); err != nil {
		Fatal(err)
	}

	origData, err := AES.AesDecrypt(result, key[:len(key)-1])
	if err != nil {
		Fatal(err)
	}
	fmt.Printf("\nDecrypted data:\n %s \n", string(origData))

}