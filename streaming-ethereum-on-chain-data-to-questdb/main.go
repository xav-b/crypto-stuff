package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v4"
)

const NEW_BLOCK_SQL = `
	INSERT INTO block
	(ts, number, hash, parent_hash, nonce, sha3_uncles, transactions_root, state_root, receipts_root, miner, difficulty, size, gas_limit, gas_used, transaction_count, base_fee_per_gas)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
`

const NEW_TX_SQL = `
	INSERT INTO tx
	(
		tx_hash, tx_index, nonce, tx_type, from_address, to_address, tx_value, tx_cost, input,
		gas, gas_price, max_fee_per_gas, max_priority_fee_per_gas,
		receipt_cumulative_gas_used, receipt_gas_used, receipt_contract_address, receipt_root, receipt_status,
		block_hash, block_number
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
`

func onTx(tx *types.Transaction, receipt *types.Receipt, conn *pgx.Conn) error {
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, _ := signer.Sender(tx)

	log.Printf("storing txn %s\n", tx.Hash().Hex())

	// TODO: receipt.Bloom
	// TODO: receipt.Logs
	_, err := conn.Exec(context.Background(), NEW_TX_SQL,
		tx.Hash().Hex(),
		receipt.TransactionIndex,
		strconv.FormatUint(tx.Nonce(), 10),
		// TODO: does it map to something?
		tx.Type(),
		sender.Hex(),
		tx.To().Hex(),
		tx.Value().String(),
		tx.Cost().Uint64(),
		// gas * gasPrice + value, is it redundant?
		// FIXME: string(tx.Data()),
		nil,

		// gas
		tx.Gas(),
		tx.GasPrice().Uint64(),
		tx.GasFeeCap().Uint64(),
		tx.GasTipCap().Uint64(),

		// receipt
		receipt.CumulativeGasUsed,
		// the fee paid to the miners to process the transaction. the
		// gasUsed is measured in Gwei wich equals 0.000000001 Ether.
		// Equally, 1 Ether equals 1 000 000 000 Gwei
		receipt.GasUsed,
		// smart contract associated to that transaction
		receipt.ContractAddress.Hex(),
		// the root hash of the rootState at the time of the
		// transaction. This is like the hash of the entire blockchain
		// until that moment
		string(receipt.PostState),
		receipt.Status,

		// block
		receipt.BlockHash.Hex(),
		// block number the transaction belongs (i.e. blockchain height)
		receipt.BlockNumber.Uint64(),
	)
	return err
}

func onBlock(block *types.Block, conn *pgx.Conn) error {
	blockTs := time.Unix(int64(block.Time()), 0)

	// convert to timestamp
	log.Printf("storing new block: %s\n", block.Hash().Hex())
	_, err := conn.Exec(context.Background(), NEW_BLOCK_SQL,
		blockTs,
		block.NumberU64(),
		block.Hash().Hex(),
		block.ParentHash().Hex(),
		// nonce is too big for big
		strconv.FormatUint(block.Nonce(), 10),
		block.UncleHash().Hex(),
		// FIXME: string(block.Bloom().Bytes()),
		block.TxHash().Hex(),
		block.Root().Hex(),
		block.ReceiptHash().Hex(),
		block.Coinbase().String(),
		block.Difficulty().Int64(),
		block.Size(),
		// FIXME: string(block.Extra()),
		block.GasLimit(),
		block.GasUsed(),
		len(block.Transactions()),
		block.BaseFee().Int64(),
	)
	// TODO: Uncles() []*Header  // wh UncleHash() then? Root uncle?
	// FIXME: unclassified
	// - MixDigest()
	// - totalDifficulty: accumulated sum of the chain difficulty before the current block

	return err
}

func streamHeaders(client *ethclient.Client, conn *pgx.Conn) {
	ctx := context.Background()

	headers := make(chan *types.Header)
	log.Println("subscribing to blocks' headers")
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatal(err)
	}

	// handling ctrl-c
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Println("listening for events...")
	for {
		select {
		case err := <-sub.Err():
			sub.Unsubscribe()
			log.Fatal(err)
		case <-interrupt:
			log.Println("catch interruption, shutting down")
			sub.Unsubscribe()
			os.Exit(0)
		case header := <-headers:
			log.Printf("received new header: %s\n", header.Hash().Hex())

			log.Println("fetching full block for this header")
			block, err := client.BlockByHash(ctx, header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			if err = onBlock(block, conn); err != nil {
				log.Fatalln("failed to process new block:", err)
			}

			for _, tx := range block.Transactions() {
				log.Println("fetching transaction receipt")
				receipt, _ := client.TransactionReceipt(ctx, tx.Hash())

				if err := onTx(tx, receipt, conn); err != nil {
					log.Fatalln("failed to process tx:", err)
				}
			}
		}
	}
}

func main() {
	ctx := context.Background()

	network := flag.String("network", "mainnet", "infura network to use")
	// TODO: use environment variables as default, and fail if also not find
	pguri := flag.String("db", "", "Postgres-compatible uri")
	blockHash := flag.String("block", "", "Specific block hash")
	flag.Parse()

	if *pguri == "" {
		log.Printf("loading postgres connection info from env")
		*pguri = pgFromEnv()
	}

	log.Printf("connecting to Ethereum: %s\n", *network)
	client, err := ethclient.Dial(ethEndpoint((*network)))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connecting to Postgres: %s\n", *pguri)
	conn, _ := pgx.Connect(ctx, *pguri)
	defer conn.Close(ctx)

	if *blockHash != "" {
		log.Printf("specific block requested: %s\n", *blockHash)
		hash := common.HexToHash(*blockHash)
		count, _ := client.TransactionCount(ctx, hash)

		log.Printf("processing %d transactions\n", count)
		for idx := uint(0); idx < count; idx++ {
			tx, _ := client.TransactionInBlock(ctx, hash, idx)
			log.Println("fetching transaction receipt")
			receipt, _ := client.TransactionReceipt(ctx, tx.Hash())
			if err := onTx(tx, receipt, conn); err != nil {
				log.Fatalln("failed to process tx:", err)
			}
		}
	} else {
		streamHeaders(client, conn)
	}
}
