/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/spf13/cobra"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/cli/clicfg"
)

func main() {
	cobra.OnInitialize(clicfg.InitConfig)
	cli.Execute()
}
