package main

import (
	"fmt"
	"io/ioutil"
)

// save save content to file
func save(c, f string) {
	fmt.Println("Saving to", f)
	ioutil.WriteFile(f, []byte(c), 0644)
}

type LocalStorage struct {
	Folder string
}

func (l LocalStorage) Save(content, name string) {
	save(content, l.Folder+"/"+name+".sql")
}
