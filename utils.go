package main

import (
	"log"
	"math/big"

	// "github.com/ethereum/go-ethereum/ethclient"
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

var weiPerEth = big.NewInt(1000000000000000000)

func fromWei(i *big.Int) *big.Float {
	f := new(big.Float).SetInt(i)
	w := new(big.Float).SetInt(weiPerEth)

	return new(big.Float).Quo(f, w)
}
