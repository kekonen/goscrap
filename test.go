package main

import (
	"fmt"
	"string"
)


func main() {
	st := "lol@kek"
	i:=string.Index(st, "@")
	fmt.Printf("lol: %v", i)
}

