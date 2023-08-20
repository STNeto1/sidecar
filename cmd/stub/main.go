package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	envs := os.Environ()

	for _, value := range envs {

		if strings.HasPrefix(value, "SIDECAR") {
			fmt.Println(value)
		}
	}
}
