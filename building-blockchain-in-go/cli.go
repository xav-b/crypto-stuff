package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// TODO: flag for difficulty mining
// TODO: flag for database
type CLI struct {
	bc *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("\tls - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("ls", flag.ExitOnError)

	// add the block's data as position argument
	// example: `blockchain_go addblock "Pay 0.031337 for a coffee"`
	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "-h":
		cli.printUsage()
		os.Exit(0)
	case "--help":
		cli.printUsage()
		os.Exit(0)
	case "addblock":
		_ = addBlockCmd.Parse(os.Args[2:])
	case "ls":
		_ = printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Block successfully added")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()
		pow := NewProofOfWork(block)

		fmt.Println()
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
