package create

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tclaudel/golang_blockchain/internal/values"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/clicfg"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/flags"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/rest"
	"go.uber.org/zap"
)

var (
	receiverAddress string
	amount          float64
)

var (
	Cmd = &cobra.Command{
		Use:   "new",
		Short: "new transaction",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if receiverAddress == "" {
				clicfg.Logger.Error("receiver address is required")
				cmd.Help()
				os.Exit(0)
			}

			if amount == 0 {
				clicfg.Logger.Error("amount is required")
				cmd.Help()
				os.Exit(0)
			}

			wallet, err := clicfg.Repositories.Wallet().Load(flags.WalletPath)
			if err != nil {
				clicfg.Logger.Fatal("failed to load wallet", zap.Error(err))

			}

			tx, err := values.NewTransaction(wallet, values.AddressFromString(receiverAddress), values.AmountFromFloat64(amount))
			if err != nil {
				clicfg.Logger.Fatal("failed to create transaction", zap.Error(err))
			}

			pk, err := tx.SenderPublicKey()
			if err != nil {
				clicfg.Logger.Fatal("failed to get sender public key", zap.Error(err))
			}

			sig, err := tx.Signature()
			if err != nil {
				clicfg.Logger.Fatal("failed to get signature", zap.Error(err))
			}

			client, err := rest.NewClientWithResponses(clicfg.Cfg.BlockchainNode.Address)
			if err != nil {
				clicfg.Logger.Fatal("failed to create client", zap.Error(err))
			}

			ts, err := tx.Timestamp()
			if err != nil {
				clicfg.Logger.Fatal("failed to get timestamp", zap.Error(err))
			}

			transaction := rest.CreateTransactionJSONRequestBody{
				Amount:           amount,
				RecipientAddress: receiverAddress,
				SenderAddress:    tx.SenderAddress().String(),
				SenderPublicKey:  pk.String(),
				Signature:        sig.String(),
				Timestamp:        ts.Time(),
			}

			resp, err := client.CreateTransactionWithResponse(ctx, transaction)
			if err != nil {
				clicfg.Logger.Fatal("failed to create transaction", zap.Error(err))
			}

			switch resp.StatusCode() {
			case 400:
				clicfg.Logger.Fatal(resp.JSON400.Message, zap.Int("errorCode", resp.JSON400.ErrCode))
			case 404:
				clicfg.Logger.Fatal("Not found")
			case 401:
				clicfg.Logger.Fatal(resp.JSON401.Message, zap.Int("errorCode", resp.JSON401.ErrCode))
			case 500:
				clicfg.Logger.Fatal(resp.JSON500.Message, zap.Int("errorCode", resp.JSON500.ErrCode))
			}

			data, err := json.MarshalIndent(resp.JSON200, "", "  ")
			if err != nil {
				clicfg.Logger.Fatal("Unable to serialize transaction", zap.Error(err))
			}

			fmt.Println(string(data))
		},
	}
)

func init() {
	Cmd.Flags().StringVarP(&flags.WalletPath, "wallet", "w", "./data/wallet/wallet.json", "path to wallet")
	Cmd.Flags().StringVar(&receiverAddress, "receiver", "", "receiver address")
	Cmd.Flags().Float64Var(&amount, "amount", 0, "amount to send")
}
