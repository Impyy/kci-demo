package cmd

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/Impyy/kci-demo/cmd/kci-demo/client"
	"github.com/Impyy/kci-demo/cmd/kci-demo/crypto"
	"github.com/Impyy/kci-demo/cmd/kci-demo/profile"
	"github.com/spf13/cobra"
)

type connectFlags struct {
	Host      string
	Port      int
	KCI       bool
	PublicKey string
}

var (
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to an instance of kci-demo",
		Run:   startConnect,
	}
	connectCmdFlags = new(connectFlags)
)

func init() {
	RootCmd.AddCommand(connectCmd)
	connectCmd.Flags().BoolVarP(&connectCmdFlags.KCI, "kci", "", false, "Whether to enable KCI mode or not")
	connectCmd.Flags().StringVarP(&connectCmdFlags.Host, "host", "", "", "The host to connect to")
	connectCmd.Flags().IntVarP(&connectCmdFlags.Port, "port", "", 0, "The port to connect to")
	connectCmd.Flags().StringVarP(&connectCmdFlags.PublicKey, "key", "", "", "The key of the host we want to connect to")
}

func startConnect(cmd *cobra.Command, args []string) {
	tempPublicKey, err := hex.DecodeString(connectCmdFlags.PublicKey)
	if err != nil {
		logger.Fatalf("bad public key: %s", err)
	}

	publicKey := new([crypto.KeySize]byte)
	copy(publicKey[:], tempPublicKey)

	prof, err := profile.Load()
	if err != nil {
		logger.Fatalf("unable to load profile: %s", err)
	}

	client, err := client.New(prof, connectCmdFlags.KCI)
	if err != nil {
		logger.Fatalf("unable to init client: %s", err)
	}
	client.Print()

	go func() {
		if err := client.Listen(); err != nil {
			logger.Fatalf("error: %s", err)
		}
	}()

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", connectCmdFlags.Host, connectCmdFlags.Port))
	if err != nil {
		logger.Fatalf("error: %s", err)
	}

	err = client.Handshake(addr, publicKey)
	if err != nil {
		logger.Fatalf("unable to send handshake packet: %s", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
