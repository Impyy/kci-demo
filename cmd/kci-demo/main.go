package main

import (
	"fmt"

	"github.com/alexbakker/kci-demo/cmd/kci-demo/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
}
