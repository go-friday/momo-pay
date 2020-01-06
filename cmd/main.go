package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-friday/momo-pay"
	"golang.org/x/crypto/openpgp"
)

func readKey(path string) (openpgp.EntityList, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return openpgp.ReadArmoredKeyRing(file)
}

func mustReadKey(path string) openpgp.EntityList {
	keys, err := readKey(path)
	if err != nil {
		panic(err)
	}

	return keys
}

func main() {
	partnerCode := flag.String("partner", "", "Momo partner code")
	password := flag.String("password", "000000", "Wallet password")
	momoKeyFile := flag.String("momoKeyFile", "./momo.asc", "Momo public key")
	partnerKeysFile := flag.String("partnerKeysFile", "./partner.asc", "Partner public & private keys")
	flag.Parse()

	client := momo.NewClient(
		"https://payment.momo.vn:2445",
		*partnerCode,
		*password,
		mustReadKey(*partnerKeysFile),
		mustReadKey(*momoKeyFile),
	)

	res, err := client.GetBalance(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Waller balance:", res)
}
