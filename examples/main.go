package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/wallet"
)

func main() {
	wallets, err := wallet.InitializeWallets()
	if err != nil {
		panic(err)
	}

	// id := wallets.AddWallet("y")
	// wallets.SaveFile()
	w, _ := wallets.GetWallet("y")
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(&w.Main.PrivateKey.PublicKey)
	pemEncodedMainPub := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: x509EncodedPub,
		})

	x509EncodedPub, _ = x509.MarshalPKIXPublicKey(&w.Main.PrivateKey.PublicKey)
	pemEncodedViewPub := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: w.Main.PublicKey,
	})

	// return string(pemEncoded), string(pemEncodedPub)
	fmt.Println(string(pemEncodedViewPub), string(pemEncodedMainPub))
	fmt.Println(string(w.Certificate[:]))
}

// -----BEGIN CERTIFICATE-----
// MIIBGTCBv6ADAgECAgEBMAoGCCqGSM49BAMCMBQxEjAQBgNVBAoTCURJRCwgSW5j
// LjAeFw0yMDA3MDMwOTE4MjZaFw0yMjA3MDMwOTE4MjZaMBQxEjAQBgNVBAoTCURJ
// RCwgSW5jLjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABM8r/BIahPrFJjWrZ2sf
// Bke3iFa3LHRUoDVpgmFCJNBdodhMOFMEbaoCcBxkrKHkYaUpxCsbcOZt+5v9zx7i
// SeujAjAAMAoGCCqGSM49BAMCA0kAMEYCIQDCYCXoc67CRoI2HwTdt5U7Qr2zn0Zg
// Ct9QrQCeLy4+YgIhAOzHWa9K20ZmovGMAEF7s+Sfk3fhOchDKKBayvz0uU5p
// -----END CERTIFICATE-----
