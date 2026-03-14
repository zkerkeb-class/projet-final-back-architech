package vault

import (
	"crypto/ed25519"
	"crypto/rand"
	"time"
)

// Le reste du code ne connaît que ça.
type SecureVault interface {
	Sign(data []byte) (signature []byte, err error)
	GetPublicKey() []byte
}

// L'implémentation mock qui simule une carte SIM.
type SIMMock struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewSIMMock() (*SIMMock, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &SIMMock{privateKey: priv, publicKey: pub}, nil
}

func (s *SIMMock) Sign(data []byte) ([]byte, error) {
	time.Sleep(200 * time.Millisecond) // latence matérielle simulée
	return ed25519.Sign(s.privateKey, data), nil
}

func (s *SIMMock) GetPublicKey() []byte {
	return s.publicKey
}