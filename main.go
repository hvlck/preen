package preen

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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
	return regexp.MustCompile("https?://([.]*.?)*/?").FindAllString(contents, -1)
}

func main() {
	files := Read(".")

	for _, file := range files {
		if strings.HasPrefix(file, ".") {
			continue
		}

		fmt.Printf("checking file %s\n", file)
		f, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("failed to read file %s: %s", file, err)
		}

		matches := Find(string(f))

		for _, match := range matches {
			r, err := http.Get(match)

			if err != nil {
				log.Printf("failed to get %s: %s\n", match, err)
			}

			if r.StatusCode != 200 {
				log.Printf("%s responded with something other than status 200\n", match)
			} else if r.Request.URL.String() != match {
				log.Printf("%s was redirected", match)
			} else {
				log.Printf("successfully retrieved %s\n", match)
			}
		}
	}
}
