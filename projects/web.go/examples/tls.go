package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/hkcoldtea/src/projects/web.go/server"
)

type cmdlineArgs struct {
	ListenIP   string
	ListenPort int
	TLS        bool
	Cert       string
	Key        string
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

// This generates a self-signed cert good for 6 hours using ecdsa P256 curve
func makeSelfSigned(cert, key string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %s", err)
	}

	// these are temporary certs, only valid for 6 hours
	notBefore := time.Now()
	notAfter := notBefore.Add(6 * time.Hour)

	var priv interface{}
	priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %s", err)
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %s", err)
	}
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Random Bits UNL"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{hostname},
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	certOut, err := os.Create(cert)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %s", cert, err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %s", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing %s: %s", cert, err)
	}
	log.Printf("wrote %s\n", cert)

	keyOut, err := os.OpenFile(key, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %s", key, err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to %s: %s", key, err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing %s: %s", key, err)
	}
	log.Printf("wrote %s\n", key)
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

	var cfg = cmdlineArgs{
		ListenIP:   "0.0.0.0",
		ListenPort: 9999,
		TLS:        false,
		Cert:       "/var/tmp/server.crt",
		Key:        "/var/tmp/server.key",
	}

	flag.StringVar(&cfg.ListenIP, "ip", cfg.ListenIP, "IP Address to Listen to (0.0.0.0)")
	flag.IntVar(&cfg.ListenPort, "port", cfg.ListenPort, "Port to listen to (9999)")
	flag.BoolVar(&cfg.TLS, "tls", cfg.TLS, "Use https instead of http")
	flag.StringVar(&cfg.Cert, "cert", cfg.Cert, "Path to temporary server.crt")
	flag.StringVar(&cfg.Key, "key", cfg.Key, "Path to temporary server.key")

	flag.Parse()

	s := server.InitServer()

	s.Get("/*m", func(c *server.Context) {
		c.SetHeader("X-DEBUG-PATH", c.Path)
		for j, k := range c.PathParams {
			c.SetHeader("X-DEBUG-PATH-"+j, k)
		}
		c.JSON(http.StatusOK, nil)
	})
	listen := fmt.Sprintf("%s:%d", cfg.ListenIP, cfg.ListenPort)

	if cfg.TLS {
		makeSelfSigned(cfg.Cert, cfg.Key)
		log.Fatal(s.RunTLS(listen, cfg.Cert, cfg.Key))
	} else {
		log.Fatal(s.Run(listen))
	}
}
