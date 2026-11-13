package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dir := ".local/keys"
	if err := os.MkdirAll(dir, 0700); err != nil {
		panic(err)
	}
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	privateDER, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		panic(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		panic(err)
	}
	for _, f := range []struct {
		name, kind string
		der        []byte
		mode       os.FileMode
	}{{"api-private.pem", "PRIVATE KEY", privateDER, 0600}, {"api-public.pem", "PUBLIC KEY", publicDER, 0644}} {
		target := filepath.Join(dir, f.name)
		file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, f.mode)
		if err != nil {
			panic(err)
		}
		if err := pem.Encode(file, &pem.Block{Type: f.kind, Bytes: f.der}); err != nil {
			_ = file.Close()
			panic(err)
		}
		if err := file.Close(); err != nil {
			panic(err)
		}
	}
	fmt.Println("Keys generated in .local/keys/. Never commit the private key.")
}
