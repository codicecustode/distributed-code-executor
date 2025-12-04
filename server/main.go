package main

import (
	"os"
	"os/signal"
	"syscall"
	"flag"
	"log"
	"time"
	"math/big"
	"encoding/pem"
	"context"
	
	"crypto/tls"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"crypto/rand"

	"github.com/quic-go/quic-go"
	
)


func main(){

	var (
			addr     = flag.String("addr", ":8443", "Server address")
			certFile = flag.String("cert", "", "TLS certificate file")
			keyFile  = flag.String("key", "", "TLS key file")
    )

  flag.Parse()

	// Generate self-signed certificate if not provided
	cert, key, err := getTLSConfig(*certFile, *keyFile)
	if err != nil {
			log.Fatalf("Failed to get TLS config: %v", err)
	}

	tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"distributed-code-executor"},
	}
	
	listener, err := quic.ListenAddr(*addr , tlsConfig, &quic.Config{

	})

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Due to OS Interruption, Shutting down server...")
		cancel()
	}()

	for {
		conn, err := listener.Accept(ctx)

		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
      continue
		}

		go HandleConnection(conn)
	}

}

func getTLSConfig(certFile, keyFile string) (tls.Certificate, *rsa.PrivateKey, error) {
    if certFile != "" && keyFile != "" {
        cert, err := tls.LoadX509KeyPair(certFile, keyFile)
        if err != nil {
            return tls.Certificate{}, nil, err
        }
        return cert, nil, nil
    }

    // Generate self-signed certificate
    priv, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return tls.Certificate{}, nil, err
    }

    template := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"Distributed Code Executor"},
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
    }

    certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
    if err != nil {
        return tls.Certificate{}, nil, err
    }

    certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
    keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

    cert, err := tls.X509KeyPair(certPEM, keyPEM)
    if err != nil {
        return tls.Certificate{}, nil, err
    }

    return cert, priv, nil
}