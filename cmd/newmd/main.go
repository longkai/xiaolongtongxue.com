package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	templ = `%s
===
Content goes here...

### EOF
` + "```yaml" + `
date: %s
summary: 
weather:
license: all-rights-reserved
location: 
background:
tags:
  - tag1
  - tag2
` + "```"
)

var mkdir = func(name string) error {
	return os.MkdirAll(name, 0755)
}

func mayFail(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s title\n", os.Args[0])
		os.Exit(1)
	}
	title := os.Args[1]
	dir := formatName(title)

	mayFail(mkdir(dir))

	f := filepath.Join(dir, "README.md")
	_, err := os.Stat(f)
	if err == nil {
		fmt.Printf("%s already existed, discard it? y/n: ", f)
		b := make([]byte, 50)
		n, err := os.Stdin.Read(b)
		mayFail(err)
		if n < 1 || b[0] != 'y' {
			fmt.Println("aborted.")
			os.Exit(0)
		}
	}

	out, err := os.Create(f)

	mayFail(err)

	mayFail(newMD(title, out))

	fmt.Println("done :)")
}

func newMD(title string, out io.Writer) error {
	_, err := out.Write([]byte(fmt.Sprintf(templ, title, time.Now().Format(time.RFC3339))))
	return err
}

func formatName(name string) string {
	return strings.Replace(strings.ToLower(name), " ", "-", -1)
}
