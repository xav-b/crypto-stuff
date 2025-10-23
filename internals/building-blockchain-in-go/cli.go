package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// TODO: flag for difficulty mining
// TODO: flag for database
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("\tls - print all the blocks of the blockchain")
	fmt.Println("\treindexutxo - Rebuilds the UTXO set")
	fmt.Println("\tcreatewallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("\twallets - Lists all addresses from the wallet file")
	fmt.Println("\tbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createBlockchain(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	// TODO: overwrite behavior or manually delete the database
	bc := NewBlockchain(address)
	defer bc.db.Close()

	fmt.Println("initializing UTXO set")
	UTXOSet := UTXOSet{bc}
	fmt.Println("reindexing UTXO set")
	UTXOSet.Reindex()

	fmt.Println("Done!")
}

func (cli *CLI) reindexUTXO() {
	bc := NewBlockchain("")
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

func (cli *CLI) getBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := NewBlockchain("")
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}

func (cli *CLI) listAddresses() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CLI) printChain() {
	// TODO: handle better new vs loading blochains. API is bad and there's too
	// much assumptions here
	bc := NewBlockchain("")
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()
		pow := NewProofOfWork(block)

		fmt.Printf("\n============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	fmt.Println("initializing a new transaction")
	bc := NewBlockchain("")
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	fmt.Printf("creating the coinbase tx, reward to %s\n", from)
	cbTx := NewCoinbaseTX(from, "")
	fmt.Printf("creating the actual transaction of %d bitcoins\n", amount)
	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	txs := []*Transaction{cbTx, tx}

	fmt.Println("mining the new block")
	newBlock := bc.AddBlock(txs)
	fmt.Println("updating UTXO set")
	UTXOSet.Update(newBlock)

	fmt.Println("Success!")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	// CLI commands
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	walletsCmd := flag.NewFlagSet("wallets", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	// CLI flags
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	// parse the right flags depending on the command
	switch os.Args[1] {
	case "-h":
		cli.printUsage()
		os.Exit(0)
	case "--help":
		cli.printUsage()
		os.Exit(0)
	case "createblockchain":
		_ = createBlockchainCmd.Parse(os.Args[2:])
	case "ls":
		_ = printChainCmd.Parse(os.Args[2:])
	case "createwallet":
		_ = createWalletCmd.Parse(os.Args[2:])
	case "wallets":
		_ = walletsCmd.Parse(os.Args[2:])
	case "balance":
		_ = getBalanceCmd.Parse(os.Args[2:])
	case "send":
		_ = sendCmd.Parse(os.Args[2:])
	case "reindexutxo":
		_ = reindexUTXOCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	// run the right command

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if walletsCmd.Parsed() {
		cli.listAddresses()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
