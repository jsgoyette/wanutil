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
			Name:        "balance",
			Aliases:     []string{"bal"},
			Usage:       "Get address balance",
			UsageText:   "wanutil balance [options]",
			Description: "Get the balance or token balance for an address. To get the token balance, make sure to set the token address in your config file.",
			Action:      getBalance,
			Flags:       []cli.Flag{addressFlag, blockFlag, tokenFlag},
		},
		{
			Name:        "block",
			Aliases:     []string{"blk"},
			Usage:       "Get block",
			UsageText:   "wanutil block [options]",
			Description: "Fetch a block by blocknumber or hash",
			Action:      getBlock,
			Flags:       []cli.Flag{blockFlag, hashFlag},
		},
		{
			Name:        "transaction",
			Aliases:     []string{"tx"},
			Usage:       "Get transaction by hash",
			UsageText:   "wanutil transaction [options]",
			Description: "Get transaction details and receipt. If it is to a recognized smart contract (in your config file) it will also try to parse the input.",
			Action:      getTransaction,
			Flags:       []cli.Flag{abiFileFlag, hashFlag},
		},
		{
			Name:        "transactionsToAddress",
			Aliases:     []string{"scan"},
			Usage:       "Scan blocks for transactions sent to a given address",
			UsageText:   "wanutil listTransactionsToAddress [options]",
			Description: "Scan blocks for transactions sent to a given address, using an optional block number range.",
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
			Name:        "validate",
			Aliases:     []string{"val"},
			Usage:       "Validate address with checksum",
			UsageText:   "wanutil validateAddress [options]",
			Description: "Check that an address is valid and return with correct capitalizations based on checksum.",
			Action:      validateAddress,
			Flags:       []cli.Flag{addressFlag},
		},
		{
			Name:        "subscribe",
			Aliases:     []string{"sub"},
			Usage:       "Subscribe to events for an address",
			UsageText:   "wanutil subscribe [options]",
			Description: "Subscribe to log events for a given address, starting from an optional block number",
			Action:      subscribe,
			Flags:       []cli.Flag{addressFlag, blockFlag},
		},
	}
)
