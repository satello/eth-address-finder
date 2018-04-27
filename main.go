package main

import (
  "fmt"
  "flag"
  "github.com/onrik/ethrpc"
  "sync"
  "log"
)

func getBlockRequest(blockNumber int, ch chan<-map[string]bool) {
  // TODO consider using same client for each request
  client := ethrpc.New("http://127.0.0.1:8545")

  _, err := client.Web3ClientVersion()
  if err != nil {
      log.Fatal(err)
  }

  // fetch block with transactions
  block, err := client.EthGetBlockByNumber(blockNumber, true)

  if err != nil {
      log.Fatal(err)
  }

  var addresses map[string]bool
  for _, element := range block.Transactions {
    sender := element.From
    receiver := element.To

    if sender != "" {
      addresses[sender] = true
    }
    if receiver != "" {
      addresses[receiver] = true
    }
  }

  ch <- addresses
}

func main() {
  var start_block int
  var end_block int
  var all_addresses map[string]bool
  // command line args
  startPtr := flag.Int("-start", 0, "start block")
  endPtr := flag.Int("-end", 0, "end block")
  flag.Parse()

  start_block = *startPtr
  end_block = *endPtr

  address_chan := make(chan map[string]bool)
  // this will allow our program to stay alive until all requests are completed
  var wg sync.WaitGroup
  wg.Add(end_block-start_block)

  // do each rpc call as a concurrent request
  for block_number := start_block; block_number <= end_block; block_number++ {
    go getBlockRequest(block_number, address_chan)
  }

  go func() {
    for address_map := range address_chan {
      // mark a response as received when we add to our master mapping of addresses
      defer wg.Done()
      for address := range address_map {
        // add address to mapping
        all_addresses[address] = true
      }
    }
  }()

  // wait until we have received all addresses
  wg.Wait()

  // print all addresses NOTE should stdout to a file or you are going to get spammed
  for address := range all_addresses {
    fmt.Print(address)
  }
}
