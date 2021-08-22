package main

import (
	"flag"

	"github.com/xav-b/goinfura"
)

func main() {
	network := flag.String("network", "mainnet", "Ethereum network to connect to")
	flag.Parse()

	println(goinfura.Endpoint(*network))
}
