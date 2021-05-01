package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

// DecryptPassword with AES GCM mode
func DecryptPassword(encPassword, masterKey []byte) ([]byte, error) {
	// trim first 3 bytes, signature "v10"
	nonce := encPassword[3:15]  // random value eg. IV
	payload := encPassword[15:] // encPassword

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, err
	}

	blockMode, _ := cipher.NewGCM(block)

	return blockMode.Open(nil, nonce, payload, nil)
}
