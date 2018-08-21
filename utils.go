package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/wanchain/go-wanchain/crypto"
	wanclient "github.com/wanchain/go-wanchain/ethclient"

	"github.com/spf13/viper"
)

func getWanchainConnection() *wanclient.Client {
	uri := viper.GetString("nodeuri")

	client, err := wanclient.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// func getEthereumConnection() *ethclient.Client {
//	uri := viper.GetString("nodeuri")

//	client, err := ethclient.Dial(uri)
//	if err != nil {
//		log.Fatal(err)
//	}

//	return client
// }

func parseAbi(abiFileName string) ([]AbiField, error) {
	abiBytes, err := ioutil.ReadFile(abiFileName)
	if err != nil {
		return nil, err
	}

	fields := []AbiField{}

	if err := json.Unmarshal(abiBytes, &fields); err != nil {
		return nil, err
	}

	return fields, nil
}

func buildSignature(f *AbiField) (string, string) {
	inputTypes := []string{}
	for _, input := range f.Inputs {
		inputTypes = append(inputTypes, input.Type.String())
	}

	str := f.Name + "(" + strings.Join(inputTypes, ",") + ")"
	hash := crypto.Keccak256Hash([]byte(str))

	return str, hash.Hex()
}

var weiPerEth = big.NewInt(1000000000000000000)

func fromWei(i *big.Int) *big.Float {
	f := new(big.Float).SetInt(i)
	w := new(big.Float).SetInt(weiPerEth)

	return new(big.Float).Quo(f, w)
}
