package security

import (
	"bytes"
	// "crypto/aes"
	// "crypto/cipher"
	// "crypto/des"
	// "encoding/base64"
	// "fmt"
	"github.com/artwebs/aogo/utils"
)

type SecurityMOD int

const (
	ECB SecurityMOD = 1 + iota
	CBC
)

type ISecurity interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	EncryptString(key, data string) (string, error)
	DecryptString(key, data string) (string, error)
}

type Security struct {
	keySize, ivSize int
}

func (this *Security) GenerateKey(size int) string {
	return string(util.RandomBytes(size))
}

func (this *Security) GetKey(key []byte, keySize int) []byte {
	rs := make([]byte, keySize)
	num := 0
	for i := 0; i < keySize; i++ {
		if num >= len(key) {
			num = 0
		}
		rs[i] = key[num]
		num++
	}
	return rs
}

func (this *Security) GetIV(iv []byte, blockSize int) []byte {
	rs := make([]byte, blockSize)
	num := 0
	for i := 0; i < blockSize; i++ {
		if num >= len(iv) {
			num = 0
		}
		rs[i] = iv[num]
		num++
	}
	return rs
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
