package values

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestBytesTransition(t *testing.T) {
	kp := GenerateKeyPair()
	b, err := json.Marshal(kp)
	assert.NoError(t, err)

	var kp2 KeyPair
	err = json.Unmarshal(b, &kp2)
	assert.NoError(t, err)

	t.Log(kp.String())
	t.Log(kp2.String())
	if diff := deep.Equal(kp, kp2); diff != nil {
		t.Error(diff)
	}
}
