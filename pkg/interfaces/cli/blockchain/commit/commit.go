package commit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/clicfg"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/rest"
	"go.uber.org/zap"
)

var (
	Cmd = &cobra.Command{
		Use:   "commit",
		Short: "Commit transactions",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			client, err := rest.NewClientWithResponses(clicfg.Cfg.BlockchainNode.Address)
			if err != nil {
				clicfg.Logger.Fatal("failed to create client", zap.Error(err))
			}

			resp, err := client.CommitTransactionsWithResponse(ctx)
			if err != nil {
				clicfg.Logger.Fatal("failed to commit transactions", zap.Error(err))
			}

			switch resp.StatusCode() {
			case 400:
				clicfg.Logger.Fatal(resp.JSON400.Message, zap.Int("errorCode", resp.JSON400.ErrCode))
			case 404:
				clicfg.Logger.Fatal("Not found")
			case 401:
				clicfg.Logger.Fatal(resp.JSON400.Message, zap.Int("errorCode", resp.JSON400.ErrCode))
			case 500:
				clicfg.Logger.Fatal(resp.JSON400.Message, zap.Int("errorCode", resp.JSON400.ErrCode))
			}

			clicfg.Logger.Info("Transactions committed")

			data, err := json.MarshalIndent(resp.JSON200, "", " ")
			if err != nil {
				clicfg.Logger.Fatal("failed to serialize transactions pool", zap.Error(err))
			}

			fmt.Println(string(data))
		},
	}
)
