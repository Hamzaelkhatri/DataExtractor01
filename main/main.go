package main

import (
	"encoding/base64"
	"io"
	"log"
	"os"
	"strings"

	dataextractor01 "github.com/Hamzaelkhatri/DataExtractor01"
)

func main() {
	base, err := dataextractor01.ExtractData(792)
	if err != nil {
		log.Println(err)
		return
	}

	// save base to png file
	file, err := os.Create("base.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	base64ToReader := func(base string) io.Reader {
		decode := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base))
		reader := io.Reader(decode)
		return reader
	}

	_, err = io.Copy(file, base64ToReader(base))
	if err != nil {
		panic(err)
	}
}
