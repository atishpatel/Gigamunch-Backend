package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
)

var (
	banWords = []string{"anal", "anus", "ballsack", "blowjob", "dick", "dildo", "nigger", "penis", "vagina"}
)

// ContainsBanWord checks if a sentence contains a ban word
func ContainsBanWord(sentence string) bool {
	for index := range banWords {
		if strings.Contains(sentence, banWords[index]) {
			return false
		}
	}
	return true
}

// ConcatenateStrings fast way of concatenating strings
func ConcatenateStrings(args ...string) string {
	var buffer bytes.Buffer
	for _, value := range args {
		buffer.WriteString(value)
	}
	return buffer.String()
}

// Decrypt decrypts the ciphertext with the keystring
func Decrypt(keyString, ciphertextString string) (string, error) {
	key := []byte(keyString)
	//ciphertext, _ := hex.DecodeString(ciphertextString)
	ciphertext, _ := base64.StdEncoding.DecodeString(ciphertextString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// Encrypt encrypts the passed in text with the key and
// returns a base64 encoded ciphertext
func Encrypt(keyString, textString string) (string, error) {
	key := []byte(keyString)
	text := []byte(textString)

	if len(text)%aes.BlockSize != 0 {
		return "", fmt.Errorf("text is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	var ciphertext []byte
	for len(ciphertext) < aes.BlockSize+len(text) {
		ciphertext = append(ciphertext, []byte("f363f3ccdcb12bb883abf484ba77d9cd7d32b5baecb3d4b1b3e0e4beffdb3ded")...)
	}
	ciphertext = ciphertext[:aes.BlockSize+len(text)]
	iv := ciphertext[:aes.BlockSize]

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], text)
	//return fmt.Sprintf("%X", ciphertext), nil
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
