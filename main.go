package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	// user : select encode or decode
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Decode or Encode (D/E) ?")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response: ", err)
		return
	}

	var opFunc func(string)
	var opStr string

	switch choice[0] {
	case 'e', 'E':
		opFunc = encode
		opStr = "encode"
	case 'd', 'D':
		opFunc = decode
		opStr = "decode"
	default:
		fmt.Println("whatever.")
		return
	}

	// user : enter a file name
	fmt.Printf("File name to %s :", opStr)
	fileName, err := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	if err != nil {
		fmt.Println("Failed to read file name:", err)
		return
	}

	switch choice[0] {
	case 'e', 'E':
		opFunc(fileName)
	default:
		opFunc(fileName)
	}
}

func encode(fileName string) {

	// c:\temp\layout.json
	// c:\temp\layout.json.encoded

	fmt.Println("Encoding", fileName)

	// read the file in
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read file: ", err)
		return
	}

	// compress it
	buf := new(bytes.Buffer)

	gz := gzip.NewWriter(buf)
	_, err = gz.Write(fileContent)
	if err != nil {
		fmt.Println("Failed to zip the file: ", err)
		return
	}
	gz.Close()

	// base64 encode it
	result := base64.RawStdEncoding.EncodeToString(buf.Bytes())

	// output it
	os.Stdout.Write([]byte(result))
	os.WriteFile(fileName+".encoded", []byte(result), os.FileMode(os.O_RDWR|os.O_CREATE|os.O_TRUNC))
}

func decode(fileName string) {

	// c:\temp\layout.json
	// c:\temp\layout.json.encoded

	fmt.Println("Decoding", fileName)

	// read the file in
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read file: ", err)
		return
	}

	// base64 decode it
	_, err = base64.RawStdEncoding.Decode(fileContent, fileContent)
	if err != nil {
		fmt.Println("Failed to base64 decode file: ", err)
		return
	}

	// decompress it
	reader := bytes.NewReader([]byte(fileContent))
	gzreader, err := gzip.NewReader(reader)
	if err != nil {
		fmt.Println("Failed to decompress file [stage 1]: ", err)
		return
	}

	result, err := ioutil.ReadAll(gzreader)
	if err != nil {
		fmt.Println("Failed to decompress file [stage 2]: ", err)
	}

	// output it
	os.Stdout.Write(result)
	os.WriteFile(fileName+".decoded", []byte(result), os.FileMode(os.O_RDWR|os.O_CREATE|os.O_TRUNC))
}
