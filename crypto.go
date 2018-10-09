package main

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

var publicKey *rsa.PublicKey = nil

func InitializePublicKey() {
	data, err := ioutil.ReadFile("./publickey")
	if err != nil {
		panic(err)
	}
	decoded, _ := pem.Decode(data)
	pub, err := x509.ParsePKIXPublicKey(decoded.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey = pub.(*rsa.PublicKey)

}

func VerifySignature(request []byte, sign []byte) bool {
	hash := sha512.Sum512(request)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hash[:], sign)
	return err == nil
}

func encryptCBC(symmetricKey []byte, inBytes []byte) ([]byte, error) {
	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, err
	}

	padLen := aes.BlockSize - len(inBytes)%aes.BlockSize
	padding := make([]byte, padLen)
	_, err = rand.Reader.Read(padding)
	if err != nil {
		return nil, err
	}
	padding[0] = byte(padLen)
	inBytes = append(padding, inBytes...)
	inBytesLen := len(inBytes)
	if inBytesLen%aes.BlockSize != 0 {
		return nil, errors.New("bad padding")
	}

	ciphertext := make([]byte, aes.BlockSize+inBytesLen)
	initializationVector := ciphertext[:aes.BlockSize]
	_, err = rand.Reader.Read(initializationVector)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCBCEncrypter(block, initializationVector)
	cfb.CryptBlocks(ciphertext[aes.BlockSize:], inBytes)

	return ciphertext, nil
}

func encrypt(data []byte) ([]byte, error) {
	randKey := make([]byte, 32)
	rand.Read(randKey)

	encrypted, err := encryptCBC(randKey, data)
	if err != nil {
		return nil, err
	}

	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, randKey)
	if err != nil {
		return nil, err
	}

	result := append(encryptedKey[:], encrypted[:]...)

	return result, nil
}
