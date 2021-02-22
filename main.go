package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Read utility function to get all files in a directory
func Read(dir string) []string {
	file, err := os.Open(dir)
	if err != nil {
		log.Fatalf("failed to read directory: %s\n", err)
	}

	defer file.Close()

	files, err := file.Readdirnames(0)
	if err != nil {
		log.Fatalf("failed to read contents of directory: %s\n", err)
	}

	return files
}

// Find finds all urls within a string
func Find(contents string) []string {
	return regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`).FindAllString(contents, -1)
}

func main() {
	files := Read(".")

	var fails uint32

	for _, file := range files {
		if strings.HasPrefix(file, ".") {
			continue
		}

		fmt.Printf("checking file %s\n", file)
		f, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("failed to read file %s: %s", file, err)
		}

		if utf8.Valid(f) == false {
			fmt.Printf("%s is not valid utf-8, skipping...\n", file)
			continue
		}

		matches := Find(string(f))

		for _, match := range matches {
			if strings.HasSuffix(file, ".md") == true {
				match = strings.TrimSuffix(match, ")") // removes excess parantheses from markdown links
			}
			r, err := http.Get(match)

			if err != nil {
				log.Printf("failed to get %s: %s\n", match, err)
			}

			if r.StatusCode != 200 {
				log.Printf("%s responded with something other than status 200\n", match)
				fails++
			} else if r.Request.URL.String() != match {
				log.Printf("%s was redirected", match)
				fails++
			} else {
				fmt.Printf("successfully retrieved %s\n", match)
			}
		}

		fmt.Printf("finished checking file %s\n", file)
	}

	fmt.Printf("finished checking all files, %v failed\n", fails)
}
