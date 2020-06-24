package main

import (
	"fmt"

	"./args"
)

func main() {
	arguments := []string{"-l", "-p", "ads", "-d", "test"}
	var a args.Args
	a.Init("l,p#,d*", arguments)
	fmt.Println(a.GetInt("p"))
	fmt.Println(a.ErrorMessage())
}
