package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sidecar/pkg"
)

func main() {
	if err := os.Setenv(pkg.LOOKUP_KEY, "SIDECAR"); err != nil {
		fmt.Println("error setting env", err)
		return
	}

	cmd := exec.Command("./stub")

	// var outputBuf bytes.Buffer
	// var errBuf bytes.Buffer
	// cmd.Stdout = &outputBuf
	// cmd.Stderr = &errBuf

	// err := cmd.Run()
	// if err != nil {
	// 	log.Println("Error executing binary:", err)
	// 	return
	// }
	//
	// fmt.Println("out -> \n", outputBuf.String())
	// fmt.Println("err -> ", errBuf.String())

	pipe, _ := cmd.StdoutPipe()

	cmd.Start()

	// var tokens []string

	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()

		log.Println(m)

		// if m == "\n" {
		// 	merged := strings.Join(tokens, " ")
		// 	fmt.Println(merged)
		// 	tokens = nil
		// 	continue
		// }
		//
		// tokens = append(tokens, m)
	}

	cmd.Wait()

}
