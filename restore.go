package main

import (
	"fmt"
	"github.com/ermyril/dockyard/list"
	"log"
	"io/ioutil"
)

func (yard *Dockyard) Restore() {

	fileList := list.List{}
	files, err := ioutil.ReadDir(yard.config.Directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileList.Items = append(fileList.Items, file.Name())
	}

	selectedItem := list.SelectItem(fileList)

	fmt.Println(selectedItem)
}
