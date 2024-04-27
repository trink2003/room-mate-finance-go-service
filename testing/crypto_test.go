package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"room-mate-finance-go-service/utils"
	"strings"
	"testing"
)

func TestCheckSumFunction(t *testing.T) {
	var testPlainTextNoUpperCase = "có thể gọi anh là đẹp trai nhất việt nam"
	var testPlainTextUpperCase = "CÓ THỂ GỌI ANH LÀ ĐẸP TRAI NHẤT VIỆT NAM"

	var hashedNoUpperCase, hashedNoUpperCaseError = utils.GenerateCheckSumUsingSha256New(testPlainTextNoUpperCase)
	var hashedUpperCase, hashedUpperCaseError = utils.GenerateCheckSumUsingSha256New(testPlainTextUpperCase)

	if hashedNoUpperCaseError != nil {
		t.Errorf("hashedNoUpperCaseError is not expected: %s", hashedNoUpperCaseError.Error())
		return
	}
	if hashedUpperCaseError != nil {
		t.Errorf("hashedUpperCaseError is not expected: %s", hashedUpperCaseError.Error())
		return
	}

	if strings.Compare(hashedNoUpperCase, hashedUpperCase) == 0 {
		t.Errorf("check sum does not work like expectation")
		return
	}

	fmt.Printf("%s\n", hashedNoUpperCase)
	fmt.Printf("%s\n", hashedUpperCase)
}

func TestAes(t *testing.T) {
	// bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	bytes := []byte("ChM1tVFwO6hWWhv6nCqNjPwSftHKPE4j")
	plaintext := "có thể gọi anh là đẹp trai nhất việt nam"
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault

	var encryptedText, encryptError = utils.AesEncryption(plaintext, key)
	if encryptError != nil {
		t.Errorf("encryptError is not expected: %s", encryptError.Error())
		return
	}

	if encryptedText == "" {
		t.Errorf("empty encryptedText is not expected")
		return
	}

	decryptedText, decryptError := utils.AesDecryption(encryptedText, key)
	if decryptError != nil {
		t.Errorf("decryptError is not expected: %s", decryptError.Error())
		return
	}

	if decryptedText == "" {
		t.Errorf("empty decryptedText is not expected")
		return
	}

	if strings.Compare(plaintext, decryptedText) != 0 {
		t.Errorf("ecrypted text does not same with plain text")
		return
	}

	fmt.Printf("\n\t- plainText: %s\n\t- encryptedText: %s\n\t- decryptedText: %s", plaintext, encryptedText, decryptedText)

}
