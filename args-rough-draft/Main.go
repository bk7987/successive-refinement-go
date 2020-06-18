package main

import (
	"./args"
)

func main() {
	arguments := []string{"-l", "-p", "32", "-d", "test"}
	var a args.Args
	a.Init("l,p#,d*", arguments)
}
