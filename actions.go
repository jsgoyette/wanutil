package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/jsgoyette/wanutil/contracts"
	"github.com/wanchain/go-wanchain/common"
	wanclient "github.com/wanchain/go-wanchain/ethclient"
)

var ZERO = big.NewInt(0)

func currentBlockNumber(client *wanclient.Client) (*big.Int, error) {
	latestBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return latestBlock.Number(), nil
}

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
	hash := c.String("hash")

	if hash == "" {
		return cli.NewExitError("No tx hash provided", 1)
	}

	client := getWanchainConnection()
	chash := common.HexToHash(hash)

	tx, isPending, err := client.TransactionByHash(
		context.Background(),
		chash,
	)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Transaction:\n%+v\n", tx)
	fmt.Printf("Pending: %v\n\n", isPending)

	if !isPending {
		receipt, err := client.TransactionReceipt(
			context.Background(),
			chash,
		)

		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Printf("Receipt:\n%+v\n", receipt)
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
