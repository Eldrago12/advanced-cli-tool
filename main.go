package main

import (
	"fmt"
	"os"

	"github.com/Eldrago12/advanced-cli-tool/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
