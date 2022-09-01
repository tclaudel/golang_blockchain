package blockchain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/clicfg"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http"
	"go.uber.org/zap"
)

var (
	Cmd = &cobra.Command{
		Use:   "blockchain",
		Short: "Manage blockchain",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			client, err := http.NewClientWithResponses(clicfg.Cfg.BlockchainNode.Address)
			if err != nil {
				clicfg.Logger.Fatal("Unable to create http client", zap.Error(err))
			}

			bcResp, err := client.GetBlockchainWithResponse(ctx)
			if err != nil {
				clicfg.Logger.Fatal("Unable to get blockchain", zap.Error(err))
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
