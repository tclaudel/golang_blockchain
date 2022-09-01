package values

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestWallet_MarshalJSON(t *testing.T) {
	wallet := NewWallet()

	t.Log(wallet.KeyPair)
	b, err := json.Marshal(wallet)
	assert.NoError(t, err)
	t.Log(string(b))

	var wallet2 Wallet
	err = json.Unmarshal(b, &wallet2)
	assert.NoError(t, err)
	t.Log(wallet2.KeyPair)
	if diff := deep.Equal(wallet, wallet2); diff != nil {
		t.Error(diff)
	}
}
