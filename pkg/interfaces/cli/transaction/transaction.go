package transaction

import (
	"github.com/spf13/cobra"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/transaction/create"
)

var (
	Cmd = &cobra.Command{
		Use:   "transaction",
		Short: "Manage transactions",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
		},
	}
)

func init() {
	Cmd.AddCommand(
		create.Cmd,
	)
}
