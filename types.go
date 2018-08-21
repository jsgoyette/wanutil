package main

import (
	"github.com/wanchain/go-wanchain/accounts/abi"
)

type AbiField struct {
	Type      string
	Name      string
	Constant  bool
	Indexed   bool
	Anonymous bool
	Inputs    []abi.Argument
	Outputs   []abi.Argument
}

type AbiMethod struct {
	AbiField
	Signature     string
	SignatureHash string
}
