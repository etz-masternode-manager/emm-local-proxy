package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"log"
	"math/big"
	"os"
)

var publicKey *rsa.PublicKey = nil

func readLength(data []byte) ([]byte, uint32, error) {
	lBuf := data[0:4]
	buf := bytes.NewBuffer(lBuf)
	var length uint32
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return nil, 0, err
	}

	return data[4:], length, nil
}

func readBigInt(data []byte, length uint32) ([]byte, *big.Int, error) {
	var bigint = new(big.Int)
	bigint.SetBytes(data[0:length])
	return data[length:], bigint, nil
}

func getRsaValues(data []byte) (format string, e *big.Int, n *big.Int, err error) {
	data, length, err := readLength(data)
	if err != nil {
		return
	}
	format = string(data[0:length])
	data = data[length:]
	data, length, err = readLength(data)
	if err != nil {
		return
	}
	data, e, err = readBigInt(data, length)
	if err != nil {
		return
	}
	data, length, err = readLength(data)
	if err != nil {
		return
	}
	data, n, err = readBigInt(data, length)
	if err != nil {
		return
	}

	return
}

func InitializePublicKey() {
	data, err := base64.StdEncoding.DecodeString("AAAAB3NzaC1yc2EAAAABJQAAAQEAoIvE/jIeRiVqaSGuSmZfqjFPaWpUcB7fF+S2mo/aXV/ht4Jl1IIJniEKwH2tA4VQPdmKH0w8kyvI+bZIXsy8Xi6ziMeb8w8NJO5PEwboXJG3cg+tdlopj5g9iXevSq1i4XifVezT56xeslHwhr0f5+NtVzvMCoIHg5c3D19UBOwEtQk27RLQ7K1qNZ2bjNB5x1tmYH7+vJSIJr6UO/6ewAnAFrsCivlPhl7RT1/euPQjnNuVEuMFBEkhDAupyowxxkXNgXbV8jJDSOvn7KPze1x+QkT40jX3MsWdgDpLbgQVLE7t4Df0teAO45qfbXj/By/EEiG+je0zhrqqW2cdBw==")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	_, e, n, err := getRsaValues(data)

	publicKey = &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}
}

func VerifySignature(request []byte, sign []byte) bool {
	hash := sha512.Sum512(request)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hash[:], sign)
	return err == nil
}
