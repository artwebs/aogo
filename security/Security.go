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
	EncryptStringWithIV(key, iv, data string) (string, error)
	DecryptString(key, data string) (string, error)
	DecryptStringWithIV(key, iv, data string) (string, error)
}

type Secret struct {
	Key string
	Iv  string
}

func Encrypt(sn, val string, desObj ISecurity, secrets map[string]Secret) string {
	if val == "" {
		return ""
	}
	secret := secrets[sn]
	crypted1, err := desObj.EncryptStringWithIV(secret.Key, secret.Iv, val)
	if err != nil {
		return ""
	}
	return sn + crypted1
}

func Decrypt(val string, desObj ISecurity, secrets map[string]Secret) string {
	if val == "" {
		return ""
	}
	sn := val[0:2]
	secret := secrets[sn]
	crypted1, err := desObj.DecryptStringWithIV(secret.Key, secret.Iv, val[2:len(val)])
	if err != nil {
		return ""
	}
	return crypted1
}

type Security struct {
	keySize, ivSize int
}

func (this *Security) GenerateKey(size int) string {
	return string(utils.RandomBytes(size))
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
