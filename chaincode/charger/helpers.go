package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"github.com/hyperledger/fabric/common/attrmgr"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/identity"
)

func SerializeIdentity(mspId string, certificate *x509.Certificate) string {
	return mspId + ":" + GetCertificateSubject(certificate)
}

func GetCertificateSubject(certificate *x509.Certificate) string {
	return certificate.Subject.CommonName
}

func GetCertificateAttribute(identity *identity.CertIdentity, attr string) (string, error) {
	manager := attrmgr.New()
	attributes, err := manager.GetAttributesFromCert(identity.Cert)
	if err != nil {
		return "", errors.New("unable to get attributes from certificate")
	}

	value, ok, err := attributes.Value(attr)
	if err != nil || ok == false {
		return "", errors.New("unable to get attributes from certificate")
	}

	return value, nil
}


func IdentityIsAdmin(identity *identity.CertIdentity) bool {
	val, _ := GetCertificateAttribute(identity, "role")
	return val == "admin"
}

func IdentityIsUser(identity *identity.CertIdentity) bool {
	val, _ := GetCertificateAttribute(identity, "role")
	return val == "user"
}

func IdentityIsEqual(identity *identity.CertIdentity, mspId, username string) bool {
	current := SerializeIdentity(identity.MspID, identity.Cert)
	given := mspId + ":" + username
	return current == given
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
