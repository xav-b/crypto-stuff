# Building Blockchain in Go

Credits to the author: [Ivan Kuznetsov](https://jeiwan.net/)

[Part 1: Basic Prototype](https://jeiwan.net/posts/building-blockchain-in-go-part-1/)
[Part 2: Proof-of-Work](https://jeiwan.net/posts/building-blockchain-in-go-part-2/)
[Part 3: Persistence and CLI](https://jeiwan.net/posts/building-blockchain-in-go-part-3/)
[Part 4: Transactions 1](https://jeiwan.net/posts/building-blockchain-in-go-part-4/)
[Part 5: Addresses](https://jeiwan.net/posts/building-blockchain-in-go-part-5/)
[Part 6: Transactions 2](https://jeiwan.net/posts/building-blockchain-in-go-part-6/)

## Installation

```console
# go dependencies
go get github.com/boltdb/bolt/...

# should have the boltdb cli installed now
bolt help
```

## Usage

```go
$ go build
$ ./blockchain -help
$ ./blockchain createblockchain --address Xavier
$ ./blockchain ls
$ ./blockchain balance -address Xavier
$ ./blockchain send -from Xavier -to Pedro -amount 6
```
