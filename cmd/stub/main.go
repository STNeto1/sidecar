package main

import (
	"fmt"
	"os"
	"sidecar/pkg"
	"time"
)

func main() {
	lookup := os.Getenv(pkg.LOOKUP_KEY)

	if lookup == "" {
		fmt.Println("no env was set")
		return
	}

	fmt.Println(lookup)

	init := 1

	for {
		if init >= 10 {
			break
		}

		fmt.Println(init)
		time.Sleep(time.Second)

		init += 1
	}
}
