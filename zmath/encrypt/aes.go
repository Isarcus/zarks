package encrypt

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/Isarcus/zarks/zmath"
)

// AESProtocol specifies the length of an AES key (bits in the name, bytes in the value)
type AESProtocol int

// AES Protocol
const (
	AES128 AESProtocol = 16
	AES192 AESProtocol = 24
	AES256 AESProtocol = 32
)

// AES is an AES cipher!
type AES cipher.Block

// NewAESCipher returns a new AES Cipher from the specified key and protocol
func NewAESCipher(key string, p AESProtocol) AES {
	keyBytes := make([]byte, p)
	for i, letter := range key {
		if i == int(p) {
			break
		}
		keyBytes[i] = byte(letter)
	}

	cipher, _ := aes.NewCipher(keyBytes)
	return cipher
}

// Encrypt will encrypt a message of any length.
func Encrypt(a AES, message []byte) []byte {
	var (
		blockSize      = a.BlockSize()
		encodedMessage = make([]byte, 0, len(message))
		encodedBytes   = 0
		blockCt        = 0
	)
	for encodedBytes < len(message) {
		var (
			src = make([]byte, blockSize)
			dst = make([]byte, blockSize)

			start = blockCt * blockSize
			end   = start + zmath.MinInt(len(message)-start, blockSize)
		)
		copy(src, message[start:end])
		a.Encrypt(dst, src)

		encodedMessage = append(encodedMessage, dst...)

		blockCt++
		encodedBytes += blockSize
	}

	return encodedMessage
}

// Decrypt will decrypt a message of any length, but will (likely incorrectly) assume that any missing parts
// of a data block are null
func Decrypt(a AES, message []byte) []byte {
	var (
		blockSize      = a.BlockSize()
		encodedMessage = make([]byte, 0, len(message))
		encodedBytes   = 0
		blockCt        = 0
	)
	for encodedBytes < len(message) {
		var (
			src = make([]byte, blockSize)
			dst = make([]byte, blockSize)

			start = blockCt * blockSize
			end   = start + zmath.MinInt(len(message)-start, blockSize)
		)
		copy(src, message[start:end])
		a.Decrypt(dst, src)

		encodedMessage = append(encodedMessage, dst...)

		blockCt++
		encodedBytes += blockSize
	}

	return encodedMessage
}
