package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/cheynewallace/tabby"
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

func createFail(url string, status int) *BadLink {
	return &BadLink{
		Url:    url,
		Status: status,
	}
}

type BadLink struct {
	Url      string
	Status   int
	Response string
}

func main() {
	links := make(map[string][]*BadLink)
	files := Read(".")

	var fails uint32

	for _, file := range files {
		if stat, _ := os.Stat(file); stat.IsDir() {
			continue
		}

		if strings.HasPrefix(file, ".") {
			continue
		}

		f, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("failed to read file %s: %s", file, err)
		}

		if !utf8.Valid(f) {
			fmt.Printf("%s is not valid utf-8, skipping...\n", file)
			continue
		}

		matches := Find(string(f))

		for _, match := range matches {
			if strings.HasSuffix(file, ".md") {
				match = strings.TrimSuffix(match, ")") // removes excess parantheses from markdown links
			}
			r, err := http.Get(match)

			if err != nil {
				log.Printf("failed to get %s: %s\n", match, err)
			}

			failure := createFail(match, r.StatusCode)
			if r.StatusCode != 200 {
				failure.Response = r.Status
				links[file] = append(links[file], failure)
				fails++
			} else if r.Request.URL.String() != match {
				failure.Response = "redirected"
				links[file] = append(links[file], failure)
				fails++
			}
		}
	}

	if fails == 0 {
		fmt.Printf("checked all files, no dead links found")
	} else {
		for name, file := range links {
			table := tabby.New()
			table.AddHeader("file", "url", "status", "code")
			for _, link := range file {
				table.AddLine(name, link.Url, link.Response, link.Status)
			}

			table.Print()
			fmt.Println()
		}

		fmt.Printf("finished checking all files, %v failed\n", fails)
	}
}
