package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var file string

func main() {

	flag.Parse()
	if flag.NArg() == 0 {
		file = os.Args[0]
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	} else {
		file = os.Args[1]
	}

	checkFile()
	getSize()
	fmt.Println(compress())
	zeroFile()
}

func checkFile() bool {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func getSize() int64 {
	f, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	return f.Size()

}

func compress() string {
	fileName := strings.Replace(file, ".log", "", -1)
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	outputFile, err := os.Create(fileName + "-" + formatted + ".gz")
	if err != nil {
		log.Fatal(err)
	}
	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	_, err = gzipWriter.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}

	return "file compressed"
}

func zeroFile() {
	d1 := []byte("")
	err := ioutil.WriteFile(file, d1, 0644)
	if err != nil {
		panic(err)
	}
}
