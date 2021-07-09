package main

import (
	"fmt"

	"github.com/workspace/evoting/ev-webservice/wallet"
)

func main() {
	wallets, err := wallet.InitializeWallets()
	if err != nil {
		panic(err)
	}

	id := wallets.AddWallet("y")
	wallets.Save()
	w, _ := wallets.GetWallet(id)

	fmt.Println(string(w.Certificate[:]))
	fmt.Println(string(wallet.Base58Encode(w.Main.PublicKey)))
	fmt.Println(w.View.PublicKey)
	x := string(wallet.Base58Encode(w.View.PublicKey))
	fmt.Println(x)
	fmt.Println(wallet.Base58Decode([]byte(x)))
}

// -----BEGIN CERTIFICATE-----
// MIIBGTCBv6ADAgECAgEBMAoGCCqGSM49BAMCMBQxEjAQBgNVBAoTCURJRCwgSW5j
// LjAeFw0yMDA3MDMwOTE4MjZaFw0yMjA3MDMwOTE4MjZaMBQxEjAQBgNVBAoTCURJ
// RCwgSW5jLjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABM8r/BIahPrFJjWrZ2sf
// Bke3iFa3LHRUoDVpgmFCJNBdodhMOFMEbaoCcBxkrKHkYaUpxCsbcOZt+5v9zx7i
// SeujAjAAMAoGCCqGSM49BAMCA0kAMEYCIQDCYCXoc67CRoI2HwTdt5U7Qr2zn0Zg
// Ct9QrQCeLy4+YgIhAOzHWa9K20ZmovGMAEF7s+Sfk3fhOchDKKBayvz0uU5p
// -----END CERTIFICATE-----
