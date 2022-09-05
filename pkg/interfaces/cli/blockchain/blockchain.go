package blockchain

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
		Use:   "blockchain",
		Short: "Manage blockchain",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			client, err := rest.NewClientWithResponses(clicfg.Cfg.BlockchainNode.Address)
			if err != nil {
				clicfg.Logger.Fatal("Unable to create http client", zap.Error(err))
			}

			bcResp, err := client.GetBlockchainWithResponse(ctx)
			if err != nil {
				clicfg.Logger.Fatal("Unable to get blockchain", zap.Error(err))
			}
			switch bcResp.StatusCode() {
			case 400:
				clicfg.Logger.Fatal(bcResp.JSON400.Message, zap.Int("errorCode", bcResp.JSON400.ErrCode))
			case 404:
				clicfg.Logger.Fatal("Not found")
			case 401:
				clicfg.Logger.Fatal(bcResp.JSON400.Message, zap.Int("errorCode", bcResp.JSON400.ErrCode))
			case 500:
				clicfg.Logger.Fatal(bcResp.JSON400.Message, zap.Int("errorCode", bcResp.JSON400.ErrCode))
			}

			data, err := json.MarshalIndent(bcResp.JSON200, "", " ")
			if err != nil {
				clicfg.Logger.Fatal("Unable to serialize blockchain", zap.Error(err))
			}

			fmt.Println(string(data))
		},
	}
)

func init() {
}
