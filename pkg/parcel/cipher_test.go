package parcel

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCipherKey []byte = []byte("thispasswordisnot32char")

func TestEncryptWithNonce(t *testing.T) {
	plaintext := []byte("exampleplaintext")

	ciphertextWithNonce, err := EncryptWithNonce(testCipherKey, plaintext)

	assert.NoError(t, err)

	// We'll only check that the output is with an OK lengths and structure
	parts := strings.Split(ciphertextWithNonce, "|")

	assert.Len(t, parts, 3)

	assert.Len(t, parts[0], 24) // 12 * 2 == 24
	assert.Len(t, parts[1], 64) // Salt is 64
	assert.Len(t, parts[2], 64) // First block is 32 * 2 == 64
}

func TestDecryptWithNonce(t *testing.T) {
	t.Run("bad nonce encoding", func(t *testing.T) {
		ciphertextWithNonce := "nonce|ff|ff"
		_, err := DecryptWithNonce(testCipherKey, ciphertextWithNonce)
		assert.ErrorIs(t, ErrUnableDecodeNonce, err)
	})

	t.Run("bad salt encoding", func(t *testing.T) {
		ciphertextWithNonce := "ff|salt|ff"
		_, err := DecryptWithNonce(testCipherKey, ciphertextWithNonce)
		assert.ErrorIs(t, ErrUnableDecodeSalt, err)
	})

	t.Run("bad ciphertext encoding", func(t *testing.T) {
		ciphertextWithNonce := "ff|ff|ciphertext"
		_, err := DecryptWithNonce(testCipherKey, ciphertextWithNonce)
		assert.ErrorIs(t, ErrUnableDecodeCiphertext, err)
	})

	t.Run("successful decryption", func(t *testing.T) {
		ciphertextWithNonce := "b0fbebc41fce8294fcbf1d30|3f6855e5d64d189ce90ae0edaf12572d0709b5555e3c3ce7629042bba7d1fcc1|440e95a87e1defa3b51ca283846ecd7b2a48ad37c912b1473ee37a2cd59c58a1"

		plaintext, err := DecryptWithNonce(testCipherKey, ciphertextWithNonce)

		assert.NoError(t, err)
		assert.Equal(t, plaintext, []byte("exampleplaintext"))
	})
}
