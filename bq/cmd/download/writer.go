package main

import (
	"fmt"
	"io/ioutil"
)

// save saves a string to a file
func save(c, f string) {
	fmt.Println("Saving to", f)
	ioutil.WriteFile(f, []byte(c), 0644)
}

type LocalStorage struct {
	Folder string
}

func (l LocalStorage) SaveView(content, name string) {
	save(content, l.Folder+"/"+name+".sql")
}

// This is not a well thought design
func (l LocalStorage) SaveSchema(content []byte, name string) {
	path := l.Folder + "/" + name + ".schema.json"
	fmt.Println("Saving to", path)
	ioutil.WriteFile(path, content, 0644)
}
