package main

import (
	"bufio"
	"fmt"
	AES "gitlab.com/MXCFoundation/payments-service/util/aes_encryption"
	"io"
	"os"
)

func readDynamicSize(initSize int, reader *os.File, readBuff []byte) ([]byte, error) {
	readLen := 0
	tmpBuff := make([]byte, initSize)

	for {
		n, err := reader.ReadAt(tmpBuff, int64(len(readBuff)))
		if err != nil && err != io.EOF {
			Fatal(err)
		}

		if err == io.EOF {
			return append(readBuff, tmpBuff...), err
		}

		readLen += n
		if readLen >= len(tmpBuff) {
			readBuff = append(readBuff, tmpBuff...)
			readBuff, err = readDynamicSize(initSize, reader, readBuff)
			if err == io.EOF {
				break
			}
		}

	}

	return readBuff, nil
}

func main() {
	originalFile := "./configuration/payment_service.toml"
	encryptedFile := "./configuration/payment_service.toml.enc"

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

	readBuff, _ := readDynamicSize(1024, reader, []byte(""))

	// read password from console
	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter password: ")
	key, err := stdinReader.ReadBytes('\n')
	if err != nil {
		Fatal(err)
	}
	//fmt.Printf("\nkey: %s \n", key)

	// encrypt file
	fmt.Printf("\nLength of readBuff: %d\n", len(readBuff))
	result, err := AES.AesEncrypt(readBuff[:len(readBuff)-1], key[:len(key)-1])
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
