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
	return &SecurityAES{Security: Security{keySize: 32}}
}

func (this *SecurityAES) EncryptString(key, data string) (string, error) {
	return this.EncryptStringWithIV(key, key, data)
}

func (this *SecurityAES) EncryptStringWithIV(key, iv, data string) (string, error) {
	crypted, err := this.EncryptWithIV([]byte(key), []byte(iv), []byte(data))
	if err == nil {
		return base64.StdEncoding.EncodeToString(crypted), err
	}
	return "", err
}

func (this *SecurityAES) Encrypt(key, data []byte) ([]byte, error) {
	return this.EncryptWithIV(key, key, data)
}

func (this *SecurityAES) EncryptWithIV(key, iv, data []byte) ([]byte, error) {
	key = this.GetKey(key, this.keySize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	originalData := PKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, this.GetIV(iv, block.BlockSize()))
	crypted := make([]byte, len(originalData))
	blockMode.CryptBlocks(crypted, originalData)
	return crypted, nil
}

func (this *SecurityAES) DecryptString(key, data string) (string, error) {
	return this.DecryptStringWithIV(key, key, data)
}

func (this *SecurityAES) DecryptStringWithIV(key, iv, data string) (string, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	originalData, err := this.DecryptWithIV([]byte(key), []byte(iv), dataByte)
	if err == nil {
		return string(originalData), err
	}
	return "", err
}

func (this *SecurityAES) Decrypt(key, data []byte) ([]byte, error) {
	return this.DecryptWithIV(key, key, data)
}

func (this *SecurityAES) DecryptWithIV(key, iv, data []byte) ([]byte, error) {
	key = this.GetKey(key, this.keySize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, this.GetIV(iv, block.BlockSize()))
	originalData := make([]byte, len(data))
	blockMode.CryptBlocks(originalData, data)
	originalData = PKCS5UnPadding(originalData)
	return originalData, nil
}
