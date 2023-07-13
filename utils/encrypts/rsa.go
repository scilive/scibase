package encrypts

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

type RSA struct{}

func (RSA) OAEPEncrypt(secretMessage []byte, key rsa.PublicKey) (string, error) {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (RSA) OAEPDecrypt(cipherText string, privKey rsa.PrivateKey) ([]byte, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func (RSA) GenPem(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	publicKey := privateKey.PublicKey
	privKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))
	pubKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&publicKey),
	}))
	return privKey, pubKey, nil

}

func (RSA) ParsePublicKeyFromPEM(pubPEM string) (rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return rsa.PublicKey{}, fmt.Errorf("failed to parse PEM block containing the key")
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return rsa.PublicKey{}, err
	}
	return *pub, nil
}

func (RSA) ParsePrivateKeyFromPEM(privPEM string) (rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return rsa.PrivateKey{}, fmt.Errorf("failed to parse PEM block containing the key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return rsa.PrivateKey{}, err
	}
	return *priv, nil
}
