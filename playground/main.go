package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func main() {
	InitLogWriter()

	reader, err := os.OpenFile("payment_service.conf", os.O_RDONLY, 0644)
	if err != nil {
		Fatal(err)
	}
	rbuff := bufio.NewReader(reader)

	writer, err := os.OpenFile("payment_service.conf.enc",
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
	read_len := 0
	for {
		read_len, err = rbuff.Read(buf)
		if err != nil && err != io.EOF {
			Fatal(err)
		}
		if read_len == 0 {
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

	fmt.Printf("key: %s", key)
	result, err := AesEncrypt(buf, key[:len(key)-1])
	if err != nil {
		Fatal(err)
	}
	fmt.Printf("result: %s",base64.StdEncoding.EncodeToString(result))
	write_len, err := wbuff.Write(result[:read_len])
	if write_len != read_len || err != nil {
		Fatal(err)
	}

	origData, err := AesDecrypt(result, key[:len(key)-1])
	if err != nil {
		Fatal(err)
	}
	fmt.Println(string(origData))

}
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
