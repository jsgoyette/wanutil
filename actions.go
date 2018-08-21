package main

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/jsgoyette/wanutil/contracts"

	"github.com/wanchain/go-wanchain/common"
	"github.com/wanchain/go-wanchain/core/types"
)

var ZERO = big.NewInt(0)

func getBalance(c *cli.Context) error {

	address := c.String("address")

	if address == "" {
		return cli.NewExitError("No address provided", 1)
	}

	client := getWanchainConnection()
	blockNumber := big.NewInt(c.Int64("block"))

	tokenSymbol := c.String("token")
	tokenAddress := viper.GetString("contracts." + tokenSymbol)

	if tokenSymbol != "" && tokenAddress == "" {
		return cli.NewExitError("Token not found", 1)
	}

	if blockNumber.Cmp(ZERO) == 0 {
		current, err := currentBlockNumber(client)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		blockNumber = current
	}

	if tokenAddress != "" {

		// get token contract instance
		ta := common.HexToAddress(tokenAddress)
		instance, err := contracts.NewStandard(ta, client)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		// check token balance on contract
		balance, err := instance.BalanceOf(
			nil,
			common.HexToAddress(address),
		)

		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Printf(
			"%s balance: %s (%s)\n",
			tokenSymbol,
			balance.String(),
			fromWei(balance).String(),
		)

	} else {

		// check WAN balance
		balance, err := client.BalanceAt(
			context.Background(),
			common.HexToAddress(address),
			blockNumber,
		)

		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Printf(
			"Balance at block %d: %s (%s)\n",
			blockNumber,
			balance.String(),
			fromWei(balance).String(),
		)
	}

	return nil
}

func getTransaction(c *cli.Context) error {
	hexHash := c.String("hash")
	abiFileName := c.String("abi")

	if hexHash == "" {
		return cli.NewExitError("No tx hash provided", 1)
	}

	client := getWanchainConnection()
	networkId, _ := client.NetworkID(context.Background())

	hash := common.HexToHash(hexHash)
	signer := types.NewEIP155Signer(networkId)

	methods := map[string]AbiMethod{}
	var from string

	tx, isPending, err := client.TransactionByHash(
		context.Background(),
		hash,
	)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if msg, err := tx.AsMessage(signer); err == nil {
		from = msg.From().Hex()
	}

	printTransaction(tx, from, isPending)

	if abiFileName != "" {

		fields, err := parseAbi(abiFileName)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for _, field := range fields {
			if field.Name == "" {
				continue
			}

			sig, sigHash := buildSignature(&field)
			methods[sigHash] = AbiMethod{
				AbiField:      field,
				Signature:     sig,
				SignatureHash: sigHash,
			}
		}

		txData := "0x" + fmt.Sprintf("%x", tx.Data())
		if txData != "" {
			for k := range methods {
				if strings.Contains(txData, k[:10]) {
					method := methods[k]
					printMethod(&method)
				}
			}
		}
	}

	if !isPending {
		receipt, err := client.TransactionReceipt(
			context.Background(),
			hash,
		)

		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if abiFileName != "" {

			for _, rlog := range receipt.Logs {
				if len(rlog.Topics) == 0 {
					continue
				}

				if method, ok := methods[rlog.Topics[0].String()]; ok {
					printEvent(rlog.Address.Hex(), &method)
				}
			}
		}

		printReceipt(receipt)
	}

	return nil
}

func listTransactionsToAddress(c *cli.Context) error {
	address := c.String("address")

	if address == "" {
		return cli.NewExitError("No address provided", 1)
	}

	client := getWanchainConnection()
	startingBlock := c.Int64("block")

	if startingBlock == 0 {
		startingBlock = 1
	}

	current, err := currentBlockNumber(client)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for c := startingBlock; c <= current.Int64(); c++ {

		block, err := client.BlockByNumber(context.Background(), big.NewInt(c))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for _, tx := range block.Transactions() {
			txaddr := tx.To()
			if txaddr != nil && address == txaddr.Hex() {
				fmt.Printf("block: %7d hash: %s\n", c, tx.Hash().Hex())
			}
		}
	}

	return nil
}

func listAbiSignatures(c *cli.Context) error {
	abiFileName := c.String("abi")
	if abiFileName == "" {
		return cli.NewExitError("ABI file path is required", 1)
	}

	fields, err := parseAbi(abiFileName)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	signatures := map[string]string{}

	for _, field := range fields {
		if field.Name == "" {
			continue
		}

		sig, sigHash := buildSignature(&field)
		signatures[sig] = sigHash
	}

	keys := make([]string, 0, len(signatures))
	for k := range signatures {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%-40s %s\n", k, signatures[k])
	}

	return nil
}

func validateAddress(c *cli.Context) error {
	address := c.String("address")

	if address == "" {
		return cli.NewExitError("No address provided", 1)
	}

	if !common.IsHexAddress(address) {
		fmt.Println("Address is INVALID")
	} else {
		addr := common.HexToAddress(address)
		fmt.Printf("Valid address: %s\n", addr.Hex())
	}

	return nil
}
