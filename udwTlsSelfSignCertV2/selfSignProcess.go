package udwTlsSelfSignCertV2

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/tachyon-protocol/udw/udwTime"
	"math/big"
	"net"
	"sync"
	"time"
)

type PemPair struct {
	Key  []byte
	Cert []byte
}

var gGetSelfSignProcessInitOnce sync.Once
var gGetSelfSignProcessInitPemPair PemPair
var gTlsConfig *tls.Config
var gCert *tls.Certificate

func getSelfSignProcessInit() {
	gGetSelfSignProcessInitOnce.Do(func() {
		startTime := time.Now()
		notBefore := startTime.Add(-100 * udwTime.Year)
		notAfter := startTime.Add(100 * udwTime.Year)
		template := x509.Certificate{
			SerialNumber: big.NewInt(1),

			NotBefore: notBefore,
			NotAfter:  notAfter,

			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}
		template.IPAddresses = []net.IP{net.IPv4(127, 0, 0, 1)}
		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			fmt.Println("rrndw9buaa", err)
			return
		}
		derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
		if err != nil {
			fmt.Println("sgmgkg25xg", err)
			return
		}
		certPem := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		})
		gGetSelfSignProcessInitPemPair.Cert = certPem
		b, err := x509.MarshalECPrivateKey(priv)
		if err != nil {
			fmt.Println("8pmrga2esq", err)
			return
		}
		privPem := pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: b,
		})
		gGetSelfSignProcessInitPemPair.Key = privPem
		cert, err := tls.X509KeyPair(gGetSelfSignProcessInitPemPair.Cert, gGetSelfSignProcessInitPemPair.Key)
		if err != nil {
			fmt.Println("3nb9k96z3p", err)
			return
		}
		gCert = &cert
		gTlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	})
}
func GetPemPair() PemPair {
	getSelfSignProcessInit()
	return gGetSelfSignProcessInitPemPair
}

func GetTlsCertificate() *tls.Certificate {
	getSelfSignProcessInit()
	return gCert
}

func GetTlsConfig() *tls.Config {
	getSelfSignProcessInit()
	return gTlsConfig
}
