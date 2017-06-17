package cmd

import (
	"github.com/alexbakker/kci-demo/cmd/kci-demo/profile"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new profile in the current directory",
		Run:   startInit,
	}
)

func init() {
	RootCmd.AddCommand(initCmd)
}

func startInit(cmd *cobra.Command, args []string) {
	if profile.Exists() {
		logger.Fatalln("profile already exists")
	}

	p, err := profile.New()
	if err != nil {
		logger.Fatalf("unable to create new profile: %s", err)
	}

	if err := profile.Save(p); err != nil {
		logger.Fatalf("unable to save new profile: %s", err)
	}
}
