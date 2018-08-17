package main

import (
	"github.com/urfave/cli"
)

var (
	addressFlag = cli.StringFlag{
		Name:  "address, a",
		Value: "",
		Usage: "Address hash",
	}
	blockFlag = cli.IntFlag{
		Name:  "block, b",
		Value: 0,
		Usage: "Block number",
	}
	countFlag = cli.IntFlag{
		Name:  "count, c",
		Value: 20,
		Usage: "Record count",
	}
	hashFlag = cli.StringFlag{
		Name:  "hash",
		Value: "",
		Usage: "Hash string",
	}
	tokenFlag = cli.StringFlag{
		Name:  "token, t",
		Value: "",
		Usage: "Token name",
	}

	commands = []cli.Command{
		{
			Name:        "getBalance",
			Aliases:     []string{"bal"},
			Usage:       "Get address balance",
			UsageText:   "wanutil getBalance [options]",
			Description: "Get the balance or token balance for an address. To get the token balance, make sure to set the token address in your config file.",
			Action:      getBalance,
			Flags:       []cli.Flag{addressFlag, blockFlag, tokenFlag},
		},
		{
			Name:        "getTransaction",
			Aliases:     []string{"tx"},
			Usage:       "Get transaction by hash",
			UsageText:   "wanutil getTransaction [options]",
			Description: "Get transaction details and receipt. If it is to a recognized smart contract (in your config file) it will also try to parse the input.",
			Action:      getTransaction,
			Flags:       []cli.Flag{hashFlag},
		},
		{
			Name:        "listTransactionsToAddress",
			Aliases:     []string{"atx"},
			Usage:       "Get transactions sent to a given address",
			UsageText:   "wanutil listTransactionsToAddress [options]",
			Description: "Get the transactions sent to a given address, using an optional block number range",
			Action:      listTransactionsToAddress,
			Flags:       []cli.Flag{addressFlag, blockFlag},
		},
	}
)
