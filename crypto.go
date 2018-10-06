package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
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
