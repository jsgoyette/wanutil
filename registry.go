package main

import (
	"github.com/urfave/cli"
)

var (
	abiFileFlag = cli.StringFlag{
		Name:  "abi",
		Value: "",
		Usage: "ABI file name",
	}
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
			Flags:       []cli.Flag{abiFileFlag, hashFlag},
		},
		{
			Name:        "listTransactionsToAddress",
			Aliases:     []string{"scan"},
			Usage:       "List transactions sent to a given address",
			UsageText:   "wanutil listTransactionsToAddress [options]",
			Description: "List the transactions sent to a given address, using an optional block number range.",
			Action:      listTransactionsToAddress,
			Flags:       []cli.Flag{addressFlag, blockFlag},
		},
		{
			Name:        "abiSignatures",
			Aliases:     []string{"sig"},
			Usage:       "Get ABI method/event signatures",
			UsageText:   "wanutil signatures [options]",
			Description: "Get the signature hashes for the methods and events for a given ABI.",
			Action:      listAbiSignatures,
			Flags:       []cli.Flag{abiFileFlag},
		},
		{
			Name:        "validateAddress",
			Aliases:     []string{"validate"},
			Usage:       "Validate address with checksum",
			UsageText:   "wanutil validateAddress [options]",
			Description: "Check that an address is valid and return with correct capitalizations based on checksum.",
			Action:      validateAddress,
			Flags:       []cli.Flag{addressFlag},
		},
	}
)
