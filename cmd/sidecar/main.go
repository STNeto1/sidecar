package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sidecar/pkg"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "sidecar",
		Usage: "Easily inject environment variables",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list profiles",
				Action: func(cCtx *cli.Context) error {
					pkg.ListProfiles()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func injectMain() {
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
