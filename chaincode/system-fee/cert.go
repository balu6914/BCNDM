package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
)

var errInvalidCert = errors.New("Failed to parse PEM certificate")

func parsePEM(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, errInvalidCert
	}

	return x509.ParseCertificate(block.Bytes)
}

func cnFromX509(certPEM string) (string, error) {
	cert, err := parsePEM(certPEM)
	if err != nil {
		return "", err
	}
	return cert.Subject.CommonName, nil
}

func callerCN(stub shim.ChaincodeStubInterface) (string, error) {
	data, err := stub.GetCreator()
	if err != nil {
		return "", err
	}

	var serializedID msp.SerializedIdentity
	if err := proto.Unmarshal(data, &serializedID); err != nil {
		return "", err
	}

	cn, err := cnFromX509(string(serializedID.IdBytes))
	if err != nil {
		return "", err
	}

	return cn, nil
}
