package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"io"
	"time"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/scrypt"
)

var (
	vaultBucket = []byte("Vault")
	keyPriv     = []byte("private_key")
	keySalt     = []byte("salt")
)

type SecureVault interface {
	Sign(data []byte) (signature []byte, err error)
	GetPublicKey() []byte
}

// SIMMock — garde pour les tests unitaires
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
	time.Sleep(200 * time.Millisecond)
	return ed25519.Sign(s.privateKey, data), nil
}

func (s *SIMMock) GetPublicKey() []byte {
	return s.publicKey
}

// PersistentVault — identité persistante chiffrée par PIN
type PersistentVault struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	db         *bolt.DB
}

func NewPersistentVault(db *bolt.DB, pin string) (*PersistentVault, error) {
	v := &PersistentVault{db: db}

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(vaultBucket)
		if err != nil {
			return err
		}

		encrypted := b.Get(keyPriv)
		if encrypted == nil {
			return v.generateAndStore(b, pin)
		}
		return v.loadAndDecrypt(b, pin, encrypted)
	})

	return v, err
}

func (v *PersistentVault) generateAndStore(b *bolt.Bucket, pin string) error {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	encrypted, err := encryptData([]byte(priv), pin, salt)
	if err != nil {
		return err
	}

	b.Put(keySalt, salt)
	b.Put(keyPriv, encrypted)

	v.privateKey = priv
	v.publicKey = pub
	return nil
}

func (v *PersistentVault) loadAndDecrypt(b *bolt.Bucket, pin string, encrypted []byte) error {
	salt := b.Get(keySalt)
	if salt == nil {
		return errors.New("salt introuvable")
	}

	decrypted, err := decryptData(encrypted, pin, salt)
	if err != nil {
		return errors.New("PIN incorrect ou données corrompues")
	}

	v.privateKey = ed25519.PrivateKey(decrypted)
	v.publicKey = v.privateKey.Public().(ed25519.PublicKey)
	return nil
}

func (v *PersistentVault) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(v.privateKey, data), nil
}

func (v *PersistentVault) GetPublicKey() []byte {
	return v.publicKey
}

func deriveKey(pin string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(pin), salt, 32768, 8, 1, 32)
}

func encryptData(data []byte, pin string, salt []byte) ([]byte, error) {
	key, err := deriveKey(pin, salt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decryptData(data []byte, pin string, salt []byte) ([]byte, error) {
	key, err := deriveKey(pin, salt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("données invalides")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}