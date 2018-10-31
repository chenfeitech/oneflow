package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

var key = []byte{12, 7, 21, 9, 8, 21, 12, 15, 7, 8, 1, 84, 95, 87, 84, 87}

func init() {
	for i := range key {
		key[i] = key[i] ^ 0x66
	}
}

func Decrypt(cipherstring string) (string, error) {
	ciphertext, err := hex.DecodeString(cipherstring)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCBCDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.CryptBlocks(ciphertext, ciphertext)
	return string(bytes.TrimRight(ciphertext, string([]byte{0}))), nil
}

func Encrypt(plainstring string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(plainstring)
	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = append(plaintext, bytes.Repeat([]byte{0}, aes.BlockSize-len(plaintext)%aes.BlockSize)...)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	return hex.EncodeToString(ciphertext), nil
}
