package main

import (
	"fmt"

	sq "github.com/SMemsky/go-squirrel"
)

func main() {
	fmt.Println(sq.HashPointer(100))
}
