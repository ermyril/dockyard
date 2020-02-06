package main

import (
	"fmt"
	"github.com/ermyril/dockyard/list"
	"strconv"
)

func (yard *Dockyard) Restore() {

	fileList := list.List{}

	for i := 0; i < 50; i++ {
		fileList.Items = append(fileList.Items, strconv.Itoa(i))
	}

	fmt.Println(list.SelectItem(fileList))

}
