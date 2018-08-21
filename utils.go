package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/wanchain/go-wanchain/core/types"
	"github.com/wanchain/go-wanchain/crypto"
	wanclient "github.com/wanchain/go-wanchain/ethclient"
	// "github.com/ethereum/go-ethereum/ethclient"

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

func currentBlockNumber(client *wanclient.Client) (*big.Int, error) {
	latestBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return latestBlock.Number(), nil
}

func printTransaction(tx *types.Transaction, from string, isPending bool) {
	v, r, s := tx.RawSignatureValues()

	fmt.Printf("Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("To: %s\n", tx.To().Hex())
	fmt.Printf("From: %s\n", from)
	fmt.Printf("TxType: 0x%x\n", tx.Txtype())
	fmt.Printf("Value: %s\n", tx.Value().String())
	fmt.Printf("Gas: %s\n", tx.Gas())
	fmt.Printf("Gas Price: %d\n", tx.GasPrice().Uint64())
	fmt.Printf("Nonce: %d\n", tx.Nonce())
	fmt.Printf("Size: %s\n", tx.Size().String())
	fmt.Printf("Data: %x\n", tx.Data())
	fmt.Printf("V: 0x%x\n", v)
	fmt.Printf("R: 0x%x\n", r)
	fmt.Printf("S: 0x%x\n", s)
	fmt.Printf("Pending: %v\n\n", isPending)
}

func printReceipt(receipt *types.Receipt) {
	fmt.Printf("Status: %d\n", receipt.Status)
	fmt.Printf("Cumulative Gas Used: %s\n", receipt.CumulativeGasUsed.String())
	fmt.Printf("Bloom: %x\n", receipt.Bloom)
	fmt.Printf("Logs:\n")

	for _, log := range receipt.Logs {
		fmt.Printf("\tAddress: %s\n", log.Address.Hex())
		fmt.Printf("\tBlock Hash: %s\n", log.BlockHash.Hex())
		fmt.Printf("\tBlock Number: %d\n", log.BlockNumber)
		fmt.Printf("\tRemoved: %v\n", log.Removed)
		fmt.Printf("\tData: %x\n", log.Data)
		fmt.Printf("\tTopics:\n")
		for _, topic := range log.Topics {
			fmt.Printf("\t\t%x\n", topic)
		}
		fmt.Println()
	}
}

func printMethod(method *AbiMethod) {
	inputNames := make([]string, len(method.Inputs))
	for i, input := range method.Inputs {
		inputNames[i] = input.Name
	}
	inputString := strings.Join(inputNames, ", ")

	fmt.Println("Method:", method.Name)
	fmt.Println("Signature:", method.Signature)
	fmt.Println("Inputs:", inputString)
	fmt.Println()
}

func printEvent(address string, method *AbiMethod) {
	inputNames := make([]string, len(method.Inputs))
	for i, input := range method.Inputs {
		inputNames[i] = input.Name
	}
	inputString := strings.Join(inputNames, ", ")

	fmt.Println("Event:", method.Name)
	fmt.Println("Address:", address)
	fmt.Println("Signature:", method.Signature)
	fmt.Println("Inputs:", inputString)
	fmt.Println()
}
