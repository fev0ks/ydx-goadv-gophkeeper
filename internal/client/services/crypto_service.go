package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"

	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/internal/logger"
)

const blockLength = 128

type CryptService interface {
	Decrypt(data []byte) ([]byte, error)
	Encrypt(data []byte) ([]byte, error)
}

type cryptService struct {
	log        *zap.SugaredLogger
	privateKey *rsa.PrivateKey
}

func NewCryptService(privateKey *rsa.PrivateKey) CryptService {
	return &cryptService{
		log:        logger.NewLogger("crypt"),
		privateKey: privateKey,
	}
}

func (e *cryptService) Decrypt(data []byte) ([]byte, error) {
	if e.privateKey == nil {
		return data, nil
	}
	decryptedData := make([]byte, 0, len(data))
	var nextBlockLength int
	for i := 0; i < len(data); i += e.privateKey.PublicKey.Size() {
		nextBlockLength = i + e.privateKey.PublicKey.Size()
		if nextBlockLength > len(data) {
			nextBlockLength = len(data)
		}
		block, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, e.privateKey, data[i:nextBlockLength], []byte("yandex"))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data '%s': %v", data, err)
		}
		decryptedData = append(decryptedData, block...)
	}
	return decryptedData, nil
}

func (e *cryptService) Encrypt(data []byte) ([]byte, error) {
	if e.privateKey == nil {
		return data, nil
	}
	encryptedData := make([]byte, 0, len(data))
	var nextBlockLength int
	for i := 0; i < len(data); i += blockLength {
		nextBlockLength = i + blockLength
		if nextBlockLength > len(data) {
			nextBlockLength = len(data)
		}
		block, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &e.privateKey.PublicKey, data[i:nextBlockLength], []byte("yandex"))
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt data '%s': %v", data, err)
		}
		encryptedData = append(encryptedData, block...)
	}
	return encryptedData, nil
}

//func (e *cryptService) EncryptFile(data []byte) ([]byte, error) {
//	if e.publicKey == nil {
//		return data, nil
//	}
//	encrypted := make([]byte, 0, len(data))
//	var nextBlockLength int
//	for i := 0; i < len(data); i += blockLength {
//		nextBlockLength = i + blockLength
//		if nextBlockLength > len(data) {
//			nextBlockLength = len(data)
//		}
//		block, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, e.publicKey, data[i:nextBlockLength], []byte("yandex"))
//		if err != nil {
//			return nil, err
//		}
//		encrypted = append(encrypted, block...)
//	}
//	log.Printf("Encrypted data '%s'", string(data))
//	return encrypted, nil
//}
