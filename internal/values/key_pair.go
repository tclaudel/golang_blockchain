package values

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
)

type KeyPair struct {
	PublicKey
	PrivateKey
}

func GenerateKeyPair() KeyPair {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return KeyPair{
		PublicKey:  NewPublicKey(privateKey.PublicKey),
		PrivateKey: NewPrivateKey(*privateKey),
	}
}

func (kp *KeyPair) UnmarshalJSON(data []byte) error {
	var dto = new(struct {
		X             string `json:"x"`
		Y             string `json:"y"`
		D             string `json:"d"`
		EllipticCurve string `json:"elliptic_curve"`
	})
	err := json.Unmarshal(data, dto)
	if err != nil {
		return err
	}

	var (
		x = new(big.Int)
		y = new(big.Int)
		d = new(big.Int)
	)

	if _, ok := x.SetString(dto.X, 10); !ok {
		return fmt.Errorf("invalid x")
	}
	if _, ok := y.SetString(dto.Y, 10); !ok {
		return fmt.Errorf("invalid y")
	}
	if _, ok := d.SetString(dto.D, 10); !ok {
		return fmt.Errorf("invalid d")
	}

	var e ecdsa.PrivateKey

	e.D = d
	e.PublicKey.Curve = elliptic.P256()
	e.PublicKey.X = x
	e.PublicKey.Y = y
	kp.PublicKey = NewPublicKey(e.PublicKey)
	kp.PrivateKey = NewPrivateKey(e)

	return nil
}

func (kp KeyPair) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		X             string `json:"x"`
		Y             string `json:"y"`
		D             string `json:"d"`
		EllipticCurve string `json:"elliptic_curve"`
	}{
		X:             kp.publicKey.X.String(),
		Y:             kp.publicKey.Y.String(),
		D:             kp.privateKey.D.String(),
		EllipticCurve: "P256",
	})
}
