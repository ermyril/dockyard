// Copyright 2015 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");

package list

import (
	"fmt"
	"github.com/gdamore/tcell"
	"os"
	"time"
)

type Pos struct {
	x int
	y int
}

type List struct {
	Items []string
	cursor int
	offset int
	width  int
	height int
	pos    Pos
}


func (list *List) draw(s tcell.Screen) {
	//w, h := s.Size()
	//
	w, h := list.width, list.height

	if w == 0 || h == 0 {
		return
	}


	if h <= list.cursor {
    	list.offset++
    	list.cursor = h - 1
	} else if list.cursor < 0 {
    	list.offset--
    	list.cursor = 0
	}

	visibleItems := list.Items[list.offset:]

	st := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorDefault+1)


	char := ' '

	for row := list.pos.y; row < list.pos.y + list.height; row++ {
		for col := list.pos.x; col < list.pos.x + list.width; col++ {
			if row == list.cursor {
				if col < len(visibleItems[row]) {
			    	s.SetContent(col, row, rune(visibleItems[row][col]), nil, st)
				} else {
			    	s.SetContent(col, row, char, nil, st)
				}
			} else if row < len(visibleItems){
				if col < len(visibleItems[row]) {
					s.SetContent(col, row, rune(visibleItems[row][col]), nil, tcell.StyleDefault)
				}
			}
		}
	}
	s.Show()
}

func SelectItem(itemsList List) string {

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	s, e := tcell.NewScreen()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.Clear()


	itemsList.offset = 0
	itemsList.cursor = 0
	itemsList.width = 40
	itemsList.height = 20
	itemsList.pos = Pos{0, 0} // still unusable


	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					close(quit)
					return
				case tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
						case tcell.KeyUp :
							if itemsList.cursor + itemsList.offset > 0 {
								itemsList.cursor--
							}
					  	case tcell.KeyDown :
							if itemsList.cursor + itemsList.offset < len(itemsList.Items) - 1 {
									itemsList.cursor++
							}
						}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	cnt := 0
	dur := time.Duration(0)



loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
		start := time.Now()

		itemsList.draw(s)
			s.Clear()

		cnt++
		dur += time.Now().Sub(start)
	}

	s.Fini()

	return itemsList.Items[itemsList.cursor + itemsList.offset]
}


