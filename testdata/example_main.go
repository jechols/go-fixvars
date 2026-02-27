package main

import (
	"fmt"
	"os"
)

func main() {
	x := 1
	var y = x + 1

	if info, err := os.Stat("no"); err == nil {
		fmt.Printf("Info: %#v\n", info)
	} else {
		fmt.Printf("y is %d\n", y)
	}
}
