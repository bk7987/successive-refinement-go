package main

import (
	"./args"
)

func main() {
	arguments := []string{"t", "3"}
	var a args.Args
	a.Init("l", arguments)
}
