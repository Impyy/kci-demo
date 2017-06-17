package cmd

import (
	"log"
	"os"

	"github.com/alexbakker/kci-demo/cmd/kci-demo/client"
	"github.com/alexbakker/kci-demo/cmd/kci-demo/profile"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "kci-demo",
		Short: "kci-demo demonstrates the KCI attack",
		Run:   startRoot,
	}
	logger = log.New(os.Stderr, "", 0)
)

func startRoot(cmd *cobra.Command, args []string) {
	prof, err := profile.Load()
	if err != nil {
		logger.Fatalf("unable to load profile: %s", err)
	}

	client, err := client.New(prof, false)
	if err != nil {
		logger.Fatalf("unable to init client: %s", err)
	}
	client.Print()

	if err := client.Listen(); err != nil {
		logger.Fatalf("error: %s", err)
	}
}
