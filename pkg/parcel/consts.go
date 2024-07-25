package parcel

import "errors"

var (
	ErrInvalidParcel = errors.New("ciphertextWithNonce did not contain required amount of parts")

	ErrUnableDecodeNonce      = errors.New("unable to decode nonce from parcel")
	ErrUnableDecodeSalt       = errors.New("unable to decode salt from parcel")
	ErrUnableDecodeCiphertext = errors.New("unable to decode ciphertext from parcel")
)
