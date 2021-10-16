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

	// c:\temp\layout.json
	// c:\temp\layout.json.encoded

	// user : enter a file name or a string to encode or decode
	fmt.Print("Enter a file name or the text to encode/decode: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if err != nil {
		fmt.Println("Failed to read the input:", err)
		return
	}

	// First we attempt to read input from a file, if that fails then treat the input as the actual text to encode or decode
	var text []byte
	fileContent, err := os.ReadFile(input)
	if err != nil {
		fmt.Println("\nTreating input as text")
		text = []byte(input)
	} else {
		fmt.Print("\nTreated the input as a file name and read the contents:\n\n", string(fileContent), "\n\n")
		text = fileContent
	}

	// Try to decode it
	op := "decode"
	result, err := decode(text)
	if err != nil {
		// The only option left is to encode it!
		op = "encode"
		result, err = encode(text)
		if err != nil {
			fmt.Println("Completely failed to do anything useful with that input")
			return
		}
	}

	fmt.Print("success!\n\n")

	// output it
	os.Stdout.Write([]byte(result))
	fmt.Print("\n\n")

	if len(fileContent) > 0 {
		fileOut := input + "." + op + "d"
		fmt.Print("Output also written to " + fileOut + "\n\n")
		os.WriteFile(fileOut, []byte(result), os.FileMode(os.O_RDWR|os.O_CREATE|os.O_TRUNC))
	}
}

func decode(text []byte) (result string, err error) {

	fmt.Print("Attempting to decode the content... ")

	// base64 decode it
	textDecoded := make([]byte, len(text))
	_, err = base64.RawStdEncoding.Decode(textDecoded, text)
	if err != nil {
		fmt.Print(err)
		return
	}

	// decompress it
	reader := bytes.NewReader(textDecoded)
	gzreader, err := gzip.NewReader(reader)
	gzreader.Multistream(false)
	if err != nil {
		fmt.Print("error at stage 1: " + err.Error())
		return
	}

	resultBytes, err := ioutil.ReadAll(gzreader)
	if err != nil {
		fmt.Print("error at stage 2: " + err.Error())
	}

	result = string(resultBytes)
	return
}

func encode(text []byte) (result string, err error) {

	fmt.Print("\nAttempting to encode the content... ")

	// compress it
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)
	_, err = gz.Write(text)
	if err != nil {
		fmt.Print(err)
		return
	}
	gz.Close()

	// base64 encode it
	result = base64.RawStdEncoding.EncodeToString(buf.Bytes())

	return
}
