package temporalclient

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type AESEncryptionServiceV1 struct {
	Cipher cipher.AEAD
}

func newAESEncryptionServiceV1(opts Options) (*AESEncryptionServiceV1, error) {
	// must be 16, 24, 32 byte length
	// this is your encryption key
	// will fail to initialize if length requirements are not met
	cipherBlock, err := aes.NewCipher(opts.EncryptionKey)
	if err != nil {
		// likely invalid key length if errors here
		return nil, err
	}
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}
	return &AESEncryptionServiceV1{
		Cipher: gcm,
	}, nil
}

// Encrypt takes a byte array and returns an encrypted byte array
// as base64 encoded
func (a AESEncryptionServiceV1) Encrypt(unencryptedBytes []byte) ([]byte, error) {
	if len(unencryptedBytes) == 0 { // prevent err on empty byte arrays - "cipher: message authentication failed"
		return []byte(""), nil
	}
	nonce := make([]byte, a.Cipher.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	encryptedBytes := a.Cipher.Seal(nonce, nonce, unencryptedBytes, nil)
	encryptedEncodedData := make([]byte, base64.RawURLEncoding.EncodedLen(len(encryptedBytes)))
	base64.RawURLEncoding.Encode(encryptedEncodedData, encryptedBytes)
	return encryptedEncodedData, nil
}

// Decrypt takes an encrypted base64 byte array then
// returns an unencrypted byte array if same key was used to encrypt it
func (a AESEncryptionServiceV1) Decrypt(encryptedBytes []byte) ([]byte, error) {
	if len(encryptedBytes) == 0 {
		return []byte(""), nil
	}
	decodedEncryptedBytes := make([]byte, base64.RawURLEncoding.DecodedLen(len(encryptedBytes)))
	if _, err := base64.RawURLEncoding.Decode(decodedEncryptedBytes, encryptedBytes); err != nil {
		return nil, err
	}
	nonceSize := a.Cipher.NonceSize()
	if len(encryptedBytes) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: %v", len(encryptedBytes))
	}
	return a.Cipher.Open(nil, decodedEncryptedBytes[:nonceSize], decodedEncryptedBytes[nonceSize:], nil)
}
