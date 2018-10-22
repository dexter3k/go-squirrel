package main

import (
	"fmt"

	sq "github.com/SMemsky/go-squirrel"
)

func main() {
	foo := "FOO"
	fmt.Println(sq.HashPointer(&foo))
}
