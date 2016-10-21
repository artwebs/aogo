package security

import (
	// "bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	// "fmt"
)

type SecurityDES struct {
	Security
}

func NewSecurityDES() *SecurityDES {
	return &SecurityDES{Security: Security{keySize: 8, ivSize: 8}}
}

func (this *SecurityDES) EncryptString(key, data string) (string, error) {
	return this.EncryptIVString(key, key, data)

}

func (this *SecurityDES) EncryptIVString(key, iv, data string) (string, error) {
	crypted, err := this.EncryptIV([]byte(key), []byte(iv), []byte(data))
	if err == nil {
		return base64.StdEncoding.EncodeToString(crypted), nil
	}
	return "", err

}

func (this *SecurityDES) Encrypt(key, data []byte) ([]byte, error) {
	return this.EncryptIV(key, key, data)
}

func (this *SecurityDES) EncryptIV(key, iv, data []byte) ([]byte, error) {
	key = this.GetKey(key, this.keySize)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data = PKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	return crypted, nil
}

func (this *SecurityDES) DecryptString(key, data string) (string, error) {
	return this.DecryptIVString(key, key, data)
}

func (this *SecurityDES) DecryptIVString(key, iv, data string) (string, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	originalData, err := this.DecryptIV([]byte(key), []byte(iv), dataByte)
	if err == nil {
		return string(originalData), err
	}
	return "", err
}

func (this *SecurityDES) Decrypt(key, data []byte) ([]byte, error) {
	return this.DecryptIV(key, key, data)
}

func (this *SecurityDES) DecryptIV(key, iv, data []byte) ([]byte, error) {
	key = this.GetKey(key, this.keySize)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	originalData := make([]byte, len(data))
	blockMode.CryptBlocks(originalData, data)
	originalData = PKCS5UnPadding(originalData)
	return originalData, nil
}

type SecurityTripleDES struct {
	Security
}

func NewSecurityTripleDES() *SecurityTripleDES {
	return &SecurityTripleDES{Security: Security{keySize: 24, ivSize: 8}}
}

func (this *SecurityTripleDES) EncryptString(key, data string) (string, error) {
	crypted, err := this.Encrypt([]byte(key), []byte(data))
	if err == nil {
		return base64.StdEncoding.EncodeToString(crypted), err
	}
	return "", err
}

func (this *SecurityTripleDES) Encrypt(key, data []byte) ([]byte, error) {
	key = this.GetKey(key, this.keySize)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	originalData := PKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, this.GetIV(key, this.ivSize))
	crypted := make([]byte, len(originalData))
	blockMode.CryptBlocks(crypted, originalData)
	return crypted, nil
}

func (this *SecurityTripleDES) DecryptString(key, data string) (string, error) {

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

func (this *SecurityTripleDES) Decrypt(key, data []byte) ([]byte, error) {
	key = this.GetKey(key, this.keySize)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, this.GetIV(key, this.ivSize))
	originalData := make([]byte, len(data))
	blockMode.CryptBlocks(originalData, data)
	originalData = ZeroUnPadding(originalData)
	return originalData, nil
}
