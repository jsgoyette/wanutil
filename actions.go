package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/jsgoyette/wanutil/contracts"

	wanchain "github.com/wanchain/go-wanchain"
	"github.com/wanchain/go-wanchain/common"
	"github.com/wanchain/go-wanchain/core/types"
	"github.com/wanchain/go-wanchain/rlp"
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

		txData := tx.Data()
		txDataStr := fmt.Sprintf("%x", txData)

		if txDataStr != "" {
			for k := range methods {
				if strings.Contains(txDataStr, k[2:10]) {
					method := methods[k]
					printMethod(&method, txData[4:])
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
					printEvent(rlog.Address.Hex(), &method, rlog)
				}
			}
		}

		printReceipt(receipt)
	}

	return nil
}

func getBlock(c *cli.Context) error {
	blockNumber := c.Int64("block")
	blockHash := c.String("hash")

	if blockNumber == 0 && blockHash == "" {
		return cli.NewExitError("Either block number or block hash must be provided", 1)
	}
	if blockNumber != 0 && blockHash != "" {
		return cli.NewExitError("Ambiguous: only a block number or a block hash should be provided", 1)
	}

	client := getWanchainConnection()

	var block *types.Block
	var err error

	if blockNumber != 0 {
		block, err = client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	} else {
		block, err = client.BlockByHash(context.Background(), common.HexToHash(blockHash))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	fmt.Println(block)

	return nil
}

func listTransactionsToAddress(c *cli.Context) error {
	return scanBlockTransactions(c, "to")
}

func listTransactionsFromAddress(c *cli.Context) error {
	return scanBlockTransactions(c, "from")
}

func scanBlockTransactions(c *cli.Context, direction string) error {
	address := strings.ToLower(c.String("address"))

	if address == "" {
		return cli.NewExitError("No address provided", 1)
	}

	client := getWanchainConnection()
	networkId, _ := client.NetworkID(context.Background())
	signer := types.NewEIP155Signer(networkId)

	startingBlock := c.Int64("block")

	if startingBlock == 0 {
		startingBlock = 1
	}

	current, err := currentBlockNumber(client)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("Block   | Hash")
	fmt.Println(strings.Repeat("-", 76))

	for c := startingBlock; c <= current.Int64(); c++ {

		block, err := client.BlockByNumber(context.Background(), big.NewInt(c))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for _, tx := range block.Transactions() {
			var txaddr *common.Address

			if direction == "to" {
				txaddr = tx.To()

				// check receipt contract address
				if txaddr == nil {

					receipt, err := client.TransactionReceipt(
						context.Background(),
						tx.Hash(),
					)
					if err != nil {
						return cli.NewExitError(err.Error(), 1)
					}

					txaddr = &receipt.ContractAddress
				}
			} else if direction == "from" {
				if msg, err := tx.AsMessage(signer); err == nil {
					from := msg.From()
					txaddr = &from
				}
			}

			if txaddr != nil && address == strings.ToLower(txaddr.Hex()) {
				fmt.Printf("%7d | %s\n", c, tx.Hash().Hex())
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
	inputNames := map[string]string{}

	for _, field := range fields {
		if field.Name == "" {
			continue
		}

		sig, sigHash := buildSignature(&field)
		names := getInputNamesString(field.Inputs)

		signatures[sigHash] = sig
		inputNames[sigHash] = names
	}

	keys := make([]string, 0, len(signatures))
	for k := range signatures {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s %s %s\n", k, signatures[k], inputNames[k])
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

func subscribe(c *cli.Context) error {
	addrString := c.String("address")

	if addrString == "" {
		return cli.NewExitError("No address provided", 1)
	}

	client := getWanchainConnection()
	startingBlock := c.Int64("block")

	address := common.HexToAddress(addrString)
	query := wanchain.FilterQuery{
		FromBlock: big.NewInt(startingBlock),
		Addresses: []common.Address{address},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for {
		select {
		case err := <-sub.Err():
			return cli.NewExitError(err.Error(), 1)
		case vLog := <-logs:
			printLog(&vLog)
		}
	}

	return nil
}

func decodeTransaction(c *cli.Context) error {
	hexString := c.String("hex")

	if hexString == "" {
		return cli.NewExitError("No hex string provided", 1)
	}

	var tx *types.Transaction
	var from string

	rawtx, err := hex.DecodeString(hexString)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	rlp.DecodeBytes(rawtx, &tx)

	client := getWanchainConnection()
	networkId, _ := client.NetworkID(context.Background())
	signer := types.NewEIP155Signer(networkId)

	if msg, err := tx.AsMessage(signer); err == nil {
		from = msg.From().Hex()
	}

	printTransaction(tx, from, false)

	return nil
}
