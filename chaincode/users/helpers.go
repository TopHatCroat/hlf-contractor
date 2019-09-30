package main

import (
	"crypto/x509"
	"errors"
	"github.com/hyperledger/fabric/common/attrmgr"
	"github.com/s7techlab/cckit/identity"
)

func SerializeIdentity(mspId string, certificate *x509.Certificate) string {
	return mspId + ":" + certificate.Subject.CommonName
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

func GetCertificateSubject(certificate *x509.Certificate) string {
	return certificate.Subject.CommonName
}
