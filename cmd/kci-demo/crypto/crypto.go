package crypto

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/nacl/box"
)

const (
	// KeySize represents the size of a public, secret or shared key in bytes.
	KeySize = 32
	// NonceSize represents the size of a nonce in bytes.
	NonceSize = 24
)

// Encrypt uses a crypto_box_afternm-equivalent function to encrypt the given data.
func Encrypt(data []byte, sharedKey *[KeySize]byte) ([]byte, *[NonceSize]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, nil, err
	}

	return box.SealAfterPrecomputation(nil, data, nonce, sharedKey), nonce, nil
}

// Decrypt uses a crypto_box_open_afternm-equivalent function to decrypt the given encrypted data.
func Decrypt(encryptedData []byte, sharedKey *[KeySize]byte, nonce *[NonceSize]byte) ([]byte, error) {
	data, success := box.OpenAfterPrecomputation(nil, encryptedData, nonce, sharedKey)
	if !success {
		return nil, errors.New("decryption failed")
	}

	return data, nil
}

// GenerateNonce generates a random nonce.
func GenerateNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)

	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// PrecomputeKey calculates the shared key between the peer's publicKey and our own secret key.
func PrecomputeKey(publicKey *[KeySize]byte, secretKey *[KeySize]byte) *[KeySize]byte {
	sharedKey := new([KeySize]byte)
	box.Precompute(sharedKey, publicKey, secretKey)
	return sharedKey
}

// GenerateKeyPair generates a new Curve25519 keypair.
func GenerateKeyPair() (*[KeySize]byte, *[KeySize]byte, error) {
	return box.GenerateKey(rand.Reader)
}
