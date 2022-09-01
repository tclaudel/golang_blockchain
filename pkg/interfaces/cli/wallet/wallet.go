package wallet

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/wallet/generate"
)

var (
	Cmd = &cobra.Command{
		Use:   "wallet",
		Short: "Manage wallets",
		Long:  "",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				err := cmd.Help()
				cobra.CheckErr(err)
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Use + " called")
		},
	}
)

func init() {
	Cmd.AddCommand(
		generate.Cmd,
	)
}
