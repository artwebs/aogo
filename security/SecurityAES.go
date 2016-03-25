package security

import (
	// "bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type SecurityAES struct {
	Security
}

func NewSecurityAES() *SecurityAES {
	return &SecurityAES{Security: Security{keySize: 24}}
}

func (this *SecurityAES) EncryptString(key, data string) (string, error) {
	crypted, err := this.Encrypt([]byte(key), []byte(data))
	if err == nil {
		return base64.StdEncoding.EncodeToString(crypted), err
	}
	return "", err
}

func (this *SecurityAES) Encrypt(key, data []byte) ([]byte, error) {
	key = this.GetKey(key, this.keySize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	originalData := PKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, this.GetIV(key, block.BlockSize()))
	crypted := make([]byte, len(originalData))
	blockMode.CryptBlocks(crypted, originalData)
	return crypted, nil
}

func (this *SecurityAES) DecryptString(key, data string) (string, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	originalData, err := this.Decrypt([]byte(key), dataByte)
	if err == nil {
		return string(originalData), err
	}
	return "", err
}

func (this *SecurityAES) Decrypt(key, data []byte) ([]byte, error) {
	key = this.GetKey(key, this.keySize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, this.GetIV(key, block.BlockSize()))
	originalData := make([]byte, len(data))
	blockMode.CryptBlocks(originalData, data)
	originalData = PKCS5UnPadding(originalData)
	return originalData, nil
}
