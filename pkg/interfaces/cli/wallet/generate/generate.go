package generate

import (
	"github.com/spf13/cobra"
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/clicfg"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/flags"
	"go.uber.org/zap"
)

var (
	Cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate wallet",
		Long:  `Generate wallet`,
		Run: func(cmd *cobra.Command, args []string) {
			wallet := values.NewWallet(values.IdentifierFromString(flags.Filename))
			clicfg.Repositories.Wallet().Save(wallet)
			clicfg.Logger.Info("Wallet generated", zap.String("address", wallet.Address().String()))
		},
	}
)

func init() {
	Cmd.Flags().StringVarP(&flags.Path, "path", "p", "./data/wallet", "path to wallet repository")
	Cmd.Flags().StringVarP(&flags.Filename, "filename", "f", "wallet.json", "filename of wallet")
}
