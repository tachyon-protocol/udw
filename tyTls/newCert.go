package tyTls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

type NewTlsCertResp struct{
	TlsCert tls.Certificate
	CertPem string
	PkPem string
}

// can not crash by user
func MustNewTlsCert(isClient bool) (resp NewTlsCertResp) {
	var ExtKeyUsage x509.ExtKeyUsage
	if isClient {
		ExtKeyUsage = x509.ExtKeyUsageClientAuth
	} else {
		ExtKeyUsage = x509.ExtKeyUsageServerAuth
	}
	const dur = 100 * 365 * 24 * time.Hour
	startTime := time.Now()
	notBefore := startTime.Add(-dur)
	notAfter := startTime.Add(dur)
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{ExtKeyUsage},
		BasicConstraintsValid: true,
	}
	if isClient == false {
		template.IPAddresses = []net.IP{net.IPv4(127, 0, 0, 1)}
	}
	//template.DNSNames = []string{"google.com"}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		panic(err)
	}
	resp.CertPem = string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}))
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		panic(err)
	}
	resp.PkPem = string(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: b,
	}))
	tlsCert,err:=tls.X509KeyPair([]byte(resp.CertPem), []byte(resp.PkPem))
	if err != nil {
		panic(err)
	}
	resp.TlsCert = tlsCert
	return resp
}

func MustNewTlsCertSimple(isClient bool) (cert tls.Certificate) {
	resp:= MustNewTlsCert(isClient)
	return resp.TlsCert
}



