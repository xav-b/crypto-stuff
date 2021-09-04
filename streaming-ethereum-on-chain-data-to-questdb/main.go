package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v4"
)

const ETH_WS_ENDPOINT = "wss://ropsten.infura.io/ws"

func ethEndpoint(network string) string {
	return fmt.Sprintf("wss://%s.infura.io/ws/v3/%s", network, os.Getenv("INFURA_PROJECT_ID"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func pgFromEnv() string {
	user := getEnv("PGUSER", "localhost")
	password := getEnv("PGPASSWORD", "root")
	host := getEnv("PGHOST", "localhost")
	port := getEnv("PGPORT", "5432")
	dbname := os.Getenv("PGDATABASE")

	// validation
	if dbname == "" {
		log.Fatalln("no database name exported: export PGDATABASE=xxxxx")
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	// const pguri = "postgresql://postgres:RDLPWbx5hM3ra@localhost:5432/crypto"
}

func main() {
	ctx := context.Background()

	network := flag.String("network", "mainnet", "infura network to use")
	// TODO: use environment variables as default, and fail if also not find
	pguri := flag.String("db", "", "Postgres-compatible uri")
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

	// const pguri = "postgresql://admin:quest@localhost:8812/qdb"
	// const pguri = "postgresql://postgres:RDLPWbx5hM3ra@localhost:5432/crypto"
	log.Printf("connecting to Postgres: %s\n", *pguri)
	conn, _ := pgx.Connect(ctx, *pguri)
	defer conn.Close(ctx)

	_, err = conn.Prepare(ctx, "newblock", `
		INSERT INTO block
		(ts, number, hash, parent_hash, nonce, sha3_uncles, transactions_root, state_root, receipts_root, miner, difficulty, size, gas_limit, gas_used, transaction_count, base_fee_per_gas)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`)
	if err != nil {
		log.Fatalln(err)
	}

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

			blockTs := time.Unix(int64(block.Time()), 0)

			log.Println("== block info:", header.Hash().Hex())
			fmt.Printf("Timestamp: %d (%s)\n", block.Time(), blockTs)
			fmt.Printf("Block number: %d\n", block.NumberU64())
			fmt.Printf("Parent hash: %s\n", block.ParentHash())
			fmt.Printf("Nonce: %d\n", block.Nonce())
			fmt.Printf("Uncle hash: %s\n", block.UncleHash())
			fmt.Printf("Logs bloom: %x\n", block.Bloom().Bytes())
			fmt.Printf("Transactions root hash: %s\n", block.TxHash())
			fmt.Printf("State root hash: %s\n", block.Root())
			fmt.Printf("Receipt hash: %s\n", block.ReceiptHash()) // receipts root == receipt hash?
			fmt.Printf("Coinbase address (miner): %s\n", block.Coinbase())
			fmt.Printf("Difficulty: %d\n", block.Difficulty().Int64())
			fmt.Printf("Block size: %.2f\n", block.Size())
			fmt.Printf("Extra: %s\n", string(block.Extra()))
			fmt.Printf("Gas limit: %d\n", block.GasLimit())
			fmt.Printf("Gas used: %d\n", block.GasUsed())
			fmt.Printf("transactions: %d\n", len(block.Transactions()))
			fmt.Printf("Base fee: %d\n", block.BaseFee()) // per gas

			// FIXME: unclassified
			// MixDigest()
			// totalDifficulty: accumulated sum of the chain difficulty before the current block
			// TODO: Uncles() []*Header  // wh UncleHash() then? Root uncle?

			/* TODO: all the stuff: https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.8/core/types#Transaction.AccessList
			for idx, transaction := range block.Transactions() {
				fmt.Printf("\ttransaction #%d: %s\n", idx, transaction.Hash())
			}*/

			fmt.Println()

			// convert to timestamp
			log.Printf("storing new block")
			_, err = conn.Exec(ctx, "newblock",
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
			if err != nil {
				log.Fatalln("failed to insert row:", err)
			}
		}
	}
}
