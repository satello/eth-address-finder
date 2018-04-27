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

  addresses := make(map[string]bool)
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
  all_addresses := make(map[string]bool)
  // command line args
  startPtr := flag.Int("start", 0, "start block")
  endPtr := flag.Int("end", 0, "end block")
  flag.Parse()

  start_block = *startPtr
  end_block = *endPtr

  address_chan := make(chan map[string]bool)
  // this will allow our program to stay alive until all requests are completed
  var wg sync.WaitGroup
  wg.Add(end_block-start_block + 1)

  // limit number of go routines running at once so we don't go over open file limit
  current_block := start_block
  // do each rpc call as a concurrent request
  for i := 0; i <= 500; i++ {
    current_block++
    go getBlockRequest(current_block, address_chan)
  }

  go func() {
    for address_map := range address_chan {
      // mark a response as received when we add to our master mapping of addresses
      for address := range address_map {
        // add address to mapping
        all_addresses[address] = true
      }

      // start a new request once one has finished
      if current_block < end_block {
        current_block++
        go getBlockRequest(current_block, address_chan)
      }
      wg.Done()
    }
  }()

  // wait until we have received all addresses
  wg.Wait()

  // print all addresses NOTE should stdout to a file or you are going to get spammed
  for address := range all_addresses {
    fmt.Println(address)
  }
}
