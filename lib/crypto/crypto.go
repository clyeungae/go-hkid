package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

func AES256Encrypt(plainText string, key string, iv string, blockSize uint) []byte {
	if blockSize == 0 {
		blockSize = aes.BlockSize
	}
	return BaseAES256Encrypt([]byte(plainText), []byte(key), []byte(iv), blockSize)
}

func AES256EncryptWithByteKey(plainText string, byteKey []byte, iv uint, blockSize uint) []byte {
	if blockSize == 0 {
		blockSize = aes.BlockSize
	}

	bIV := make([]byte, blockSize)
	binary.LittleEndian.PutUint64(bIV, uint64(iv))

	return BaseAES256Encrypt([]byte(plainText), byteKey, bIV, blockSize)
}

func BaseAES256Encrypt(plainText []byte, byteKey []byte, byteIV []byte, blockSize uint) []byte {
	if blockSize <= 0 {
		blockSize = aes.BlockSize
	}
	plainTextBlock := PKCS5Padding(plainText, int(blockSize), len(plainText))
	block, err := aes.NewCipher(byteKey)

	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, byteIV)
	mode.CryptBlocks(cipherText, plainTextBlock)

	return cipherText
}

func PKCS5Padding(cipherText []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(cipherText)%blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func AES256Decrypt(cipherText string, key string, iv string) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	cipherTextDecoded, err := hex.DecodeString(cipherText)
	if err != nil {
		panic(err)
	}

	return AES256DecryptWithByte([]byte(cipherTextDecoded), bKey, bIV)
}

func AES256DecryptWithByte(cipherText []byte, bKey []byte, bIV []byte) string {
	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks(cipherText, cipherText)
	return string(PKCS5UnPadding(cipherText))
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func SHA256Hashing(text string) string {
	return SHA256HashingWithByte([]byte(text))
}

func SHA256HashingWithByte(byteText []byte) string {
	hash := sha256.New()
	hash.Write(byteText)
	return hex.EncodeToString(hash.Sum(nil))
}
