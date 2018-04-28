package main

import (
  "fmt"
  "flag"
  "github.com/onrik/ethrpc"
  "sync"
  "log"
  "time"
)

func getBlockRequest(blockNumber int, ch chan<-bool, client *ethrpc.EthRPC) {
  // fetch block with transactions
  block, err := client.EthGetBlockByNumber(blockNumber, true)

  if err != nil {
      log.Fatal(err)
  }

  for _, element := range block.Transactions {
    sender := element.From
    receiver := element.To

    if sender != "" {
      fmt.Println(sender)
    }
    if receiver != "" {
      fmt.Println(receiver)
    }
  }

  ch <- true
}

func main() {
  start := time.Now()
  var start_block int
  var end_block int
  all_addresses := make(map[string]bool)
  // command line args
  startPtr := flag.Int("start", 0, "start block")
  endPtr := flag.Int("end", 0, "end block")
  flag.Parse()

  start_block = *startPtr
  end_block = *endPtr

  client := ethrpc.New("http://127.0.0.1:8545")

  _, err := client.Web3ClientVersion()
  if err != nil {
      log.Fatal(err)
  }

  block_done_chan := make(chan bool)
  // this will allow our program to stay alive until all requests are completed
  var wg sync.WaitGroup
  block_range_size := end_block-start_block
  wg.Add(block_range_size + 1)

  // limit number of go routines running at once so we don't go over open file limit
  current_block := start_block
  batch_size := 500
  if block_range_size < batch_size {
    batch_size = block_range_size
  }
  // do each rpc call as a concurrent request
  for i := 0; i <= batch_size; i++ {
    current_block++
    go getBlockRequest(current_block, block_done_chan, client)
  }

  go func() {
    for success := range block_done_chan {
      // start a new request once one has finished
      if success {
        if current_block <= end_block {
          current_block++
          go getBlockRequest(current_block, block_done_chan, client)
        }
        wg.Done()
      }
    }
  }()

  // wait until we have received all addresses
  wg.Wait()
  // print all addresses NOTE should stdout to a file or you are going to get spammed
  fmt.Println()
  fmt.Println("All done")
  fmt.Println(time.Since(start))

  for address := range all_addresses {
    fmt.Println(address)
  }
}
