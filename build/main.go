package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const template = `package main
func init() {
	files[%#v] = %#v
}`

func try(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	file := os.Args[1]
	b, err := ioutil.ReadFile(file)
	try(err)

	try(ioutil.WriteFile(file+".go", []byte(fmt.Sprintf(template, file, b)), 0644))
}
