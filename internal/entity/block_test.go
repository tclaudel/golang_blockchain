package entity

import (
	"testing"

	"github.com/go-playground/assert"
	"github.com/tclaudel/golang_blockchain/internal/values"
)

func TestGenesis(t *testing.T) {
	genesisBlock := Genesis

	t.Log("genesis block:", genesisBlock)

	assert.IsEqual(true, genesisBlock.nonce.Equal(values.GenesisNonce))
	assert.IsEqual(true, genesisBlock.previousHash.Equal(values.GenesisHash))
	assert.IsEqual(0, len(genesisBlock.transactions))
}
