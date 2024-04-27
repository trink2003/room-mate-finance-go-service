package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not rsa")
}

func main() {

	// Create the keys
	priv, pub := GenerateRsaKeyPair()

	// Export the keys to pem string
	priv_pem := ExportRsaPrivateKeyAsPemStr(priv)
	pub_pem, _ := ExportRsaPublicKeyAsPemStr(pub)

	// Import the keys from pem string
	priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
	pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)

	// Export the newly imported keys
	priv_parsed_pem := ExportRsaPrivateKeyAsPemStr(priv_parsed)
	pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)

	fmt.Println(priv_parsed_pem)
	fmt.Println(pub_parsed_pem)

	// Check that the exported/imported keys match the original keys
	if priv_pem != priv_parsed_pem || pub_pem != pub_parsed_pem {
		fmt.Println("Failure: Export and Import did not result in same Keys")
	} else {
		fmt.Println("Success")
	}

	encryptedBytes, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		pub_parsed,
		[]byte("super secret message"),
		nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("encrypted bytes: ", encryptedBytes)
	fmt.Printf("encrypted string: %x\n", encryptedBytes)

	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OEAPOptions in the end signify that we encrypted the data using OEAP, and that we used
	// SHA256 to hash the input.
	decryptedBytes, err := priv.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA512})
	if err != nil {
		log.Fatalln(err)
	}

	// We get back the original information in the form of bytes, which we
	// the cast to a string and print
	fmt.Println("decrypted message: ", string(decryptedBytes))

	msg := []byte("verifiable message")

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha512.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		log.Fatalln(err)
	}
	msgHashSum := msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA512, msgHashSum, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(pub_parsed, crypto.SHA512, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	}
	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")

}
