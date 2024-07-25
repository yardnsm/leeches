package parcel

// Mostly taken from https://gist.github.com/kkirsche/e28da6754c39d5e7ea10

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/scrypt"
)

const DELIMETER = "|"

func EncryptWithNonce(key []byte, plaintext []byte) (string, error) {
	key, salt, err := deriveKey(key, nil)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	return fmt.Sprintf("%x%s%x%s%x", nonce, DELIMETER, salt, DELIMETER, ciphertext), nil
}

func DecryptWithNonce(key []byte, ciphertextWithNonce string) ([]byte, error) {
	parts := strings.Split(ciphertextWithNonce, DELIMETER)

	if len(parts) != 3 {
		return nil, ErrInvalidParcel
	}

	nonce, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, ErrUnableDecodeNonce
	}

	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, ErrUnableDecodeSalt
	}

	ciphertext, err := hex.DecodeString(strings.TrimSpace(parts[2]))
	if err != nil {
		return nil, ErrUnableDecodeCiphertext
	}

	key, _, err = deriveKey(key, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return plaintext, nil
}

// Please not that we are using scrypt, which is DAMN SLOW. This is fine for us, as we want it to be
// FUCKING DAMN SLOW.
func deriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
