package main

import (
	"crypto/x509"
)

func SerializeIdentity(mspId string, certificate *x509.Certificate) string {
	return mspId + ":" + certificate.Subject.CommonName
}
