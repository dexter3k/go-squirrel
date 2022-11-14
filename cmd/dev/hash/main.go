package main

import (
	"fmt"

	sq "github.com/dexter3k/go-squirrel"
)

func main() {
	fmt.Println(sq.HashPointer(100))
}
