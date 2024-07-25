package parcel

import (
	"encoding/json"
)

// The key should be 32-bytes long
type ParcelKey []byte

func Marshal(v any, key ParcelKey) ([]byte, error) {
	encoded, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	encrypted, err := EncryptWithNonce(key, encoded)
	if err != nil {
		return nil, err
	}

	return []byte(encrypted), nil
}

func Unmarshal(data []byte, key ParcelKey, v any) error {
	decrypted, err := DecryptWithNonce(key, string(data))
	if err != nil {
		return err
	}

	return json.Unmarshal(decrypted, v)
}
