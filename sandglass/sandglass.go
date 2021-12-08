package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Indexes int

const (
	size = Indexes(iota)
	char
	color
	back
	colorDef = string("\033[0m")
	backDef  = string("\033[40m")
)

type Options []string

type Mod func(*Options)

func setsize(newsize int) Mod {
	return func(glass *Options) {
		(*glass)[size] = fmt.Sprint(newsize)
	}
}

func setchar(newchar rune) Mod {
	return func(glass *Options) {
		(*glass)[char] = string(newchar)
	}
}

func setcolor(newcolor string) Mod {
	return func(glass *Options) {
		(*glass)[color] = newcolor
	}
}

func setback(newback string) Mod {
	return func(glass *Options) {
		(*glass)[back] = newback
	}
}

func sandglass(mods ...Mod) {
	glass := Options{"7", "X", colorDef, backDef}

	for _, mod := range mods {
		mod(&glass)
	}

	var width int
	width, err := strconv.Atoi(glass[size])
	if err != nil || width < 0 {
		fmt.Println("sandlass: size option error")
		return
	}

	// config blank space according to character width
	blank := " "
	if len(glass[char]) > 1 {
		blank = strings.Repeat(" ", len(glass[char])-1)
	}

	fmt.Println(glass[color])

	var j int
	for i := -width + 1; i < width; i++ {
		if i < 0 {
			j = -i + 1
		} else {
			j = i + 1
		}

		switch {
		// print borders
		case j == width:
			fmt.Println(strings.Repeat(glass[char], j))
		// print centre
		case i == -width/2:
			fmt.Println(strings.Repeat(blank, width-j) + glass[char])
			i = -i
		// print middle
		default:
			fmt.Println(strings.Repeat(blank, width-j) + glass[char] + glass[back] +
				strings.Repeat(blank, 2*j-width-2) + backDef + glass[char])
		}
	}

	fmt.Println(colorDef)
}

func main() {
	sandglass()
	sandglass(setcolor("\033[1;32m"), setchar('A'), setsize(9))
	sandglass(setsize(15), setcolor("\033[1;31m"))
	sandglass(setcolor("\033[1;33m"), setchar('u'))
	sandglass(setchar('£'), setsize(13))
	sandglass(setsize(11), setchar('你'), setcolor("\033[1;34m"))
	sandglass(setsize(-42))
	sandglass(setchar('$'), setsize(0))

	// additional parameter "back": filling the inside of the sandglass
	sandglass(setback("\033[43m"), setcolor("\033[0;32m"))
	sandglass(setchar('歲'), setback("\033[44m"), setsize(21), setcolor("\033[0;36m"))
}
