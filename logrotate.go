package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var maxAge = flag.Int64("maxage", 0, "maxAge")
	var file = flag.String("file", "", "The file we want to compress. i.e. ../files/pocket_error.lg")
	flag.Parse()

	// The flags are required.
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			log.Fatal("-"+f.Name, " not set!")
		}
	})

	// Check the file and compress.
	checkFile(*file)
	getSize(*file)
	fmt.Println(compress(*file))
	zeroFile(*file)

	// Delete old files.
	if *maxAge > 0 {
		filePattern := strings.Replace(*file, ".log", "*", 1)
		filename := strings.Replace(*file, ".log", "", 1)

		// Search for files with the pattern.
		files, err := filepath.Glob(filePattern)
		if err != nil {
			fmt.Println(err)
		}

		// No files found.
		if len(files) == 0 {
			os.Exit(0)
		}

		for _, file := range files {
			if strings.Contains(file, filename+"-") {
				fileToRemove := strings.Replace(file, ".gz", "", 1)
				date := strings.Replace(fileToRemove, filename+"-", "", 1)
				today := time.Now()

				fileDate, err := time.Parse("2006-01-02T15:04:05", date)
				if err != nil {
					panic(err)
				}
				difference := today.Sub(fileDate)
				days := int64(difference.Hours() / 24)
				if days >= *maxAge {
					e := os.Remove(file)
					if e != nil {
						log.Fatal(e)
					}
					fmt.Println("Deleted " + file)
				}
			}
		}
	}
}

func checkFile(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func getSize(file string) int64 {
	f, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	return f.Size()

}

func compress(file string) string {
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

func zeroFile(file string) {
	d1 := []byte("")
	err := ioutil.WriteFile(file, d1, 0644)
	if err != nil {
		panic(err)
	}
}
